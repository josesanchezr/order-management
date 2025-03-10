package mappers

import (
	"order_management/internal/dtos"
	"order_management/internal/models"
)

func ConvertOrderRequestDTOToOrder(orderRequestDTO dtos.OrderRequestDTO) models.Order {
	order := models.Order{
		CustomerName: orderRequestDTO.CustomerName,
		TotalAmount:  0, // Se calculará después
		OrderItems:   make([]models.OrderItem, len(orderRequestDTO.Items)),
	}

	// Convertir los items
	for i, item := range orderRequestDTO.Items {
		order.OrderItems[i] = models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Subtotal:  0, // Se calculará después
		}
	}

	return order
}

func ConvertOrderToOrderResponseDTO(order models.Order) dtos.OrderResponseDTO {
	orderDTO := dtos.OrderResponseDTO{
		ID:           order.ID,
		CustomerName: order.CustomerName,
		TotalAmount:  order.TotalAmount,
		Items:        make([]dtos.OrderItemResponseDTO, len(order.OrderItems)),
	}

	// Convertir los items
	for i, item := range order.OrderItems {
		orderDTO.Items[i] = dtos.OrderItemResponseDTO{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			Quantity:    item.Quantity,
			Subtotal:    item.Subtotal,
		}
	}

	return orderDTO
}
