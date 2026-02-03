package repository

import (
	"freepass-2026/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IMenuRepository interface {
	CreateMenu(tx *gorm.DB, menu *entity.Menu) error
	GetMenuByCanteen(tx *gorm.DB, canteenID uuid.UUID) ([]entity.Menu, error)
	GetMenuByID(tx *gorm.DB, menuID uuid.UUID) (*entity.Menu, error)
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

func (r *MenuRepository) GetMenuByCanteen(tx *gorm.DB, canteenID uuid.UUID) ([]entity.Menu, error) {
	var menus []entity.Menu
	err := tx.Debug().Where("canteen_id = ?", canteenID).Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *MenuRepository) GetMenuByID(tx *gorm.DB, menuID uuid.UUID) (*entity.Menu, error) {
	var menu entity.Menu
	err := tx.Debug().Where("menu_id = ?", menuID).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}
