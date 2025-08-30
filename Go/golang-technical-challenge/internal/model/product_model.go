package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateProductRequest struct {
	ItemName   string          `json:"item_name" validate:"required,min=5,max=255"`
	Quantity   int             `json:"quantity" validate:"required,min=1"`
	TotalCost  decimal.Decimal `json:"total_cost" validate:"required"`
	TotalPrice decimal.Decimal `json:"total_price" validate:"required"`
}

type ProductResponse struct {
	ID         string          `json:"id"`
	ItemName   string          `json:"item_name"`
	Quantity   int             `json:"quantity"`
	TotalCost  decimal.Decimal `json:"total_cost"`
	TotalPrice decimal.Decimal `json:"total_price"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}
