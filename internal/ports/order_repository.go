package ports

import (
	"order_management/internal/models"

	"gorm.io/gorm"
)

// OrderRepository define las operaciones disponibles para gestionar órdenes.
type OrderRepository interface {
	Create(order *models.Order, tx *gorm.DB) error
	FindByID(id uint) (*models.Order, error)
}
