package ports

import (
	"order_management/internal/models"
)

type ProductService interface {
	GetAllProducts() ([]models.Product, error)
	UpdateStock(id uint, stock int) error
}
