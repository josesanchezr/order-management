package models

import (
	"time"
)

// Order representa un pedido realizado por un cliente.
type Order struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerName string    `gorm:"type:varchar(255);not null" json:"customer_name"`
	TotalAmount  float64   `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relación con OrderItems
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_items"`
}

func (Order) TableName() string {
	return "orders"
}
