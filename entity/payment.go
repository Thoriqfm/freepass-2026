package entity

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	PaymentID     uuid.UUID `json:"payment_id" gorm:"type:varchar(36);primary_key"`
	OrderID       uuid.UUID `json:"order_id" gorm:"type:varchar(36);not null"`
	Amount        int       `json:"amount" gorm:"not null"`
	PaymentMethod string    `json:"payment_method" gorm:"type:varchar(50);not null"`
	Status        string    `json:"status" gorm:"type:varchar(50);not null"`
	PaidAt        time.Time `json:"paid_at" gorm:"autoCreateTime"`
}
