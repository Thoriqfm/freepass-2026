package service

import (
	"errors"
	"freepass-2026/entity"
	"freepass-2026/internal/repository"
	"freepass-2026/model"
	"freepass-2026/pkg/bcrypt"
	"freepass-2026/pkg/database"
	"freepass-2026/pkg/jwt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserService interface {
	Register(param model.UserRegisterParam) error
	GetUser(param model.UserParam) (*entity.User, error)
	Login(param model.UserLoginParam) (*model.UserLoginResponse, error)
	GetUserProfile(userID uuid.UUID) (*model.UserProfile, error)
	UpdateUserProfile(userId uuid.UUID, param model.UpdateUserProfile) (*model.UserProfile, error)
	LoginAdmin(param model.UserLoginParam) (*model.UserLoginResponse, error)
}

type UserService struct {
	db             *gorm.DB
	userRepository repository.IUserRepository
	bcrypt         bcrypt.Interface
	jwtAuth        jwt.Interface
}

func NewUserService(userRepository repository.IUserRepository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface) IUserService {
	return &UserService{
		db:             database.Connection,
		userRepository: userRepository,
		bcrypt:         bcrypt,
		jwtAuth:        jwtAuth,
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

func (u *UserService) GetUser(param model.UserParam) (*entity.User, error) {
	return u.userRepository.GetUser(param)
}

func (u *UserService) Login(param model.UserLoginParam) (*model.UserLoginResponse, error) {
	tx := u.db.Begin()
	defer tx.Rollback()

	user, err := u.userRepository.GetUser(model.UserParam{
		Email: param.Email,
	})
	if err != nil {
		return nil, errors.New("email or password is wrong")
	}

	err = u.bcrypt.CompareAndHashPassword(user.Password, param.Password)
	if err != nil {
		return nil, errors.New("email or password is wrong")
	}

	token, err := u.jwtAuth.CreateJWTToken(user.UserID, false)
	if err != nil {
		return nil, err
	}

	response := &model.UserLoginResponse{
		Token:  token,
		RoleID: user.RoleID,
	}

	return response, nil
}

func (u *UserService) GetUserProfile(userID uuid.UUID) (*model.UserProfile, error) {
	user, err := u.userRepository.GetUser(model.UserParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	response := &model.UserProfile{
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}

	return response, nil
}

func (u *UserService) UpdateUserProfile(userId uuid.UUID, param model.UpdateUserProfile) (*model.UserProfile, error) {
	tx := u.db.Begin()
	defer tx.Rollback()

	user, err := u.userRepository.GetUserByID(tx, userId)
	if err != nil {
		return nil, errors.New("failed to get user data")
	}

	user.Name = param.Name
	user.Phone = param.Phone

	err = u.userRepository.UpdateUserProfile(tx, user)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &model.UserProfile{
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}, nil

}

/*
* ADMIN FEATURES
 */

func (u *UserService) LoginAdmin(param model.UserLoginParam) (*model.UserLoginResponse, error) {
	tx := u.db.Begin()
	defer tx.Rollback()

	user, err := u.userRepository.GetUser(model.UserParam{
		Email: param.Email,
	})
	if err != nil {
		return nil, errors.New("email or password is wrong")
	}

	// validate for admin role (1)
	if user.RoleID != 1 {
		return nil, errors.New("access denied: admin only")
	}

	err = u.bcrypt.CompareAndHashPassword(user.Password, param.Password)
	if err != nil {
		return nil, errors.New("email or password is wrong")
	}

	token, err := u.jwtAuth.CreateJWTToken(user.UserID, false)
	if err != nil {
		return nil, err
	}

	response := &model.UserLoginResponse{
		Token:  token,
		RoleID: user.RoleID,
	}

	return response, nil
}
