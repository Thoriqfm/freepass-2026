package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateFeedbackParam struct {
	OrderID uuid.UUID `json:"order_id" binding:"required"`
	Rating  int       `json:"rating" binding:"required,min=1,max=5"`
	Comment string    `json:"comment"`
}

type FeedbackResponse struct {
	FeedbackID uuid.UUID `json:"feedback_id"`
	OrderID    uuid.UUID `json:"order_id"`
	UserID     uuid.UUID `json:"user_id"`
	CanteenID  uuid.UUID `json:"canteen_id"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FeedbackListResponse struct {
	Feedbacks []FeedbackResponse `json:"feedbacks"`
	Total     int                `json:"total"`
}

type CanteenFeedbackResponse struct {
	FeedbackID uuid.UUID `json:"feedback_id"`
	UserID     uuid.UUID `json:"user_id"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
}

type CanteenFeedbackListResponse struct {
	Feedbacks []CanteenFeedbackResponse `json:"feedbacks"`
	Average   float64                   `json:"average_rating"`
	Total     int                       `json:"total"`
}

type DeleteFeedbackResponse struct {
	FeedbackID uuid.UUID `json:"feedback_id"`
	Message    string    `json:"message"`
	DeletedAt  time.Time `json:"deleted_at"`
}
