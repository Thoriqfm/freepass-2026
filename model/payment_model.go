package model

import (
	"time"

	"github.com/google/uuid"
)

type CreatePaymentParam struct {
	OrderID       uuid.UUID `json:"order_id" binding:"required"`
	PaymentMethod string    `json:"payment_method" binding:"required,oneof=credit_card bank_transfer e_wallet"`
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
