package integration_test

import (
	"encoding/json"
	"net/http"
	"order_management/internal/dtos"
	"order_management/internal/handlers"
	"order_management/internal/models"
	"order_management/internal/repositories"
	"order_management/internal/services"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// setupOrderRoutes configura las rutas específicas para OrderHandler
func setupOrderRoutes(e *echo.Echo, db *gorm.DB, redisClient *redis.Client) {
	orderRepo := repositories.NewOrderRepository(db)
	productRepo := repositories.NewProductRepository(db)
	orderService := services.NewOrderService(orderRepo, productRepo, db)

	apiGroup := e.Group("/api")
	handlers.NewOrderHandler(apiGroup, orderService, redisClient)
}

// TestCreateOrderSuccess: Creación exitosa de una orden
func TestCreateOrderSuccess(t *testing.T) {
	SetupTestServer(t, setupOrderRoutes)
	defer TearDown()

	// Insertar un producto en la DB antes de la prueba
	product := models.Product{
		Name:  "Producto de prueba",
		Price: 100.0,
		Stock: 10,
	}
	err := db.Create(&product).Error
	assert.NoError(t, err) // Verificar que no hubo error al insertar el producto

	client := resty.New()
	orderRequest := dtos.OrderRequestDTO{
		CustomerName: "Customer 1",
		Items: []dtos.OrderItemRequestDTO{
			{ProductID: product.ID, Quantity: 2}, // Usamos el ID del producto insertado
		},
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(orderRequest).
		Post(server.URL + "/api/orders")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode())

	var responseData map[string]string
	err = json.Unmarshal(resp.Body(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, "Orden creada correctamente", responseData["message"])
}

// TestCreateOrderInvalidPayload: Creación de orden fallida por payload inválido
func TestCreateOrderInvalidPayload(t *testing.T) {
	SetupTestServer(t, setupOrderRoutes)
	defer TearDown()

	client := resty.New()
	invalidOrderRequest := map[string]interface{}{
		"items": []map[string]interface{}{
			{"product_id": 1, "quantity": 2}, // Aquí falta "customer_name"
		},
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(invalidOrderRequest).
		Post(server.URL + "/api/orders")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode())

	var responseData map[string]string
	err = json.Unmarshal(resp.Body(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, "Key: 'OrderRequestDTO.CustomerName' Error:Field validation for 'CustomerName' failed on the 'required' tag", responseData["error"])
}

// TestGetOrderByIdSuccess: Obtener una orden por ID
func TestGetOrderByIdSuccess(t *testing.T) {
	SetupTestServer(t, setupOrderRoutes)
	defer TearDown()

	// Insertar un producto en la DB antes de la prueba
	product := models.Product{
		Name:  "Producto de prueba",
		Price: 100.0,
		Stock: 10,
	}
	err := db.Create(&product).Error
	assert.NoError(t, err) // Verificar que no hubo error al insertar el producto

	// Crear una orden primero
	client := resty.New()
	orderRequest := dtos.OrderRequestDTO{
		CustomerName: "Customer 1",
		Items: []dtos.OrderItemRequestDTO{
			{ProductID: 1, Quantity: 2},
		},
	}

	createResp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(orderRequest).
		Post(server.URL + "/api/orders")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, createResp.StatusCode())

	// Ahora obtener la orden con ID 1
	getResp, err := client.R().
		Get(server.URL + "/api/orders/1")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, getResp.StatusCode())

	var orderResponse dtos.OrderResponseDTO
	err = json.Unmarshal(getResp.Body(), &orderResponse)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), orderResponse.ID)
	assert.Equal(t, 200.0, orderResponse.TotalAmount)
}

// TestGetOrderByIdInvalidID: Obtener una orden con ID inválido
func TestGetOrderByIdInvalidID(t *testing.T) {
	SetupTestServer(t, setupOrderRoutes)
	defer TearDown()

	client := resty.New()

	resp, err := client.R().
		Get(server.URL + "/api/orders/invalidID")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode())

	var responseData map[string]string
	err = json.Unmarshal(resp.Body(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, "ID inválido", responseData["error"])
}

// TestGetOrderByIdNotFound: Obtener una orden que no existe
func TestGetOrderByIdNotFound(t *testing.T) {
	SetupTestServer(t, setupOrderRoutes)
	defer TearDown()

	client := resty.New()

	resp, err := client.R().
		Get(server.URL + "/api/orders/999") // ID que no existe

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode())

	var responseData map[string]string
	err = json.Unmarshal(resp.Body(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, "Order not found", responseData["error"])
}
