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

// func (r *PaymentRepository) GetPaymentByID(tx *gorm.DB, paymentID uuid.UUID) (*entity.Payment, error) {
// 	var payment entity.Payment
// 	err := tx.Debug().Where("payment_id = ?", paymentID).First(&payment).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &payment, nil
// }

// func (r *PaymentRepository) GetPaymentByOrderID(tx *gorm.DB, orderID uuid.UUID) (*entity.Payment, error) {
// 	var payment entity.Payment
// 	err := tx.Debug().Where("order_id = ?", orderID).First(&payment).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &payment, nil
// }

// func (r *PaymentRepository) GetUserPayments(tx *gorm.DB, userID uuid.UUID) ([]entity.Payment, error) {
// 	var payments []entity.Payment

// }
