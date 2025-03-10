package repositories

import (
	"order_management/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// productRepositoryImpl es la implementación concreta de ProductRepository
type ProductRepositoryImpl struct {
	db *gorm.DB
}

// NewProductRepository crea una nueva instancia de ProductRepositoryImpl
func NewProductRepository(db *gorm.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db}
}

// GetAll obtiene todos los productos de la base de datos
func (r *ProductRepositoryImpl) GetAll() ([]models.Product, error) {
	var products []models.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetByID obtiene un producto por su ID
func (r *ProductRepositoryImpl) GetByID(id uint, tx *gorm.DB) (*models.Product, error) {
	var product models.Product
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// Update actualiza un producto existente en la base de datos
func (r *ProductRepositoryImpl) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// UpdateStock actualiza el stock de un producto.
func (r *ProductRepositoryImpl) UpdateStock(id uint, newStock int, tx *gorm.DB) error {
	return tx.Model(&models.Product{}).
		Where("id = ?", id).
		Update("stock", newStock).Error
}
