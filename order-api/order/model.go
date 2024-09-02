package order

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID          uint64     `json:"id"`
	CustomerID  uuid.UUID  `json:"customer_id"`
	LineItems   []LineItem `json:"line_items"`
	Total       float32    `json:"total"`
	CreatedAt   *time.Time `json:"created_at"`
	ShippedAt   *time.Time `json:"shipped_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Currency    string     `json:"currency"`
}

type LineItem struct {
	ID       uuid.UUID `json:"id"`
	Quantity int       `json:"quantity"`
	Price    float32   `json:"price"`
}
