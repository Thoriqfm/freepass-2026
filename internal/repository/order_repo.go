package repository

import (
	"freepass-2026/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	CreateOrder(tx *gorm.DB, order *entity.Order) error
	CreateOrderItem(tx *gorm.DB, orderItem *entity.OrderItem) error
	GetOrderByID(tx *gorm.DB, orderID uuid.UUID) (*entity.Order, error)
	GetUserOrders(tx *gorm.DB, userID uuid.UUID) ([]entity.Order, error)
	GetCanteenOrders(tx *gorm.DB, canteenID uuid.UUID) ([]entity.Order, error)
	UpdateOrder(tx *gorm.DB, order *entity.Order) error
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(tx *gorm.DB, order *entity.Order) error {
	err := tx.Debug().Create(&order).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) CreateOrderItem(tx *gorm.DB, orderItem *entity.OrderItem) error {
	err := tx.Debug().Create(&orderItem).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetOrderByID(tx *gorm.DB, orderID uuid.UUID) (*entity.Order, error) {
	var order entity.Order
	err := tx.Debug().Preload("OrderItems").Where("order_id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetUserOrders(tx *gorm.DB, userID uuid.UUID) ([]entity.Order, error) {
	var orders []entity.Order
	err := tx.Debug().Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetCanteenOrders(tx *gorm.DB, canteenID uuid.UUID) ([]entity.Order, error) {
	var orders []entity.Order
	err := tx.Debug().Where("canteen_id = ?", canteenID).Preload("OrderItems").Order("created_at DESC").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) UpdateOrder(tx *gorm.DB, order *entity.Order) error {
	err := tx.Debug().Save(&order).Error
	if err != nil {
		return err
	}
	return nil
}
