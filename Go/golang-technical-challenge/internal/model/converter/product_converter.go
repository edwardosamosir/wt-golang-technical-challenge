package converter

import (
	"golang-technical-challenge/internal/entity"
	"golang-technical-challenge/internal/model"
)

func ProductToResponse(product *entity.Product) model.ProductResponse {
	return model.ProductResponse{
		ID:         product.ID,
		ItemName:   product.ItemName,
		Quantity:   product.Quantity,
		TotalCost:  product.TotalCost,
		TotalPrice: product.TotalPrice,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
	}
}

func ProductsToResponseList(products []entity.Product) []model.ProductResponse {
	responses := make([]model.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = ProductToResponse(&p)
	}
	return responses
}
