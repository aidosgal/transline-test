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

type Usecase interface{
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
	log := u.log.With("method", "GetCustomer")
	
	log.Debug("started storage.GetCustomerByIDN", slog.String("idn", idn))
	
	customer, err := u.storage.GetCustomerByIDN(ctx, idn)
	if err != nil {
		log.Error("failed to storage.GetCustomerByIDN", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to storage.GetCustomerByIDN: %w", err)
	}
	
	return customer, nil
}

func (u *usecase) UpsertCustomer(ctx context.Context, idn string) (*entity.Customer, error) {
	log := u.log.With("method", "UpsertCustomer")
	
	log.Debug("started storage.UpsertCustomer", slog.String("idn", idn))
	
	customer, err := u.storage.UpsertCustomer(ctx, idn)
	if err != nil {
		log.Error("failed to storage.UpsertCustomer", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to storage.UpsertCustomer: %w", err)
	}
	
	return customer, nil
}
