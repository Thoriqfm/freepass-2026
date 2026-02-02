package service

import (
	"errors"
	"freepass-2026/entity"
	"freepass-2026/internal/repository"
	"freepass-2026/model"
	"freepass-2026/pkg/bcrypt"
	"freepass-2026/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserService interface {
	Register(param model.UserRegisterParam) error
}

type UserService struct {
	db             *gorm.DB
	userRepository repository.IUserRepository
	bcrypt         bcrypt.Interface
}

func NewUserService(userRepository repository.IUserRepository, bcrypt bcrypt.Interface) IUserService {
	return &UserService{
		db:             database.Connection,
		userRepository: userRepository,
		bcrypt:         bcrypt,
	}
}

func (u *UserService) Register(param model.UserRegisterParam) error {
	tx := u.db.Begin()
	defer tx.Rollback()

	_, err := u.userRepository.GetUser(model.UserParam{
		Email: param.Email,
	})

	if err == nil {
		return errors.New("email already exists")
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	if param.Password != param.ConfirmPassword {
		return errors.New("password not match")
	}

	hashPassword, err := u.bcrypt.GenerateFromPassword(param.Password)
	if err != nil {
		return err
	}

	user := &entity.User{
		UserID:   userID,
		RoleID:   2,
		Name:     param.Name,
		Email:    param.Email,
		Phone:    param.Phone,
		Password: hashPassword,
	}

	err = u.userRepository.CreateUser(tx, user)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}
	return nil

}
