package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID       uuid.UUID `json:"order_id" gorm:"type:varchar(36);primary_key"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:varchar(36);not null"`
	CanteenID     uuid.UUID `json:"canteen_id" gorm:"type:varchar(36);not null"`
	TotalPrice    int       `json:"total_price" gorm:"not null"`
	PaymentStatus string    `json:"payment_status" gorm:"type:varchar(50);not null"`
	OrderStatus   string    `json:"order_status" gorm:"type:varchar(50);not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
	Payment    []Payment   `json:"payment" gorm:"foreignKey:OrderID"`
	Feedbacks  []Feedback  `json:"feedbacks" gorm:"foreignKey:OrderID"`
}
