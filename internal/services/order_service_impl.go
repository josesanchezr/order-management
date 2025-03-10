package services

import (
	"errors"
	"log"
	"order_management/internal/models"
	"order_management/internal/ports"

	"gorm.io/gorm"
)

// OrderServiceImpl implementa OrderService.
type OrderServiceImpl struct {
	repo        ports.OrderRepository
	productRepo ports.ProductRepository
	db          *gorm.DB
}

// NewOrderService crea una nueva instancia de OrderService.
func NewOrderService(repo ports.OrderRepository, productRepo ports.ProductRepository, db *gorm.DB) ports.OrderService {
	return &OrderServiceImpl{repo: repo, productRepo: productRepo, db: db}
}

func (s *OrderServiceImpl) CreateOrder(order *models.Order) error {
	var totalAmount float64

	// Iniciar transacción
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Validar stock y calcular el total
	for i, item := range order.OrderItems {
		// Obtener el producto con la transacción activa
		product, err := s.productRepo.GetByID(item.ProductID, tx)
		if err != nil {
			log.Printf("Error al buscar producto ID %d: %v", item.ProductID, err)
			tx.Rollback()
			return errors.New("producto no encontrado")
		}

		// Verificar stock disponible
		if product.Stock < item.Quantity {
			log.Printf("Stock insuficiente para el producto ID %d", product.ID)
			tx.Rollback()
			return errors.New("stock insuficiente para un producto")
		}

		// Calcular subtotal
		item.Subtotal = float64(item.Quantity) * product.Price
		totalAmount += item.Subtotal

		// Reducir stock del producto y actualizar en la BD con la transacción activa
		product.Stock -= item.Quantity
		if err := s.productRepo.UpdateStock(product.ID, product.Stock, tx); err != nil {
			log.Printf("Error al actualizar stock del producto ID %d: %v", product.ID, err)
			tx.Rollback()
			return errors.New("error al actualizar stock")
		}

		// Actualizar el pedido con el subtotal corregido
		order.OrderItems[i] = item
	}

	// Asignar total a la orden
	order.TotalAmount = totalAmount

	// Guardar la orden dentro de la transacción
	if err := s.repo.Create(order, tx); err != nil {
		log.Printf("Error al guardar la orden: %v", err)
		tx.Rollback()
		return errors.New("error al crear la orden")
	}

	// Commit si todo fue exitoso
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error al confirmar la transacción: %v", err)
		return errors.New("error al confirmar la transacción")
	}

	log.Println("Orden creada con éxito")
	return nil
}

// GetOrderById busca una orden por su ID.
func (s *OrderServiceImpl) GetOrderById(id uint) (*models.Order, error) {
	return s.repo.FindByID(id)
}
