package dtos

// ProductResponseDTO representa la respuesta que se envía al cliente
type ProductResponseDTO struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

// UpdateStocRequestkDTO representa el payload recibido en la actualización del stock de productos
type UpdateStocRequestkDTO struct {
	Stock int `json:"stock"`
}
