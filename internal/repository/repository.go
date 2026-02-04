package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepository    IUserRepository
	MenuRepository    IMenuRepository
	CanteenRepository ICanteenRepository
	OrderRepository   IOrderRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository:    NewUserRepository(db),
		MenuRepository:    NewMenuRepository(db),
		CanteenRepository: NewCanteenRepository(db),
		OrderRepository:   NewOrderRepository(db),
	}
}
