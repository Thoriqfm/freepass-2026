package repository

import (
	"freepass-2026/entity"

	"gorm.io/gorm"
)

type IFeedbackRepository interface {
	CreateFeedback(tx *gorm.DB, feedback *entity.Feedback) error
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
