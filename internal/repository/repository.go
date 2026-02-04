package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepository     IUserRepository
	MenuRepository     IMenuRepository
	CanteenRepository  ICanteenRepository
	OrderRepository    IOrderRepository
	PaymentRepository  IPaymentRepository
	FeedbackRepository IFeedbackRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository:     NewUserRepository(db),
		MenuRepository:     NewMenuRepository(db),
		CanteenRepository:  NewCanteenRepository(db),
		OrderRepository:    NewOrderRepository(db),
		PaymentRepository:  NewPaymentRepository(db),
		FeedbackRepository: NewFeedbackRepository(db),
	}
}
