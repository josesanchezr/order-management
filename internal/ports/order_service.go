package ports

import (
	"order_management/internal/models"
)

// OrderService define los métodos disponibles para manejar órdenes.
type OrderService interface {
	CreateOrder(order *models.Order) error
	GetOrderById(id uint) (*models.Order, error)
}
