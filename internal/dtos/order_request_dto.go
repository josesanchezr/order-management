package dtos

// OrderRequestDTO representa el payload recibido para crear una orden
type OrderRequestDTO struct {
	CustomerName string                `json:"customer_name" validate:"required"`
	Items        []OrderItemRequestDTO `json:"items" validate:"required,dive"`
}

// OrderItemRequestDTO representa los items relacionados a la orden
type OrderItemRequestDTO struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,gt=0"`
}
