package server

import (
	"context"
	"log/slog"

	"github.com/aidosgal/transline-test/services/customer/entity"
	"github.com/aidosgal/transline-test/services/customer/usecase"
	pb "github.com/aidosgal/transline-test/specs/proto/customer"
)

type server struct {
	pb.UnimplementedCustomerServer
	log     *slog.Logger
	usecase usecase.Usecase
}

func New(log *slog.Logger, usecase usecase.Usecase) *server {
	return &server{
		log:     log.With("layer", "server"),
		usecase: usecase,
	}
}

func (s *server) UpsertCustomer(ctx context.Context, req *pb.UpsertCustomerRequest) (*pb.CustomerResponse, error) {
	log := s.log.With("method", "UpsertCustomer")

	log.InfoContext(ctx, "received upsert customer request", slog.String("idn", req.GetIdn()))

	resp, err := s.usecase.UpsertCustomer(ctx, req.GetIdn())
	if err != nil {
		log.ErrorContext(ctx, "failed to upsert customer", slog.String("error", err.Error()))
		return nil, err
	}

	log.InfoContext(ctx, "customer upserted successfully",
		slog.String("customer_id", resp.ID),
		slog.String("idn", resp.IDN))
	return entity.MakeCustomerEntityToPb(resp), nil
}

func (s *server) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.CustomerResponse, error) {
	log := s.log.With("method", "GetCustomer")

	log.InfoContext(ctx, "received get customer request", slog.String("idn", req.Idn))

	resp, err := s.usecase.GetCustomer(ctx, req.Idn)
	if err != nil {
		log.ErrorContext(ctx, "failed to get customer", slog.String("error", err.Error()))
		return nil, err
	}

	log.InfoContext(ctx, "customer retrieved successfully",
		slog.String("customer_id", resp.ID),
		slog.String("idn", resp.IDN))
	return entity.MakeCustomerEntityToPb(resp), nil
}
