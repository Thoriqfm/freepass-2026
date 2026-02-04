package repository

import (
	"freepass-2026/entity"

	"gorm.io/gorm"
)

type IPaymentRepository interface {
	CreatePayment(tx *gorm.DB, payment *entity.Payment) error
}

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreatePayment(tx *gorm.DB, payment *entity.Payment) error {
	err := tx.Debug().Create(&payment).Error
	if err != nil {
		return err
	}
	return nil
}
