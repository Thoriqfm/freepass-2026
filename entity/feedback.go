package entity

import (
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	FeedbackID uuid.UUID `json:"feedback_id" gorm:"type:varchar(36);primary_key"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:varchar(36);not null"`
	CanteenID  uuid.UUID `json:"canteen_id" gorm:"type:varchar(36);not null"`
	OrderID    uuid.UUID `json:"order_id" gorm:"type:varchar(36);not null"`
	Rating     int       `json:"rating" gorm:"not null"`
	Comment    string    `json:"comment" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`

	// OrderID   uuid.UUID `json:"order_id" gorm:"type:varchar(36);not null;foreignKey:OrderID;references:OrderID"`
	// UserID    uuid.UUID `json:"user_id" gorm:"type:varchar(36);not null;foreignKey:UserID;references:UserID"`
	// CanteenID uuid.UUID `json:"canteen_id" gorm:"type:varchar(36);not null;foreignKey:CanteenID;references:CanteenID"`
}
