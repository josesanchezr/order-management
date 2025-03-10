package mappers

import (
	"order_management/internal/dtos"
	"order_management/internal/models"
)

func ConvertProductToProductResponseDTO(products []models.Product) []dtos.ProductResponseDTO {
	var productResponseDTOs []dtos.ProductResponseDTO

	for _, product := range products {
		productResponseDTOs = append(productResponseDTOs, dtos.ProductResponseDTO{
			ID:    product.ID,
			Name:  product.Name,
			Price: product.Price,
			Stock: product.Stock,
		})
	}

	return productResponseDTOs
}
