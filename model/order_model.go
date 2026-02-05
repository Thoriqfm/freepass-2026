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
	OrderID          uuid.UUID           `json:"order_id"`
	CanteenID        uuid.UUID           `json:"canteen_id"`
	TotalPrice       int                 `json:"total_price"`
	PaymentStatus    string              `json:"payment_status"`
	OrderStatus      string              `json:"order_status"`
	Items            []OrderItemResponse `json:"items"`
	CreatedAt        time.Time           `json:"created_at"`
	PaymentDeadline  time.Time           `json:"payment_deadline"`
	CountdownMinutes int                 `json:"countdown_minutes"`
	Message          string              `json:"message"`
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
	OrderID       uuid.UUID `json:"order_id"`
	CanteenID     uuid.UUID `json:"canteen_id"`
	TotalPrice    int       `json:"total_price"`
	PaymentStatus string    `json:"payment_status"`
	OrderStatus   string    `json:"order_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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

// Response for owner to view all orders from all canteens
type OwnerOrderResponse struct {
	OrderID       uuid.UUID                `json:"order_id"`
	UserID        uuid.UUID                `json:"user_id"`
	UserName      string                   `json:"user_name"`
	UserEmail     string                   `json:"user_email"`
	CanteenID     uuid.UUID                `json:"canteen_id"`
	CanteenName   string                   `json:"canteen_name"`
	TotalPrice    int                      `json:"total_price"`
	PaymentStatus string                   `json:"payment_status"`
	OrderStatus   string                   `json:"order_status"`
	Items         []OwnerOrderItemResponse `json:"items"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
}

type OwnerOrderItemResponse struct {
	MenuID   uuid.UUID `json:"menu_id"`
	MenuName string    `json:"menu_name"`
	Quantity int       `json:"quantity"`
	Price    int       `json:"price"`
	Subtotal int       `json:"subtotal"`
}

type OwnerOrderListResponse struct {
	Orders   []OwnerOrderResponse `json:"orders"`
	Total    int                  `json:"total"`
	OwnerID  uuid.UUID            `json:"owner_id"`
	FilterBy string               `json:"filter_by,omitempty"`
}

// query param for filter
// Query parameters untuk filter orders
type OwnerOrderQueryParam struct {
	Status string `form:"status"` // pending, confirmed, completed, canceled, all
	Limit  int    `form:"limit"`  // default 50
	Page   int    `form:"page"`   // default 1
}
