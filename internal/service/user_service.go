package service

import (
	"errors"
	"freepass-2026/entity"
	"freepass-2026/internal/repository"
	"freepass-2026/model"
	"freepass-2026/pkg/bcrypt"
	"freepass-2026/pkg/database"
	"freepass-2026/pkg/jwt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserService interface {
	Register(param model.UserRegisterParam) error
	GetUser(param model.UserParam) (*entity.User, error)
	Login(param model.UserLoginParam) (*model.UserLoginResponse, error)
	GetUserProfile(userID uuid.UUID) (*model.UserProfile, error)
	UpdateUserProfile(userId uuid.UUID, param model.UpdateUserProfile) (*model.UserProfile, error)
	// ADMIN FEATURES
	LoginAdmin(param model.UserLoginParam) (*model.UserLoginResponse, error)
	CreateCanteenOwner(param model.CreateCanteenOwnerParam) error
	UpdateCanteenOwnerProfile(ownerID uuid.UUID, param model.UpdateCanteenOwnerProfile) (*model.UpdateCanteenOwnerResponse, error)
	DeleteUser(adminID uuid.UUID, targetUserID uuid.UUID) (*model.DeleteUserResponse, error)
	// OWNER FEATURES
	LoginOwner(param model.UserLoginParam) (*model.UserLoginResponse, error)
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

func (U *UserService) CreateCanteenOwner(param model.CreateCanteenOwnerParam) error {
	tx := U.db.Begin()
	defer tx.Rollback()

	_, err := U.userRepository.GetUser(model.UserParam{
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

	hashPassword, err := U.bcrypt.GenerateFromPassword(param.Password)
	if err != nil {
		return err
	}

	user := &entity.User{
		UserID:   userID,
		RoleID:   3,
		Name:     param.Name,
		Email:    param.Email,
		Phone:    param.Phone,
		Password: hashPassword,
	}

	err = U.userRepository.CreateUser(tx, user)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) UpdateCanteenOwnerProfile(ownerID uuid.UUID, param model.UpdateCanteenOwnerProfile) (*model.UpdateCanteenOwnerResponse, error) {
	tx := u.db.Begin()
	defer tx.Rollback()

	owner, err := u.userRepository.GetUserByID(tx, ownerID)
	if err != nil {
		return nil, errors.New("canteen owner not found")
	}

	// validate if user was canteen owner (role = 3)
	if owner.RoleID != 3 {
		return nil, errors.New("user is not a canteen owner")
	}

	// update fields
	if param.Name != "" {
		owner.Name = param.Name
	}
	if param.Email != "" {
		owner.Email = param.Email
	}
	if param.Phone != "" {
		owner.Phone = param.Phone
	}
	if param.Password != "" {
		hashPassword, err := u.bcrypt.GenerateFromPassword(param.Password)
		if err != nil {
			return nil, err
		}
		owner.Password = hashPassword
	}

	err = u.userRepository.UpdateUserProfile(tx, owner)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.UpdateCanteenOwnerResponse{
		UserID: owner.UserID,
		Name:   owner.Name,
		Email:  owner.Email,
		Phone:  owner.Phone,
	}

	return response, nil
}

func (u *UserService) DeleteUser(adminID uuid.UUID, targetUserID uuid.UUID) (*model.DeleteUserResponse, error) {
	tx := u.db.Begin()
	defer tx.Rollback()

	// validate admin role
	admin, err := u.userRepository.GetUserByID(tx, adminID)
	if err != nil {
		return nil, errors.New("admin not found")
	}

	if admin.RoleID != 1 {
		return nil, errors.New("access denied: admin only")
	}

	// get target user
	targetUser, err := u.userRepository.GetUserByID(tx, targetUserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if targetUser.RoleID == 3 {
		if len(targetUser.CanteenOwner) > 0 { // check if owner has active canteens
			return nil, errors.New("cannot delete canteen owner with active canteens")
		}
	}

	// delete user
	err = u.userRepository.DeleteUser(tx, targetUserID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	response := &model.DeleteUserResponse{
		Message:   "user deleted successfully",
		UserID:    targetUserID,
		DeletedAt: time.Now(),
	}

	return response, nil

}

/*
* OWNER FEATURES
 */

func (u *UserService) LoginOwner(param model.UserLoginParam) (*model.UserLoginResponse, error) {
	user, err := u.userRepository.GetUser(model.UserParam{
		Email: param.Email,
	})
	if err != nil {
		return nil, errors.New("email or password is wrong")
	}

	// Validate owner role (3)
	if user.RoleID != 3 {
		return nil, errors.New("access denied: owner only")
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
