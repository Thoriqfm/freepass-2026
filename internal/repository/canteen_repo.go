package repository

import (
	"freepass-2026/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICanteenRepository interface {
	GetCanteenByID(tx *gorm.DB, canteenID uuid.UUID) (*entity.Canteen, error)
}

type CanteenRepository struct {
	db *gorm.DB
}

func NewCanteenRepository(db *gorm.DB) ICanteenRepository {
	return &CanteenRepository{db: db}
}

func (r *CanteenRepository) GetCanteenByID(tx *gorm.DB, canteenID uuid.UUID) (*entity.Canteen, error) {
	var canteen entity.Canteen
	err := tx.Where("canteen_id = ?", canteenID).First(&canteen).Error
	if err != nil {
		return nil, err
	}
	return &canteen, nil
}
