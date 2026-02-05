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
	GetOrderItems(tx *gorm.DB, orderID uuid.UUID) ([]entity.OrderItem, error)
	GetUserOrders(tx *gorm.DB, userID uuid.UUID) ([]entity.Order, error)
	GetCanteenOrders(tx *gorm.DB, canteenID uuid.UUID) ([]entity.Order, error)
	UpdateOrder(tx *gorm.DB, order *entity.Order) error
	GetOwnerAllOrders(tx *gorm.DB, ownerID uuid.UUID, status string, limit, offset int) ([]entity.Order, error)
	CountOwnerAllOrders(tx *gorm.DB, ownerID uuid.UUID, status string) (int, error)
	GetCanceledOrders(tx *gorm.DB) ([]entity.Order, error)
	DeleteOrder(tx *gorm.DB, orderID uuid.UUID) error
	DeleteOrderItems(tx *gorm.DB, orderID uuid.UUID) error
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

func (r *OrderRepository) GetOrderItems(tx *gorm.DB, orderID uuid.UUID) ([]entity.OrderItem, error) {
	var orderItems []entity.OrderItem
	err := tx.Debug().Where("order_id = ?", orderID).Find(&orderItems).Error
	if err != nil {
		return nil, err
	}
	return orderItems, nil
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

func (r *OrderRepository) GetOwnerAllOrders(tx *gorm.DB, ownerID uuid.UUID, status string, limit, offset int) ([]entity.Order, error) {
	var orders []entity.Order

	query := tx.Debug().
		Table("orders").
		Select("orders.*").
		Joins("JOIN canteens ON orders.canteen_id = canteens.canteen_id").
		Where("canteens.owner_id = ?", ownerID).
		Preload("OrderItems")

	if status != "all" && status != "" {
		query = query.Where("orders.order_status = ?", status)
	}

	err := query.
		Order("orders.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) CountOwnerAllOrders(tx *gorm.DB, ownerID uuid.UUID, status string) (int, error) {
	var count int64

	query := tx.Debug().
		Table("orders").
		Joins("JOIN canteens ON orders.canteen_id = canteens.canteen_id").
		Where("canteens.owner_id = ?", ownerID)

	if status != "all" && status != "" {
		query = query.Where("orders.order_status = ?", status)
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *OrderRepository) GetCanceledOrders(tx *gorm.DB) ([]entity.Order, error) {
	var orders []entity.Order
	if err := tx.Where("order_status = ?", "canceled").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) DeleteOrder(tx *gorm.DB, orderID uuid.UUID) error {
	if err := tx.Where("order_id = ?", orderID).Delete(&entity.Order{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) DeleteOrderItems(tx *gorm.DB, orderID uuid.UUID) error {
	if err := tx.Where("order_id = ?", orderID).Delete(&entity.OrderItem{}).Error; err != nil {
		return err
	}
	return nil
}
