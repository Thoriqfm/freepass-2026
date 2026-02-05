package repository

import (
	"freepass-2026/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IFeedbackRepository interface {
	CreateFeedback(tx *gorm.DB, feedback *entity.Feedback) error
	GetFeedbackByID(tx *gorm.DB, feedbackID uuid.UUID) (*entity.Feedback, error)
	DeleteFeedback(tx *gorm.DB, feedbackID uuid.UUID) error
}

type FeedbackRepository struct {
	db *gorm.DB
}

func NewFeedbackRepository(db *gorm.DB) IFeedbackRepository {
	return &FeedbackRepository{db: db}
}

func (r *FeedbackRepository) CreateFeedback(tx *gorm.DB, feedback *entity.Feedback) error {
	err := tx.Debug().Create(&feedback).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *FeedbackRepository) GetFeedbackByID(tx *gorm.DB, feedbackID uuid.UUID) (*entity.Feedback, error) {
	var feedback entity.Feedback
	err := tx.Debug().Where("feedback_id = ?", feedbackID).First(&feedback).Error
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

func (r *FeedbackRepository) DeleteFeedback(tx *gorm.DB, feedbackID uuid.UUID) error {
	err := tx.Debug().Where("feedback_id = ?", feedbackID).Delete(&entity.Feedback{}).Error
	if err != nil {
		return err
	}
	return nil
}
