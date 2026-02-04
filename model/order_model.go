package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderItemParam struct {
	MenuID   uuid.UUID `json:"menu_id" binding:"required"`
	Quantity int       `json:"quantity" binding:"required,min=1"`
}

type CreateOrderParam struct {
	CanteenID uuid.UUID        `json:"canteen_id" binding:"required"`
	Items     []OrderItemParam `json:"items" binding:"required,min=1"`
}

type OrderItemResponse struct {
	MenuID   uuid.UUID `json:"menu_id"`
	Quantity int       `json:"quantity"`
	Price    int       `json:"price"`
	Subtotal int       `json:"subtotal"`
}

type CreateOrderResponse struct {
	OrderID       uuid.UUID           `json:"order_id"`
	CanteenID     uuid.UUID           `json:"canteen_id"`
	TotalPrice    int                 `json:"total_price"`
	PaymentStatus string              `json:"payment_status"`
	OrderStatus   string              `json:"order_status"`
	Items         []OrderItemResponse `json:"items"`
	CreatedAt     time.Time           `json:"created_at"`
}

type OrderDetailResponse struct {
	OrderID       uuid.UUID           `json:"order_id"`
	CanteenID     uuid.UUID           `json:"canteen_id"`
	TotalPrice    int                 `json:"total_price"`
	PaymentStatus string              `json:"payment_status"`
	OrderStatus   string              `json:"order_status"`
	Items         []OrderItemResponse `json:"items"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}

type UserOrderResponse struct {
	OrderID       uuid.UUID           `json:"order_id"`
	CanteenID     uuid.UUID           `json:"canteen_id"`
	TotalPrice    int                 `json:"total_price"`
	PaymentStatus string              `json:"payment_status"`
	OrderStatus   string              `json:"order_status"`
	Items         []OrderItemResponse `json:"items"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}

type UserOrderListResponse struct {
	Orders []UserOrderResponse `json:"orders"`
	Total  int                 `json:"total"`
}

type OrderStatusResponse struct {
	OrderID       uuid.UUID `json:"order_id"`
	PaymentStatus string    `json:"payment_status"`
	OrderStatus   string    `json:"order_status"`
	UpdatedAt     time.Time `json:"updated_at"`
}
