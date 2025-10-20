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
		log: log.With("layer", "usecase"),
		storage: storage,
		customer: customer,
	}
}

func (u *usecase) CreateShipment(ctx context.Context, req *entity.CreateReq) (*entity.CreateResp, error) {
	log := u.log.With("method", "CreateShipment")
	customer, err := u.customer.UpsertCustomer(ctx, &customer.UpsertCustomerRequest{
		Idn: req.Customer.IDN,
	})
	if err != nil {
		log.Error("failed to customer.UpsertCustomer", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to customer.UpsertCustomer: %w", err)
	}

	shipmentID, err := u.storage.CreateShipment(ctx, req, customer.Id)
	if err != nil {
		log.Error("failed to storage.CreateShipment", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to storage.CreateShipment: %w", err)
	}

	shipment, err := u.storage.GetShipment(ctx, shipmentID)
	if err != nil {
		log.Error("failed to storage.GetShipment", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to storage.GetShipment: %w", err)
	}
	
	return &entity.CreateResp{
		Shipment: *shipment,
	}, nil
}

func (u *usecase) GetShipment(ctx context.Context, id string) (*entity.Shipment, error) {
	log := u.log.With("method", "GetShipment")
	shipment, err := u.storage.GetShipment(ctx, id)
	if err != nil {
		log.Error("failed to storage.GetShipment", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to storage.GetShipment: %w", err)
	}

	return shipment, nil
}
