package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aidosgal/transline-test/services/customer/entity"
	"github.com/aidosgal/transline-test/services/customer/storage"
)

type usecase struct {
	log     *slog.Logger
	storage storage.Storage
}

type Usecase interface {
	GetCustomer(ctx context.Context, idn string) (*entity.Customer, error)
	UpsertCustomer(ctx context.Context, idn string) (*entity.Customer, error)
}

func New(log *slog.Logger, storage storage.Storage) Usecase {
	return &usecase{
		log:     log.With("layer", "usecase"),
		storage: storage,
	}
}

func (u *usecase) GetCustomer(ctx context.Context, idn string) (*entity.Customer, error) {
	log := u.log.With("method", "GetCustomer", "idn", idn)

	log.Info("getting customer from storage")

	customer, err := u.storage.GetCustomerByIDN(ctx, idn)
	if err != nil {
		log.Error("failed to get customer from storage", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to storage.GetCustomerByIDN: %w", err)
	}

	log.Info("customer found successfully",
		slog.String("customer_id", customer.ID),
		slog.String("created_at", customer.CreatedAt.String()))
	return customer, nil
}

func (u *usecase) UpsertCustomer(ctx context.Context, idn string) (*entity.Customer, error) {
	log := u.log.With("method", "UpsertCustomer", "idn", idn)

	log.Info("upserting customer in storage")

	customer, err := u.storage.UpsertCustomer(ctx, idn)
	if err != nil {
		log.Error("failed to upsert customer in storage", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to storage.UpsertCustomer: %w", err)
	}

	log.Info("customer upserted successfully",
		slog.String("customer_id", customer.ID),
		slog.String("idn", customer.IDN),
		slog.String("created_at", customer.CreatedAt.String()))
	return customer, nil
}
