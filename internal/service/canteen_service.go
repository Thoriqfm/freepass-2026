package service

import (
	"freepass-2026/entity"
	"freepass-2026/internal/repository"
	"freepass-2026/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICanteenService interface {
	CreateCanteen(ownerID uuid.UUID, param model.CreateCanteenParam) (*model.CreateCanteenResponse, error)
	GetAllCanteens() ([]model.CanteenListResponse, error)
}

type CanteenService struct {
	db                *gorm.DB
	canteenRepository repository.ICanteenRepository
}

func NewCanteenService(canteenRepository repository.ICanteenRepository, db *gorm.DB) ICanteenService {
	return &CanteenService{
		db:                db,
		canteenRepository: canteenRepository,
	}
}

func (c *CanteenService) CreateCanteen(ownerID uuid.UUID, param model.CreateCanteenParam) (*model.CreateCanteenResponse, error) {
	tx := c.db.Begin()
	defer tx.Rollback()

	canteenID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	canteen := &entity.Canteen{
		CanteenID: canteenID,
		OwnerID:   ownerID,
		Name:      param.Name,
		Location:  param.Location,
		IsOpen:    param.IsOpen,
	}

	err = c.canteenRepository.CreateCanteen(tx, canteen)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.CreateCanteenResponse{
		CanteenID: canteen.CanteenID,
		OwnerID:   canteen.OwnerID,
		Name:      canteen.Name,
		Location:  canteen.Location,
		IsOpen:    canteen.IsOpen,
		CreatedAt: canteen.CreatedAt,
	}

	return response, nil
}

func (c *CanteenService) GetAllCanteens() ([]model.CanteenListResponse, error) {
	tx := c.db.Begin()
	defer tx.Rollback()

	canteens, err := c.canteenRepository.GetAllCanteens(tx)
	if err != nil {
		return nil, err
	}

	var responses []model.CanteenListResponse
	for _, canteen := range canteens {
		responses = append(responses, model.CanteenListResponse{
			Name:     canteen.Name,
			Location: canteen.Location,
			IsOpen:   canteen.IsOpen,
		})
	}

	return responses, nil
}
