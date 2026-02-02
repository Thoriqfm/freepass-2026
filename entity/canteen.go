package entity

import (
	"time"

	"github.com/google/uuid"
)

type Canteen struct {
	CanteenID uuid.UUID `json:"canteen_id" gorm:"type:varchar(36);primary_key"`
	OwnerID   uuid.UUID `json:"owner_id" gorm:"type:varchar(36);not null"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Location  string    `json:"location" gorm:"type:varchar(255);not null"`
	IsOpen    bool      `json:"is_open" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Menus     []Menu     `json:"menus" gorm:"foreignKey:CanteenID"`
	Orders    []Order    `json:"orders" gorm:"foreignKey:CanteenID"`
	Feedbacks []Feedback `json:"feedbacks" gorm:"foreignKey:CanteenID"`
}
