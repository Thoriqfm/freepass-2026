package service

import (
	"errors"
	"freepass-2026/entity"
	"freepass-2026/internal/repository"
	"freepass-2026/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IFeedbackService interface {
	CreateFeedback(userID uuid.UUID, param model.CreateFeedbackParam) (*model.FeedbackResponse, error)
}

type FeedbackService struct {
	db                 *gorm.DB
	feedbackRepository repository.IFeedbackRepository
	orderRepository    repository.IOrderRepository
	canteenRepository  repository.ICanteenRepository
}

func NewFeedbackService(feedbackRepository repository.IFeedbackRepository, orderRepository repository.IOrderRepository, canteenRepository repository.ICanteenRepository, db *gorm.DB) IFeedbackService {
	return &FeedbackService{
		db:                 db,
		feedbackRepository: feedbackRepository,
		orderRepository:    orderRepository,
		canteenRepository:  canteenRepository,
	}
}

func (f *FeedbackService) CreateFeedback(userID uuid.UUID, param model.CreateFeedbackParam) (*model.FeedbackResponse, error) {
	tx := f.db.Begin()
	defer tx.Rollback()

	order, err := f.orderRepository.GetOrderByID(tx, param.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.UserID != userID {
		return nil, errors.New("access denied: order not owned by this user")
	}

	// validate only for completed orders
	if order.OrderStatus != "completed" {
		return nil, errors.New("cannot give feedback for incomplete order")
	}

	if order.PaymentStatus != "paid" {
		return nil, errors.New("cannot give feedback for unpaid order")
	}

	// create
	feedbackID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	feedback := &entity.Feedback{
		FeedbackID: feedbackID,
		UserID:     userID,
		OrderID:    param.OrderID,
		CanteenID:  order.CanteenID,
		Rating:     param.Rating,
		Comment:    param.Comment,
	}

	err = f.feedbackRepository.CreateFeedback(tx, feedback)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.FeedbackResponse{
		FeedbackID: feedback.FeedbackID,
		OrderID:    feedback.OrderID,
		UserID:     feedback.UserID,
		CanteenID:  feedback.CanteenID,
		Rating:     feedback.Rating,
		Comment:    feedback.Comment,
		CreatedAt:  feedback.CreatedAt,
		UpdatedAt:  feedback.UpdatedAt,
	}

	return response, nil

}
