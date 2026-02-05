package service

import (
	"errors"
	"freepass-2026/entity"
	"freepass-2026/internal/repository"
	"freepass-2026/model"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IOrderService interface {
	CreateOrder(userID uuid.UUID, param model.CreateOrderParam) (*model.CreateOrderResponse, error)
	GetUserOrders(userID uuid.UUID) (*model.UserOrderListResponse, error)
	GetOrderDetail(userID uuid.UUID, orderID uuid.UUID) (*model.OrderDetailResponse, error)
	GetCanteenOrders(ownerID uuid.UUID, canteenID uuid.UUID) ([]entity.Order, error)
	UpdateOrderStatus(ownerID uuid.UUID, orderID uuid.UUID, newStatus string) (*model.OrderStatusResponse, error)
	GetOwnerAllOrders(ownerID uuid.UUID, queryParam model.OwnerOrderQueryParam) (*model.OwnerOrderListResponse, error)
	CleanupCanceledOrders() error
}

type OrderService struct {
	db                *gorm.DB
	orderRepository   repository.IOrderRepository
	menuRepository    repository.IMenuRepository
	canteenRepository repository.ICanteenRepository
	userRepository    repository.IUserRepository
}

func NewOrderService(orderRepository repository.IOrderRepository, menuRepository repository.IMenuRepository, canteenRepository repository.ICanteenRepository, userRepository repository.IUserRepository, db *gorm.DB) IOrderService {
	return &OrderService{
		db:                db,
		orderRepository:   orderRepository,
		menuRepository:    menuRepository,
		canteenRepository: canteenRepository,
		userRepository:    userRepository,
	}
}

func (o *OrderService) CreateOrder(userID uuid.UUID, param model.CreateOrderParam) (*model.CreateOrderResponse, error) {
	tx := o.db.Begin()
	defer tx.Rollback()

	canteen, err := o.canteenRepository.GetCanteenByID(tx, param.CanteenID)
	if err != nil {
		return nil, errors.New("canteen not found")
	}

	if !canteen.IsOpen {
		return nil, errors.New("canteen is closed")
	}

	orderID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	var totalPrice int
	var orderItems []entity.OrderItem
	var itemResponses []model.OrderItemResponse

	for _, item := range param.Items {
		// get menu by id
		menu, err := o.menuRepository.GetMenuByID(tx, item.MenuID)
		if err != nil {
			return nil, errors.New("menu not found")
		}

		if menu.CanteenID != param.CanteenID {
			return nil, errors.New("menu does not belong to this canteen")
		}

		if !menu.IsAvailable {
			return nil, errors.New("menu is not available: " + menu.Name)
		}

		if menu.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for menu: " + menu.Name)
		}

		// calculate subtotal
		subtotal := menu.Price * item.Quantity
		totalPrice += subtotal

		// create order item
		orderItemID, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}

		orderItem := entity.OrderItem{
			OrderItemID: orderItemID,
			MenuID:      menu.MenuID,
			OrderID:     orderID,
			Quantity:    item.Quantity,
			Price:       menu.Price,
			Subtotal:    subtotal,
		}

		orderItems = append(orderItems, orderItem)

		itemResponses = append(itemResponses, model.OrderItemResponse{
			MenuID:   item.MenuID,
			Quantity: item.Quantity,
			Price:    menu.Price,
			Subtotal: subtotal,
		})

		if menu.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for menu: " + menu.Name)
		}
	}

	// create order
	order := &entity.Order{
		OrderID:       orderID,
		UserID:        userID,
		CanteenID:     param.CanteenID,
		TotalPrice:    totalPrice,
		PaymentStatus: "unpaid",
		OrderStatus:   "pending",
	}

	err = o.orderRepository.CreateOrder(tx, order)
	if err != nil {
		return nil, err
	}

	// create order items
	for _, orderItem := range orderItems {
		err = o.orderRepository.CreateOrderItem(tx, &orderItem)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	// go rout : timer for autocancel
	go o.startOrderTimer(orderID)
	paymentDeadline := time.Now().Add(15 * time.Minute)
	countdownMinutes := 15

	response := &model.CreateOrderResponse{
		OrderID:          order.OrderID,
		CanteenID:        order.CanteenID,
		TotalPrice:       order.TotalPrice,
		PaymentStatus:    order.PaymentStatus,
		OrderStatus:      order.OrderStatus,
		Items:            itemResponses,
		CreatedAt:        order.CreatedAt,
		PaymentDeadline:  paymentDeadline,
		CountdownMinutes: countdownMinutes,
		Message:          "Please complete payment within 15 minutes or your order will be automatically cancelled",
	}

	return response, nil
}

