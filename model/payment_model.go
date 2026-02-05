package model

import (
	"time"

	"github.com/google/uuid"
)

type CreatePaymentParam struct {
	OrderID       uuid.UUID `json:"order_id,omitempty"` // omitempty karena akan diisi dari URL
	PaymentMethod string    `json:"payment_method" binding:"required"`
}

type CreatePaymentResponse struct {
	PaymentID     uuid.UUID `json:"payment_id"`
	OrderID       uuid.UUID `json:"order_id"`
	Amount        int       `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	PaidAt        time.Time `json:"paid_at"`
}

type PaymentListResponse struct {
	Payments []CreatePaymentResponse `json:"payments"`
	Total    int                     `json:"total"`
}

// /* For Owner: verify payment */
// type VerifyPaymentParam struct {
// 	PaymentID uuid.UUID `json:"payment_id" binding:"required"`
// }

// type VerifyPaymentResponse struct {
// 	PaymentID uuid.UUID `json:"payment_id"`
// 	OrderID   uuid.UUID `json:"order_id"`
// 	Status    string    `json:"status"`
// 	Message   string    `json:"message"`
// 	UpdatedAt time.Time `json:"updated_at"`
// }
