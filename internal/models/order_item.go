package models

import "time"

// OrderItem representa los productos dentro de un pedido.
type OrderItem struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   uint      `gorm:"not null" json:"order_id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Subtotal  float64   `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relación con Order
	Order Order `gorm:"foreignKey:OrderID" json:"order"`

	// Relación con Product
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
