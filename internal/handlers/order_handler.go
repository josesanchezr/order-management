package handlers

import (
	"net/http"
	"order_management/internal/dtos"
	"order_management/internal/mappers"
	"order_management/internal/middlewares"
	"order_management/internal/ports"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderService ports.OrderService
}

func NewOrderHandler(apiGroup *echo.Group, orderService ports.OrderService, redisClient *redis.Client) {
	handler := &OrderHandler{
		orderService: orderService,
	}

	// Rutas
	// Aplicar middleware de idempotencia solo en POST /orders
	apiGroup.POST("/orders", handler.CreateOrder, middlewares.IdempotencyMiddleware(redisClient))
	apiGroup.GET("/orders/:id", handler.GetOrderById)
}

// CreateOrder maneja la creación de un nuevo pedido
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var orderRequest dtos.OrderRequestDTO
	if err := c.Bind(&orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	// Validar la estructura después de Bind()
	if err := c.Validate(orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Convertir DTO a modelo
	order := mappers.ConvertOrderRequestDTOToOrder(orderRequest)

	// Llamar al servicio para crear la orden
	err := h.orderService.CreateOrder(&order)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "Orden creada correctamente"})
}

// GetOrderById maneja la obtención de una orden por su ID
func (h *OrderHandler) GetOrderById(c echo.Context) error {
	orderIDInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ID inválido"})
	}

	orderID := uint(orderIDInt) // Conversión segura de int a uint

	order, err := h.orderService.GetOrderById(orderID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found"})
	}

	// Convertir model a DTO
	orderDTO := mappers.ConvertOrderToOrderResponseDTO(*order)

	return c.JSON(http.StatusOK, orderDTO)
}
