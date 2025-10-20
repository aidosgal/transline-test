package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/aidosgal/transline-test/services/customer/entity"
)

type storage struct {
	log *slog.Logger
	db  *sql.DB
}

type Storage interface {
	GetCustomerByIDN(ctx context.Context, idn string) (*entity.Customer, error)
	UpsertCustomer(ctx context.Context, idn string) (*entity.Customer, error)
}

func New(log *slog.Logger, db *sql.DB) Storage {
	return &storage{
		log: log.With("layer", "Storage"),
		db:  db,
	}
}

func (s *storage) GetCustomerByIDN(ctx context.Context, idn string) (*entity.Customer, error) {
	log := s.log.With("method", "GetCustomerByIDN")
	customer := &entity.Customer{}

	log.Debug("select query started", slog.String("idn", idn))
	err := s.db.QueryRowContext(ctx, `SELECT id, idn, created_at FROM customers WHERE idn=$1`, idn).
		Scan(&customer.ID, &customer.IDN, &customer.CreatedAt)
	if err != nil {
		log.Error("failed to SELECT id, idn, created_at FROM customers WHERE idn=%s", idn,
			slog.String("idn", idn),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	log.Debug("finished query", slog.Any("customer", customer))
	
	return customer, nil
}

func (s *storage) UpsertCustomer(ctx context.Context, idn string) (*entity.Customer, error) {
	log := s.log.With("method", "UpsertCustomer")
	customer := &entity.Customer{}

	query := `
		INSERT INTO customers (id, idn)
		VALUES (gen_random_uuid(), $1)
		ON CONFLICT (idn) DO UPDATE SET idn = EXCLUDED.idn
		RETURNING id, idn, created_at;
	`

	log.Debug("executing upsert", slog.String("idn", idn))

	err := s.db.QueryRowContext(ctx, query, idn).Scan(
		&customer.ID,
		&customer.IDN,
		&customer.CreatedAt,
	)
	if err != nil {
		log.Error("upsert failed",
			slog.String("idn", idn),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to upsert customer: %w", err)
	}

	log.Debug("upsert successful", slog.Any("customer", customer))
	return customer, nil
}
