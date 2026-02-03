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
	GetMenusByCanteen(ownerID uuid.UUID, canteenID uuid.UUID) (*model.MenuListResponse, error)
	GetMenuByMenuID(ownerID uuid.UUID, menuID uuid.UUID) (*model.MenuResponse, error)
	UpdateMenu(ownerID uuid.UUID, menuID uuid.UUID, param model.UpdateMenuParam) (*model.UpdateMenuResponse, error)
	DeleteMenu(ownerID uuid.UUID, menuID uuid.UUID) (*model.DeleteMenuResponse, error)
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

func (m *MenuService) GetMenusByCanteen(ownerID uuid.UUID, canteenID uuid.UUID) (*model.MenuListResponse, error) {
	tx := m.db.Begin()
	defer tx.Rollback()

	canteen, err := m.canteenRepository.GetCanteenByID(tx, canteenID)
	if err != nil {
		return nil, errors.New("canteen not found")
	}
	if canteen.OwnerID != ownerID {
		return nil, errors.New("access denied: canteen not owned by this owner")
	}

	menus, err := m.menuRepository.GetMenuByCanteen(tx, canteenID)
	if err != nil {
		return nil, err
	}

	// convert to response model
	var menuResponses []model.MenuResponse
	for _, menu := range menus {
		menuResponses = append(menuResponses, model.MenuResponse{
			MenuID:      menu.MenuID,
			CanteenID:   menu.CanteenID,
			Name:        menu.Name,
			Description: menu.Description,
			Price:       menu.Price,
			Stock:       menu.Stock,
			IsAvailable: menu.IsAvailable,
			CreatedAt:   menu.CreatedAt,
			UpdatedAt:   menu.UpdatedAt,
		})
	}

	response := &model.MenuListResponse{
		Menus: menuResponses,
		Total: len(menuResponses),
	}

	return response, nil
}

func (m *MenuService) GetMenuByMenuID(ownerID uuid.UUID, menuID uuid.UUID) (*model.MenuResponse, error) {
	tx := m.db.Begin()
	defer tx.Rollback()

	menu, err := m.menuRepository.GetMenuByID(tx, menuID)
	if err != nil {
		return nil, errors.New("menu not found")
	}

	canteen, err := m.canteenRepository.GetCanteenByID(tx, menu.CanteenID)
	if err != nil {
		return nil, errors.New("canteen not found")
	}

	if canteen.OwnerID != ownerID {
		return nil, errors.New("access denied: canteen not owned by this owner")
	}

	// response
	response := &model.MenuResponse{
		MenuID:      menu.MenuID,
		CanteenID:   menu.CanteenID,
		Name:        menu.Name,
		Description: menu.Description,
		Price:       menu.Price,
		Stock:       menu.Stock,
		IsAvailable: menu.IsAvailable,
		CreatedAt:   menu.CreatedAt,
		UpdatedAt:   menu.UpdatedAt,
	}

	return response, nil
}

func (m *MenuService) UpdateMenu(ownerID uuid.UUID, menuID uuid.UUID, param model.UpdateMenuParam) (*model.UpdateMenuResponse, error) {
	tx := m.db.Begin()
	defer tx.Rollback()

	// get menu
	menu, err := m.menuRepository.GetMenuByID(tx, menuID)
	if err != nil {
		return nil, errors.New("menu not found")
	}

	canteen, err := m.canteenRepository.GetCanteenByID(tx, menu.CanteenID)
	if err != nil {
		return nil, errors.New("canteen not found")
	}

	if canteen.OwnerID != ownerID {
		return nil, errors.New("access denied: canteen not owned by this owner")
	}

	// update fields
	if param.Name != "" {
		menu.Name = param.Name
	}
	if param.Description != "" {
		menu.Description = param.Description
	}
	if param.Price > 0 {
		menu.Price = param.Price
	}
	if param.Stock >= 0 {
		menu.Stock = param.Stock
		// auto set by stock
		if param.Stock == 0 {
			menu.IsAvailable = false
		} else if param.Stock > 0 {
			if !menu.IsAvailable {
				menu.IsAvailable = true
			}
		}
	}
	if param.IsAvailable != menu.IsAvailable {
		menu.IsAvailable = param.IsAvailable
	}

	err = m.menuRepository.UpdateMenu(tx, menu)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}
	response := &model.UpdateMenuResponse{
		MenuID:      menu.MenuID,
		CanteenID:   menu.CanteenID,
		Name:        menu.Name,
		Description: menu.Description,
		Price:       menu.Price,
		Stock:       menu.Stock,
		IsAvailable: menu.IsAvailable,
		UpdatedAt:   menu.UpdatedAt,
	}
	return response, nil
}

func (m *MenuService) DeleteMenu(ownerID uuid.UUID, menuID uuid.UUID) (*model.DeleteMenuResponse, error) {
	tx := m.db.Begin()
	defer tx.Rollback()

	menu, err := m.menuRepository.GetMenuByID(tx, menuID)
	if err != nil {
		return nil, errors.New("menu not found")
	}

	canteen, err := m.canteenRepository.GetCanteenByID(tx, menu.CanteenID)
	if err != nil {
		return nil, errors.New("canteen not found")
	}
	if canteen.OwnerID != ownerID {
		return nil, errors.New("access denied: canteen not owned by this owner")
	}

	err = m.menuRepository.DeleteMenu(tx, menuID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.DeleteMenuResponse{
		MenuID:  menuID,
		Message: "menu deleted successfully",
	}

	return response, nil
}
