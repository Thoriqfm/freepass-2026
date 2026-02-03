package service

import (
	"errors"
	"freepass-2026/entity"
	"freepass-2026/internal/repository"
	"freepass-2026/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IMenuService interface {
	CreateMenu(ownerID uuid.UUID, canteenID uuid.UUID, param model.CreateMenuParam) error
}

type MenuService struct {
	db                *gorm.DB
	menuRepository    repository.IMenuRepository
	canteenRepository repository.ICanteenRepository
}

func NewMenuService(menuRepository repository.IMenuRepository, canteenRepository repository.ICanteenRepository, db *gorm.DB) IMenuService {
	return &MenuService{
		menuRepository:    menuRepository,
		canteenRepository: canteenRepository,
		db:                db,
	}
}

func (m *MenuService) CreateMenu(ownerID uuid.UUID, canteenID uuid.UUID, param model.CreateMenuParam) error {
	tx := m.db.Begin()
	defer tx.Rollback()

	canteen, err := m.canteenRepository.GetCanteenByID(tx, canteenID)
	if err != nil {
		return errors.New("canteen not found")
	}

	if canteen.OwnerID != ownerID {
		return errors.New("access denied: canteen not owned by this owner")
	}

	menuID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	menu := &entity.Menu{
		MenuID:      menuID,
		CanteenID:   canteenID,
		Name:        param.Name,
		Description: param.Description,
		Price:       param.Price,
		Stock:       param.Stock,
		IsAvailable: param.IsAvailable,
	}

	err = m.menuRepository.CreateMenu(tx, menu)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}
