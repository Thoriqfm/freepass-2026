package service

import (
	"errors"
	"freepass-2026/entity"
	"freepass-2026/internal/repository"
	"freepass-2026/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPaymentService interface {
	CreatePayment(userID uuid.UUID, param model.CreatePaymentParam) (*model.CreatePaymentResponse, error)
}

type PaymentService struct {
	db                *gorm.DB
	paymentRepository repository.IPaymentRepository
	orderRepository   repository.IOrderRepository
	orderService      IOrderService
}

func NewPaymentService(db *gorm.DB, paymentRepository repository.IPaymentRepository, orderRepository repository.IOrderRepository, orderService IOrderService) IPaymentService {
	return &PaymentService{
		db:                db,
		paymentRepository: paymentRepository,
		orderRepository:   orderRepository,
		orderService:      orderService,
	}
}

func (p *PaymentService) CreatePayment(userID uuid.UUID, param model.CreatePaymentParam) (*model.CreatePaymentResponse, error) {
	tx := p.db.Begin()
	defer tx.Rollback()

	order, err := p.orderRepository.GetOrderByID(tx, param.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.UserID != userID {
		return nil, errors.New("access denied: order not owned by this user")
	}

	// validate order
	if order.PaymentStatus == "paid" {
		return nil, errors.New("order already paid")
	}

	paymentID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	payment := &entity.Payment{
		PaymentID:     paymentID,
		OrderID:       order.OrderID,
		Amount:        order.TotalPrice,
		PaymentMethod: param.PaymentMethod,
		Status:        "paid",
		PaidAt:        order.CreatedAt,
	}

	err = p.paymentRepository.CreatePayment(tx, payment)
	if err != nil {
		return nil, err
	}

	order.PaymentStatus = "paid"
	err = p.orderRepository.UpdateOrder(tx, order)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.CreatePaymentResponse{
		PaymentID:     payment.PaymentID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		PaymentMethod: payment.PaymentMethod,
		Status:        payment.Status,
		PaidAt:        payment.PaidAt,
	}

	return response, nil
}
