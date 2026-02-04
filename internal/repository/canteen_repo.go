package repository

import (
	"freepass-2026/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICanteenRepository interface {
	CreateCanteen(tx *gorm.DB, canteen *entity.Canteen) error
	GetCanteenByID(tx *gorm.DB, canteenID uuid.UUID) (*entity.Canteen, error)
	GetAllCanteens(tx *gorm.DB) ([]entity.Canteen, error)
	UpdateCanteenStatus(tx *gorm.DB, canteen *entity.Canteen) error
}

type CanteenRepository struct {
	db *gorm.DB
}

func NewCanteenRepository(db *gorm.DB) ICanteenRepository {
	return &CanteenRepository{db: db}
}

func (r *CanteenRepository) CreateCanteen(tx *gorm.DB, canteen *entity.Canteen) error {
	err := tx.Debug().Create(&canteen).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CanteenRepository) GetCanteenByID(tx *gorm.DB, canteenID uuid.UUID) (*entity.Canteen, error) {
	var canteen entity.Canteen
	err := tx.Where("canteen_id = ?", canteenID).First(&canteen).Error
	if err != nil {
		return nil, err
	}
	return &canteen, nil
}

func (r *CanteenRepository) GetAllCanteens(tx *gorm.DB) ([]entity.Canteen, error) {
	var canteens []entity.Canteen
	err := tx.Debug().Find(&canteens).Error
	if err != nil {
		return nil, err
	}
	return canteens, nil
}

func (r *CanteenRepository) UpdateCanteenStatus(tx *gorm.DB, canteen *entity.Canteen) error {
	err := tx.Debug().Save(&canteen).Error
	if err != nil {
		return err
	}
	return nil
}
