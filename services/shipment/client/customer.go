package client

import (
	"fmt"

	"github.com/aidosgal/transline-test/pkg/config"
	"github.com/aidosgal/transline-test/specs/proto/customer"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CustomerClient struct {
	conn *grpc.ClientConn
	customer.CustomerClient
}

func New(cfg *config.Config) (*CustomerClient, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", "customer-service", cfg.CustomerService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc connection: %w", err)
	}

	return &CustomerClient{
		conn:           conn,
		CustomerClient: customer.NewCustomerClient(conn),
	}, nil
}

func (c *CustomerClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
