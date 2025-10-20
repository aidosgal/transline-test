package entity

import (
	"time"

	customerv1 "github.com/aidosgal/transline-test/specs/proto/customer"
)

type (
	Customer struct {
		ID        string    `json:"id"`
		IDN       string    `json:"idn"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func MakeCustomerEntityToPb(customer *Customer) *customerv1.CustomerResponse {
	return &customerv1.CustomerResponse{
		Id: customer.ID,
		Idn: customer.IDN,
		CreatedAt: customer.CreatedAt.String(),
	}
}
