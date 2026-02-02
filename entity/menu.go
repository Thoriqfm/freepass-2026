package entity

import (
	"time"

	"github.com/google/uuid"
)

type Menu struct {
	MenuID      uuid.UUID `json:"menu_id" gorm:"type:varchar(36);primary_key"`
	CanteenID   uuid.UUID `json:"canteen_id" gorm:"type:varchar(36);not null"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text"`
	Price       int       `json:"price" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"not null"`
	IsAvailable bool      `json:"is_available" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:MenuID"`

	// CanteenID uuid.UUID `json:"canteen_id" gorm:"type:varchar(36);not null;foreignKey:CanteenID;references:CanteenID"`
}
