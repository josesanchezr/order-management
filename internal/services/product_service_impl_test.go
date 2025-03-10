package services

import (
	"errors"
	"order_management/internal/models"
	"order_management/test/mocks"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestGetAllProducts_Success verifica que GetAllProducts() retorne correctamente la lista de productos.
func TestGetAllProducts_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Crear base de datos en memoria para pruebas
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	productService := NewProductService(mockProductRepo, db)

	// Datos simulados
	expectedProducts := []models.Product{
		{ID: 1, Name: "Producto 1", Stock: 10, Price: 100.0},
		{ID: 2, Name: "Producto 2", Stock: 5, Price: 50.0},
	}

	// Simula la respuesta exitosa del repositorio
	mockProductRepo.
		EXPECT().
		GetAll().
		Return(expectedProducts, nil)

	// Ejecutar
	products, err := productService.GetAllProducts()

	// Verificar
	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
}

// TestGetAllProducts_Failure verifica que GetAllProducts() retorne un error si falla la consulta a la base de datos.
func TestGetAllProducts_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Crear base de datos en memoria para pruebas
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	productService := NewProductService(mockProductRepo, db)

	// Simula un error en la base de datos
	mockProductRepo.
		EXPECT().
		GetAll().
		Return(nil, errors.New("error en base de datos"))

	// Ejecutar
	products, err := productService.GetAllProducts()

	// Verificar
	assert.Error(t, err)
	assert.Nil(t, products)
}

// TestUpdateStock_Success verifica que UpdateStock() actualice correctamente el stock del producto.
func TestUpdateStock_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Crear base de datos en memoria para pruebas
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	productService := NewProductService(mockProductRepo, db)

	// Datos simulados
	productId := uint(1)
	newStock := 5

	// Simula la actualización exitosa del stock
	mockProductRepo.
		EXPECT().
		UpdateStock(uint(1), 5, gomock.Any()).
		Return(nil)

	// Ejecutar
	err := productService.UpdateStock(productId, newStock)

	// Verificar
	assert.NoError(t, err)
	assert.Equal(t, newStock, 5)
}

// TestUpdateStock_UpdateFailure verifica que UpdateStock() retorne un error si la actualización del stock falla.
func TestUpdateStock_UpdateFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Crear base de datos en memoria para pruebas
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	productService := NewProductService(mockProductRepo, db)

	// Datos simulados
	product := &models.Product{ID: 1, Name: "Producto 1", Stock: 10, Price: 100.0}
	newStock := 5

	// Simula un error en la actualización del stock
	mockProductRepo.
		EXPECT().
		UpdateStock(product.ID, newStock, gomock.Any()).
		Return(errors.New("error al actualizar stock"))

	// Ejecutar
	err := productService.UpdateStock(product.ID, newStock)

	// Verificar
	assert.Error(t, err)
	assert.Equal(t, "error al actualizar stock", err.Error())
}
