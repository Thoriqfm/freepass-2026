package service

import (
	"freepass-2026/internal/repository"
	"freepass-2026/pkg/bcrypt"
)

type Service struct {
	UserService IUserService
}

func NewService(repository *repository.Repository, bcrypt bcrypt.Interface) *Service {
	return &Service{
		UserService: NewUserService(repository.UserRepository, bcrypt),
	}
}
