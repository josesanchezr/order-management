package dtos

// OrderResponseDTO representa la respuesta que se envia al cliente con la información de la orden
type OrderResponseDTO struct {
	ID           uint                   `json:"id"`
	CustomerName string                 `json:"customer_name"`
	TotalAmount  float64                `json:"total_amount"`
	Items        []OrderItemResponseDTO `json:"items"`
}

// OrderItemResponseDTO representa los items relacionados con la orden que se le envia al cliente
type OrderItemResponseDTO struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}
