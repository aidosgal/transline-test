package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/aidosgal/transline-test/services/shipment/entity"
)

type storage struct {
	log *slog.Logger
	db  *sql.DB
}

type Storage interface {
	GetShipment(ctx context.Context, id string) (*entity.Shipment, error)
	CreateShipment(ctx context.Context, req *entity.CreateReq, customerID string) (string, error)
}

func New(log *slog.Logger, db *sql.DB) Storage {
	return &storage{
		log: log.With("layer", "storage"),
		db:  db,
	}
}

func (s *storage) GetShipment(ctx context.Context, id string) (*entity.Shipment, error) {
	log := s.log.With("method", "CreateShipment")

	shipment := &entity.Shipment{}
	err := s.db.QueryRowContext(ctx, `SELECT route, price, status, customer_id, created_at FROM shipments WHERE id=$1`, id).
		Scan(&shipment.Route, &shipment.Price, &shipment.Status, &shipment.CustomerID, &shipment.CreatedAt)
	if err != nil {
		log.Error("failed db select shipment", slog.String("id", id), slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}
	
	return shipment, nil
}

func (s *storage) CreateShipment(ctx context.Context, req *entity.CreateReq, customerID string) (string, error) {
	log := s.log.With("method", "CreateShipment")
	var shipID string
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO shipments (route, price, customer_id) VALUES ($1,$2,$3) RETURNING id`,
		req.Route, req.Price, customerID).Scan(&shipID)
	if err != nil {
		log.Error("failed db insert shipment", slog.String("error", err.Error()))
		return "", fmt.Errorf("failed db insert shipment: %w", err)
	}
	
	return shipID, nil
}
