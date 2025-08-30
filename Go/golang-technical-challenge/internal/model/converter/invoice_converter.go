package converter

import (
	"golang-technical-challenge/internal/entity"
	"golang-technical-challenge/internal/model"
)

func InvoiceToResponse(invoice *entity.Invoice) *model.InvoiceResponse {
	return &model.InvoiceResponse{
		InvoiceNo:       invoice.InvoiceNo,
		Date:            invoice.Date,
		CustomerName:    invoice.CustomerName,
		SalespersonName: invoice.SalespersonName,
		PaymentType:     invoice.PaymentType,
		Notes:           invoice.Notes,
		CreatedAt:       invoice.CreatedAt,
		UpdatedAt:       invoice.UpdatedAt,
		Products:        ProductsToResponseList(invoice.Products),
	}
}

func InvoicesToResponseList(invoices []entity.Invoice) []model.InvoiceResponse {
	responses := make([]model.InvoiceResponse, len(invoices))
	for i, inv := range invoices {
		responses[i] = *InvoiceToResponse(&inv)
	}
	return responses
}
