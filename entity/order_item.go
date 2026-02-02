package entity

import "github.com/google/uuid"

type OrderItem struct {
	OrderItemID uuid.UUID `json:"order_item_id" gorm:"type:varchar(36);primary_key"`
	MenuID      uuid.UUID `json:"menu_id" gorm:"type:varchar(36);not null"`
	OrderID     uuid.UUID `json:"order_id" gorm:"type:varchar(36);not null"`
	Quantity    int       `json:"quantity" gorm:"not null"`
	Price       int       `json:"price" gorm:"not null"`
	Subtotal    int       `json:"subtotal" gorm:"not null"`
}
