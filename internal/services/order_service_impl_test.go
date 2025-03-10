package services

import (
	"errors"
	"log"
	"order_management/internal/models"
	"order_management/test/mocks"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)

	// Crear base de datos en memoria para pruebas
	db, errDB := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if errDB != nil {
		log.Fatalf("Error al abrir la base de datos en memoria: %v", errDB)
	}
	db.AutoMigrate(&models.Product{}, &models.Order{}, &models.OrderItem{})

	// Insertar producto de prueba en la base de datos
	product := &models.Product{ID: 1, Name: "Laptop", Price: 500, Stock: 10}
	db.Create(product)

	service := NewOrderService(mockOrderRepo, mockProductRepo, db)

	order := &models.Order{
		ID:          1,
		TotalAmount: 0,
		OrderItems: []models.OrderItem{
			{ProductID: 1, Quantity: 2, Subtotal: 0},
		},
	}

	// Mocks esperados (se usan los métodos del ORM, no los mocks)
	mockProductRepo.EXPECT().GetByID(uint(1), gomock.Any()).Return(product, nil).Times(1)
	mockProductRepo.EXPECT().UpdateStock(uint(1), int(8), gomock.Any()).Return(nil).Times(1)
	mockOrderRepo.EXPECT().Create(order, gomock.Any()).Return(nil).Times(1)

	err := service.CreateOrder(order)

	// **Validaciones**
	assert.NoError(t, err)
	assert.Equal(t, float64(1000), order.TotalAmount) // 2 * 500
}

// Test para CreateOrder con producto no encontrado
func TestCreateOrder_ProductNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Crear base de datos en memoria para pruebas
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)

	service := NewOrderService(mockOrderRepo, mockProductRepo, db)

	order := &models.Order{
		ID:         1,
		OrderItems: []models.OrderItem{{ProductID: 1, Quantity: 2}},
	}

	mockProductRepo.EXPECT().GetByID(uint(1), gomock.Any()).Return(nil, errors.New("not found")).Times(1)

	err := service.CreateOrder(order)

	assert.Error(t, err)
	assert.Equal(t, "producto no encontrado", err.Error())
}

// Test para CreateOrder con stock insuficiente
func TestCreateOrder_InsufficientStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Crear base de datos en memoria para pruebas
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)

	service := NewOrderService(mockOrderRepo, mockProductRepo, db)

	order := &models.Order{
		ID:         1,
		OrderItems: []models.OrderItem{{ProductID: 1, Quantity: 5}},
	}

	product := &models.Product{ID: 1, Name: "Laptop", Price: 500, Stock: 2}

	mockProductRepo.EXPECT().GetByID(uint(1), gomock.Any()).Return(product, nil).Times(1)

	err := service.CreateOrder(order)

	assert.Error(t, err)
	assert.Equal(t, "stock insuficiente para un producto", err.Error())
}

// Test para CreateOrder con error al actualizar stock
func TestCreateOrder_UpdateStockFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Crear base de datos en memoria para pruebas
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)

	orderService := NewOrderService(mockOrderRepo, mockProductRepo, db)

	// Simulación de datos
	product := &models.Product{
		ID:    1,
		Stock: 10,
		Price: 100.0,
	}
	order := &models.Order{
		OrderItems: []models.OrderItem{
			{ProductID: product.ID, Quantity: 5},
		},
	}

	// Mock de la obtención del producto
	mockProductRepo.
		EXPECT().
		GetByID(product.ID, gomock.Any()).
		Return(product, nil)

	// Mock de UpdateStock que falla
	mockProductRepo.
		EXPECT().
		UpdateStock(product.ID, gomock.Any(), gomock.Any()).
		Return(errors.New("error en la base de datos"))

	// Ejecutar la prueba
	err := orderService.CreateOrder(order)

	// Verificar resultado esperado
	assert.Error(t, err)
	assert.Equal(t, "error al actualizar stock", err.Error())
}

// Test para CreateOrder con error al crear la orden
func TestCreateOrder_SaveOrderFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Crear base de datos en memoria para pruebas
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)

	orderService := NewOrderService(mockOrderRepo, mockProductRepo, db)

	// Simulación de datos
	product := &models.Product{
		ID:    1,
		Stock: 10,
		Price: 100.0,
	}
	order := &models.Order{
		OrderItems: []models.OrderItem{
			{ProductID: product.ID, Quantity: 5},
		},
	}

	// Mock de la obtención del producto
	mockProductRepo.
		EXPECT().
		GetByID(product.ID, gomock.Any()).
		Return(product, nil)

	// Mock de actualización de stock con éxito
	mockProductRepo.
		EXPECT().
		UpdateStock(product.ID, gomock.Any(), gomock.Any()).
		Return(nil)

	// Mock de error en la creación de la orden
	mockOrderRepo.
		EXPECT().
		Create(order, gomock.Any()).
		Return(errors.New("error en la base de datos"))

	// Ejecutar la prueba
	err := orderService.CreateOrder(order)

	// Verificar resultado esperado
	assert.Error(t, err)
	assert.Equal(t, "error al crear la orden", err.Error())
}

// Test para GetOrderById con éxito
func TestGetOrderById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Crear base de datos en memoria para pruebas
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)

	service := NewOrderService(mockOrderRepo, nil, db)

	expectedOrder := &models.Order{ID: 1, TotalAmount: 100}

	mockOrderRepo.EXPECT().FindByID(uint(1)).Return(expectedOrder, nil).Times(1)

	order, err := service.GetOrderById(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
}

// Test para GetOrderById cuando la orden no existe
func TestGetOrderById_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Crear base de datos en memoria para pruebas
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)

	service := NewOrderService(mockOrderRepo, nil, db)

	mockOrderRepo.EXPECT().FindByID(uint(1)).Return(nil, errors.New("not found")).Times(1)

	order, err := service.GetOrderById(1)

	assert.Error(t, err)
	assert.Nil(t, order)
}
