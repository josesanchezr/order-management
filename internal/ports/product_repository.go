package ports

import (
	"order_management/internal/models"

	"gorm.io/gorm"
)

// ProductRepository define las operaciones que pueden realizarse sobre la entidad Product
type ProductRepository interface {
	GetAll() ([]models.Product, error)
	GetByID(id uint, tx *gorm.DB) (*models.Product, error)
	Update(product *models.Product) error
	UpdateStock(id uint, newStock int, tx *gorm.DB) error
}
