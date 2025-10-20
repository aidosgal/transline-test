package entity

import (
	"time"
)

type (
	Shipment struct {
		ID         string    `json:"id"`
		Route      string    `json:"route"`
		Price      int       `json:"price"`
		Status     string    `json:"status"`
		CustomerID string    `json:"customer_id"`
		CreatedAt  time.Time `json:"created_at"`
	}
)
