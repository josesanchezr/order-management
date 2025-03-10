package services

import (
	"errors"
	"log"
	"order_management/internal/models"
	"order_management/internal/ports"

	"gorm.io/gorm"
)

type ProductServiceImpl struct {
	productRepo ports.ProductRepository
	db          *gorm.DB
}

func NewProductService(productRepo ports.ProductRepository, db *gorm.DB) ports.ProductService {
	return &ProductServiceImpl{productRepo: productRepo, db: db}
}

func (s *ProductServiceImpl) GetAllProducts() ([]models.Product, error) {
	return s.productRepo.GetAll()
}

func (s *ProductServiceImpl) UpdateStock(id uint, stock int) error {
	// Iniciar transacción
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.productRepo.UpdateStock(id, stock, tx); err != nil {
		return err
	}

	// Commit si todo fue exitoso
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error al confirmar la transacción: %v", err)
		return errors.New("error al confirmar la transacción")
	}
	return nil
}
