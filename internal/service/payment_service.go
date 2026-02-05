package service

import (
	"errors"
	"freepass-2026/entity"
	"freepass-2026/internal/repository"
	"freepass-2026/model"
	"time"

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
	menuRepository    repository.IMenuRepository
	orderService      IOrderService
}

func NewPaymentService(db *gorm.DB, paymentRepository repository.IPaymentRepository, orderRepository repository.IOrderRepository, menuRepository repository.IMenuRepository, orderService IOrderService) IPaymentService {
	return &PaymentService{
		db:                db,
		paymentRepository: paymentRepository,
		orderRepository:   orderRepository,
		menuRepository:    menuRepository,
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

	if order.PaymentStatus == "paid" {
		return nil, errors.New("order already paid")
	}

	// validate stock
	orderItems, err := p.orderRepository.GetOrderItems(tx, param.OrderID)
	if err != nil {
		return nil, errors.New("failed to get order items")
	}

	for _, item := range orderItems {
		menu, err := p.menuRepository.GetMenuByID(tx, item.MenuID)
		if err != nil {
			return nil, errors.New("menu not found: " + item.MenuID.String())
		}

		if menu.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for menu: " + menu.Name)
		}
	}

	// reduce stock after validation
	for _, item := range orderItems {
		menu, _ := p.menuRepository.GetMenuByID(tx, item.MenuID)

		menu.Stock -= item.Quantity
		if menu.Stock == 0 {
			menu.IsAvailable = false
		}

		err = p.menuRepository.UpdateMenu(tx, menu)
		if err != nil {
			return nil, err
		}
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
		PaidAt:        time.Now(),
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
