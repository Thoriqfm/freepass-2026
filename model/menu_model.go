package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateMenuParam struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       int    `json:"price" binding:"required,min=1000"`
	Stock       int    `json:"stock" binding:"required,min=0"`
	IsAvailable bool   `json:"is_available"`
}

type MenuResponse struct {
	MenuID      uuid.UUID `json:"menu_id"`
	CanteenID   uuid.UUID `json:"canteen_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Stock       int       `json:"stock"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MenuListResponse struct {
	Menus []MenuResponse `json:"menus"`
	Total int            `json:"total"`
}

type UpdateMenuParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	IsAvailable bool   `json:"is_available"`
}

type UpdateMenuResponse struct {
	MenuID      uuid.UUID `json:"menu_id"`
	CanteenID   uuid.UUID `json:"canteen_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Stock       int       `json:"stock"`
	IsAvailable bool      `json:"is_available"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DeleteMenuResponse struct {
	Message string    `json:"message"`
	MenuID  uuid.UUID `json:"menu_id"`
}
