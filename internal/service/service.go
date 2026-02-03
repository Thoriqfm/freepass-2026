package service

import (
	"freepass-2026/internal/repository"
	"freepass-2026/pkg/bcrypt"
	"freepass-2026/pkg/database"
	"freepass-2026/pkg/jwt"
)

type Service struct {
	UserService    IUserService
	MenuService    IMenuService
	CanteenService ICanteenService
}

func NewService(repository *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface) *Service {
	return &Service{
		UserService:    NewUserService(repository.UserRepository, bcrypt, jwtAuth),
		MenuService:    NewMenuService(repository.MenuRepository, repository.CanteenRepository, database.Connection),
		CanteenService: NewCanteenService(repository.CanteenRepository, database.Connection),
	}
}
