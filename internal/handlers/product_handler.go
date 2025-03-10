package handlers

import (
	"net/http"
	"strconv"

	"order_management/internal/dtos"
	"order_management/internal/mappers"
	"order_management/internal/ports"

	"github.com/labstack/echo/v4"
)

// ProductHandler maneja las solicitudes HTTP relacionadas con productos
type ProductHandler struct {
	productService ports.ProductService
}

// NewProductHandler registra los endpoints de productos en Echo
func NewProductHandler(apiGroup *echo.Group, productService ports.ProductService) {
	handler := &ProductHandler{productService: productService}

	apiGroup.GET("/products", handler.GetAllProducts)
	//apiGroup.GET("/products/:id", handler.GetProductByID)
	apiGroup.PUT("/products/:id/stock", handler.UpdateStock)
}

// GetAllProducts maneja la solicitud para obtener todos los productos
func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error al obtener productos"})
	}

	productResponseDTOs := mappers.ConvertProductToProductResponseDTO(products)

	return c.JSON(http.StatusOK, productResponseDTOs)
}

// UpdateStock maneja la solicitud para actualizar el stock de un producto
func (h *ProductHandler) UpdateStock(c echo.Context) error {
	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ID inválido"})
	}

	var stockDTO dtos.UpdateStocRequestkDTO
	if err := c.Bind(&stockDTO); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Datos de entrada inválidos"})
	}

	id := uint(idInt) // Conversión segura de int a uint

	if err := h.productService.UpdateStock(id, stockDTO.Stock); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error al actualizar stock"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Stock actualizado correctamente"})
}
