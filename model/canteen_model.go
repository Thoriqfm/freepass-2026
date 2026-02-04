package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateCanteenParam struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
	IsOpen   bool   `json:"is_open"`
}

type CreateCanteenResponse struct {
	CanteenID uuid.UUID `json:"canteen_id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	IsOpen    bool      `json:"is_open"`
	CreatedAt time.Time `json:"created_at"`
}

type CanteenListResponse struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	IsOpen   bool   `json:"is_open"`
}

type UpdateCantenStatusParam struct {
	IsOpen bool `json:"is_open"`
}

type UpdateCanteenStatusResponse struct {
	CanteenID uuid.UUID `json:"canteen_id"`
	Name      string    `json:"name"`
	IsOpen    bool      `json:"is_open"`
	UpdatedAt time.Time `json:"updated_at"`
}
