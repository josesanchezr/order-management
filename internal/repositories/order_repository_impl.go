package repositories

import (
	"order_management/internal/models"
	"order_management/internal/ports"

	"gorm.io/gorm"
)

// OrderRepositoryImpl implementa OrderRepository usando GORM.
type OrderRepositoryImpl struct {
	db *gorm.DB
}

// NewOrderRepository crea una nueva instancia de OrderRepositoryImpl.
func NewOrderRepository(db *gorm.DB) ports.OrderRepository {
	return &OrderRepositoryImpl{db: db}
}

// Create inserta una nueva orden en la base de datos.
func (r *OrderRepositoryImpl) Create(order *models.Order, tx *gorm.DB) error {
	return tx.Create(order).Error
}

// FindByID busca una orden por ID.
func (r *OrderRepositoryImpl) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("OrderItems.Product").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
