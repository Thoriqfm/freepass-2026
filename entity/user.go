package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `json:"user_id" gorm:"type:varchar(36);primary_key"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Email     string    `json:"email" gorm:"type:varchar(255);unique;not null"`
	Password  string    `json:"password" gorm:"type:varchar(255);not null"`
	RoleID    int       `json:"role_id"`
	Phone     string    `json:"phone" gorm:"type:varchar(255)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	CanteenOwner []Canteen  `json:"canteen_owner" gorm:"foreignKey:OwnerID"`
	Feedbacks    []Feedback `json:"feedbacks" gorm:"foreignKey:UserID"`
}
