package model

import (
	"time"
)

type InvoiceResponse struct {
	InvoiceNo       string            `json:"invoice_no"`
	Date            time.Time         `json:"date"`
	CustomerName    string            `json:"customer_name"`
	SalespersonName string            `json:"salesperson_name"`
	PaymentType     string            `json:"payment_type"`
	Notes           *string           `json:"notes,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Products        []ProductResponse `json:"products"`
}

type InvoiceListResponse struct {
	Invoices    []InvoiceResponse `json:"invoices"`
	TotalProfit string            `json:"total_profit"`
	TotalCash   string            `json:"total_cash"`
	Paging      PageMetadata      `json:"paging"`
}

type CreateInvoiceRequest struct {
	InvoiceNo       string                 `json:"invoice_no" validate:"required,max=50"`
	Date            string                 `json:"date" validate:"required,datetime=2006-01-02"`
	CustomerName    string                 `json:"customer_name" validate:"required,min=2,max=255"`
	SalespersonName string                 `json:"salesperson_name" validate:"required,min=2,max=255"`
	PaymentType     string                 `json:"payment_type" validate:"required,oneof=CASH CREDIT"`
	Notes           *string                `json:"notes,omitempty" validate:"omitempty,min=5"`
	Products        []CreateProductRequest `json:"products" validate:"required,dive"`
}

type UpdateInvoiceRequest struct {
	Date            string                 `json:"date" validate:"required,datetime=2006-01-02"`
	CustomerName    string                 `json:"customer_name" validate:"required,min=2,max=255"`
	SalespersonName string                 `json:"salesperson_name" validate:"required,min=2,max=255"`
	PaymentType     string                 `json:"payment_type" validate:"required,oneof=CASH CREDIT"`
	Notes           *string                `json:"notes,omitempty" validate:"omitempty,min=5"`
	Products        []CreateProductRequest `json:"products" validate:"required,dive"`
}

type DeleteInvoiceRequest struct {
	InvoiceNo string `json:"-" validate:"required"`
}

type ImportError struct {
	InvoiceNo string `json:"invoice_no"`
	Message   string `json:"message"`
}
