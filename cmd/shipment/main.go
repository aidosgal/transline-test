package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aidosgal/transline-test/pkg/config"
	"github.com/aidosgal/transline-test/services/shipment/client"
	"github.com/aidosgal/transline-test/services/shipment/server"
	"github.com/aidosgal/transline-test/services/shipment/storage"
	"github.com/aidosgal/transline-test/services/shipment/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-migrate/migrate/v4"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp, err := setupTracer(ctx, log)
	if err != nil {
		log.Error("failed to initialize tracer", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Error("failed to shutdown tracer", slog.String("error", err.Error()))
		}
	}()

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	postgresURL := cfg.Shipment.Postgres.BuildPostgresURL()
	log.Info("connecting to database", slog.String("url", postgresURL))

	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		log.Error("failed to open db", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Error("failed to connect to db", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("connected to database")

	migrationURL := cfg.Shipment.Postgres.BuildPostgresMigrationURL()
	migrationPath := "/app/migrations"

	m, err := migrate.New("file://"+migrationPath, migrationURL)
	if err != nil {
		log.Error("failed to init migrations service", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("no migrations to apply", slog.String("error", err.Error()))
		}
	}

	log.Info("migration applied successfully")

	customerClient, err := client.New(cfg)
	if err != nil {
		log.Error("failed to connect to customer GRPC", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer customerClient.Close()

	shipmentStorage := storage.New(log, db)
	shipmentUsecase := usecase.New(log, shipmentStorage, customerClient)
	shipmentServer := server.New(shipmentUsecase)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.URLFormat)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Route("/api/v1", func(apiRouter chi.Router) {
		apiRouter.Route("/shipments", func(authRouter chi.Router) {
			authRouter.Post("/", shipmentServer.CreateShipment)
			authRouter.Post("/{id}", shipmentServer.GetShipment)
		})
	})

	wrappedChi := otelhttp.NewHandler(router, "shipment-service")

	address := fmt.Sprintf(":%d", 8080)
	server := &http.Server{
		Addr:    address,
		Handler: wrappedChi,
	}

	go func() {
		log.Info("HTTP server listenipng", slog.String("address", address))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", slog.String("error", err.Error()))
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", slog.String("error", err.Error()))
	}

	log.Info("server stopped")
}

func setupTracer(ctx context.Context, log *slog.Logger) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("otel-collector:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("shipment-service"),
		)),
	)

	log.Info("OpenTelemetry tracer initialized", slog.String("endpoint", "otel-collector:4318"))
	return tp, nil
}
