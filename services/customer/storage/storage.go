package storage

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/aidosgal/transline-test/services/customer/entity"
)

type storage struct {
	log *slog.Logger
	db  *sql.DB
}

type Storage interface {
	GetCustomerByIDN(ctx context.Context, idn string) (*entity.Customer, error)
}

func New(log *slog.Logger, db *sql.DB) Storage {
	return &storage{
		log: log,
		db:  db,
	}
}

func (s *storage) GetCustomerByIDN(ctx context.Context, idn string) (*entity.Customer, error) {
	return nil, nil
}