func (o *OrderService) GetUserOrders(userID uuid.UUID) (*model.UserOrderListResponse, error) {
	tx := o.db.Begin()
	defer tx.Rollback()

	orders, err := o.orderRepository.GetUserOrders(tx, userID)
	if err != nil {
		return nil, err
	}

	var orderResponses []model.UserOrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, model.UserOrderResponse{
			OrderID:       order.OrderID,
			CanteenID:     order.CanteenID,
			TotalPrice:    order.TotalPrice,
			PaymentStatus: order.PaymentStatus,
			OrderStatus:   order.OrderStatus,
			CreatedAt:     order.CreatedAt,
			UpdatedAt:     order.UpdatedAt,
		})
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.UserOrderListResponse{
		Orders: orderResponses,
		Total:  len(orderResponses),
	}

	return response, nil
}

func (o *OrderService) GetOrderDetail(userID uuid.UUID, orderID uuid.UUID) (*model.OrderDetailResponse, error) {
	tx := o.db.Begin()
	defer tx.Rollback()

	order, err := o.orderRepository.GetOrderByID(tx, orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.UserID != userID {
		return nil, errors.New("access denied: order does not belong to this user")
	}

	var itemResponses []model.OrderItemResponse
	for _, orderItem := range order.OrderItems {
		itemResponses = append(itemResponses, model.OrderItemResponse{
			MenuID:   orderItem.MenuID,
			Quantity: orderItem.Quantity,
			Price:    orderItem.Price,
			Subtotal: orderItem.Subtotal,
		})
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.OrderDetailResponse{
		OrderID:       order.OrderID,
		CanteenID:     order.CanteenID,
		TotalPrice:    order.TotalPrice,
		PaymentStatus: order.PaymentStatus,
		OrderStatus:   order.OrderStatus,
		Items:         itemResponses,
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
	}

	return response, nil
}

func (o *OrderService) GetCanteenOrders(ownerID uuid.UUID, canteenID uuid.UUID) ([]entity.Order, error) {
	tx := o.db.Begin()
	defer tx.Rollback()

	canteen, err := o.canteenRepository.GetCanteenByID(tx, canteenID)
	if err != nil {
		return nil, errors.New("canteen not found")
	}

	if canteen.OwnerID != ownerID {
		return nil, errors.New("access denied: canteen not owned by this owner")
	}

	orders, err := o.orderRepository.GetCanteenOrders(tx, canteenID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

/*
* OWNER ONLY
 */

func (o *OrderService) UpdateOrderStatus(ownerID uuid.UUID, orderID uuid.UUID, newStatus string) (*model.OrderStatusResponse, error) {
	tx := o.db.Begin()
	defer tx.Rollback()

	order, err := o.orderRepository.GetOrderByID(tx, orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	canteen, err := o.canteenRepository.GetCanteenByID(tx, order.CanteenID)
	if err != nil {
		return nil, errors.New("canteen not found")
	}

	if canteen.OwnerID != ownerID {
		return nil, errors.New("access denied: canteen not owned by this owner")
	}

	if order.PaymentStatus != "paid" {
		return nil, errors.New("cannot update order status: payment not completed")
	}

	validStatuses := map[string]bool{
		"pending":   true,
		"cooking":   true,
		"ready":     true,
		"completed": true,
	}

	if !validStatuses[newStatus] {
		return nil, errors.New("invalid order status")
	}

	order.OrderStatus = newStatus

	err = o.orderRepository.UpdateOrder(tx, order)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.OrderStatusResponse{
		OrderID:       order.OrderID,
		PaymentStatus: order.PaymentStatus,
		OrderStatus:   order.OrderStatus,
		UpdatedAt:     order.UpdatedAt,
	}

	return response, nil
}

// timer start 30 minute for user for auto cancel order if not paid
func (o *OrderService) startOrderTimer(orderID uuid.UUID) {
	time.Sleep(15 * time.Minute)

	o.autoCancelOrder(orderID)
}

// auto cancel order
func (o *OrderService) autoCancelOrder(orderID uuid.UUID) {
	tx := o.db.Begin()
	defer tx.Rollback()

	order, err := o.orderRepository.GetOrderByID(tx, orderID)
	if err != nil {
		return // not found, skip
	}

	if order.PaymentStatus != "unpaid" {
		return // already paid, skip
	}

	order.PaymentStatus = "canceled"
	order.OrderStatus = "canceled"

	err = o.orderRepository.UpdateOrder(tx, order)
	if err != nil {
		return
	}

	tx.Commit()
}

func (o *OrderService) GetOwnerAllOrders(ownerID uuid.UUID, queryParam model.OwnerOrderQueryParam) (*model.OwnerOrderListResponse, error) {
	tx := o.db.Begin()
	defer tx.Rollback()

	// Set default values
	if queryParam.Limit <= 0 {
		queryParam.Limit = 50
	}
	if queryParam.Page <= 0 {
		queryParam.Page = 1
	}
	if queryParam.Status == "" {
		queryParam.Status = "all"
	}

	offset := (queryParam.Page - 1) * queryParam.Limit

	// Get orders from all owner canteens
	orders, err := o.orderRepository.GetOwnerAllOrders(tx, ownerID, queryParam.Status, queryParam.Limit, offset)
	if err != nil {
		return nil, err
	}

	// Get total count
	totalCount, err := o.orderRepository.CountOwnerAllOrders(tx, ownerID, queryParam.Status)
	if err != nil {
		return nil, err
	}

	// Transform to response with canteen name
	var orderResponses []model.OwnerOrderResponse
	for _, order := range orders {
		var itemResponses []model.OwnerOrderItemResponse
		for _, item := range order.OrderItems {
			menu, _ := o.menuRepository.GetMenuByID(tx, item.MenuID)
			menuName := ""
			if menu != nil {
				menuName = menu.Name
			}

			itemResponses = append(itemResponses, model.OwnerOrderItemResponse{
				MenuID:   item.MenuID,
				MenuName: menuName,
				Quantity: item.Quantity,
				Price:    item.Price,
				Subtotal: item.Subtotal,
			})
		}

		// Get canteen name
		canteen, _ := o.canteenRepository.GetCanteenByID(tx, order.CanteenID)
		canteenName := ""
		if canteen != nil {
			canteenName = canteen.Name
		}

		user, _ := o.userRepository.GetUserByID(tx, order.UserID)
		userName := ""
		userEmail := ""
		if user != nil {
			userName = user.Name
			userEmail = user.Email
		}

		orderResponses = append(orderResponses, model.OwnerOrderResponse{
			OrderID:       order.OrderID,
			UserID:        order.UserID,
			UserName:      userName,
			UserEmail:     userEmail,
			CanteenID:     order.CanteenID,
			CanteenName:   canteenName,
			TotalPrice:    order.TotalPrice,
			PaymentStatus: order.PaymentStatus,
			OrderStatus:   order.OrderStatus,
			Items:         itemResponses,
			CreatedAt:     order.CreatedAt,
			UpdatedAt:     order.UpdatedAt,
		})
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.OwnerOrderListResponse{
		Orders:  orderResponses,
		Total:   totalCount,
		OwnerID: ownerID,
	}

	return response, nil
}

func (o *OrderService) CleanupCanceledOrders() error {
	tx := o.db.Begin()
	defer tx.Rollback()

	canceledOrders, err := o.orderRepository.GetCanceledOrders(tx)
	if err != nil {
		return errors.New("failed to get canceled orders")
	}

	if len(canceledOrders) == 0 {
		log.Println("No canceled orders to cleanup")
		return nil // no canceled orders to clean up
	}

	log.Printf("Found %d canceled orders to delete\n", len(canceledOrders))

	for _, order := range canceledOrders {
		// delete order items
		if err := o.orderRepository.DeleteOrderItems(tx, order.OrderID); err != nil {
			tx.Rollback()
			return errors.New("failed to delete order items: " + err.Error())
		}
		// delete order
		if err := o.orderRepository.DeleteOrder(tx, order.OrderID); err != nil {
			tx.Rollback()
			return errors.New("failed to delete order: " + err.Error())
		}

		log.Printf("Deleted canceled order: %s\n", order.OrderID.String())
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}
