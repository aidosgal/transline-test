package usecase

import (
	"context"
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
		log:     log,
		storage: storage,
	}
}

func (u *usecase) GetCustomer(ctx context.Context, idn string) (*entity.Customer, error) {
	return nil, nil
}

func (u *usecase) UpsertCustomer(ctx context.Context, idn string) (*entity.Customer, error) {
	return nil, nil
}
