package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aidosgal/transline-test/services/shipment/client"
	"github.com/aidosgal/transline-test/services/shipment/entity"
	"github.com/aidosgal/transline-test/services/shipment/storage"
	"github.com/aidosgal/transline-test/specs/proto/customer"
)

type usecase struct {
	log      *slog.Logger
	storage  storage.Storage
	customer *client.CustomerClient
}

type Usecase interface {
	CreateShipment(ctx context.Context, req *entity.CreateReq) (*entity.CreateResp, error)
	GetShipment(ctx context.Context, id string) (*entity.Shipment, error)
}

func New(log *slog.Logger, storage storage.Storage, customer *client.CustomerClient) Usecase {
	return &usecase{
		log:      log.With("layer", "usecase"),
		storage:  storage,
		customer: customer,
	}
}

func (u *usecase) CreateShipment(ctx context.Context, req *entity.CreateReq) (*entity.CreateResp, error) {
	log := u.log.With("method", "CreateShipment",
		"route", req.Route,
		"price", req.Price,
		"customer_idn", req.Customer.IDN)

	log.InfoContext(ctx, "starting shipment creation process")

	log.InfoContext(ctx, "calling customer service to upsert customer")
	customer, err := u.customer.UpsertCustomer(ctx, &customer.UpsertCustomerRequest{
		Idn: req.Customer.IDN,
	})
	if err != nil {
		log.ErrorContext(ctx, "failed to upsert customer via gRPC", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to customer.UpsertCustomer: %w", err)
	}
	log.InfoContext(ctx, "customer upserted via gRPC", slog.String("customer_id", customer.Id))

	log.InfoContext(ctx, "saving shipment to storage")
	shipmentID, err := u.storage.CreateShipment(ctx, req, customer.Id)
	if err != nil {
		log.ErrorContext(ctx, "failed to save shipment to storage", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to storage.CreateShipment: %w", err)
	}
	log.InfoContext(ctx, "shipment saved to storage", slog.String("shipment_id", shipmentID))

	log.InfoContext(ctx, "retrieving created shipment")
	shipment, err := u.storage.GetShipment(ctx, shipmentID)
	if err != nil {
		log.ErrorContext(ctx, "failed to retrieve created shipment", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to storage.GetShipment: %w", err)
	}

	log.InfoContext(ctx, "shipment creation completed successfully",
		slog.String("shipment_id", shipment.ID),
		slog.String("customer_id", shipment.CustomerID),
		slog.String("status", shipment.Status))

	return &entity.CreateResp{
		Shipment: *shipment,
	}, nil
}

func (u *usecase) GetShipment(ctx context.Context, id string) (*entity.Shipment, error) {
	log := u.log.With("method", "GetShipment", "shipment_id", id)

	log.InfoContext(ctx, "retrieving shipment from storage")

	shipment, err := u.storage.GetShipment(ctx, id)
	if err != nil {
		log.ErrorContext(ctx, "failed to retrieve shipment from storage", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to storage.GetShipment: %w", err)
	}

	log.InfoContext(ctx, "shipment retrieved successfully",
		slog.String("route", shipment.Route),
		slog.String("status", shipment.Status),
		slog.Int("price", shipment.Price))

	return shipment, nil
}
