package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/aidosgal/transline-test/pkg/config"
	"github.com/aidosgal/transline-test/services/customer/server"
	"github.com/aidosgal/transline-test/services/customer/storage"
	"github.com/aidosgal/transline-test/services/customer/usecase"
	pb "github.com/aidosgal/transline-test/specs/proto/customer"
	"github.com/golang-migrate/migrate/v4"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"

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

	postgresURL := cfg.CustomerService.Postgres.BuildPostgresURL()
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

	migrationURL := cfg.CustomerService.Postgres.BuildPostgresMigrationURL()
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

	customerStorage := storage.New(log, db)
	customerUsecase := usecase.New(log, customerStorage)
	customerServer := server.New(customerUsecase)

	address := fmt.Sprintf(":%d", cfg.CustomerService.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Error("failed to listen", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("gRPC server listening", slog.String("address", address))

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	pb.RegisterCustomerServer(grpcServer, customerServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("server stopped", slog.String("error", err.Error()))
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server gracefully...")
	grpcServer.GracefulStop()
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
			semconv.ServiceName("customer-service"),
		)),
	)

	log.Info("OpenTelemetry tracer initialized", slog.String("endpoint", "otel-collector:4318"))
	return tp, nil
}
