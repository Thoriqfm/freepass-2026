package repository

import (
	"freepass-2026/entity"

	"gorm.io/gorm"
)

type IMenuRepository interface {
	CreateMenu(tx *gorm.DB, menu *entity.Menu) error
}

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) IMenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) CreateMenu(tx *gorm.DB, menu *entity.Menu) error {
	err := tx.Debug().Create(&menu).Error
	if err != nil {
		return err
	}
	return nil
}
