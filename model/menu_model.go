package model

type CreateMenuParam struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       int    `json:"price" binding:"required,min=1000"`
	Stock       int    `json:"stock" binding:"required,min=0"`
	IsAvailable bool   `json:"is_available"`
}
