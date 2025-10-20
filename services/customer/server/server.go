package server

import (
	"context"

	"github.com/aidosgal/transline-test/services/customer/entity"
	"github.com/aidosgal/transline-test/services/customer/usecase"
	pb "github.com/aidosgal/transline-test/specs/proto/customer"
)

type server struct {
	pb.UnimplementedCustomerServer
	usecase usecase.Usecase
}

func New (usecase usecase.Usecase) *server {
	return &server{
		usecase: usecase,
	}
}

func (s *server) UpsertCustomer(ctx context.Context, req *pb.UpsertCustomerRequest) (*pb.CustomerResponse, error) {
	resp, err := s.usecase.UpsertCustomer(ctx, req.GetIdn())
	if err != nil {
		return nil, err
	}
	
	return entity.MakeCustomerEntityToPb(resp), nil
}

func (s *server) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.CustomerResponse, error) {
	resp, err := s.usecase.GetCustomer(ctx, req.Idn)
	if err != nil {
		return nil, err
	}
	
	return entity.MakeCustomerEntityToPb(resp), nil
}
