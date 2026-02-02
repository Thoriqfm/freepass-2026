package model

import "github.com/google/uuid"

type UserParam struct {
	UserID uuid.UUID `json:"-"`
	Name   string    `json:"-"`
	Email  string    `json:"-"`
	Phone  string    `json:"-"`
}

type UserRegisterParam struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Phone           string `json:"phone" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

type UserRegisterResponse struct {
}

type UserLoginParam struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	Token  string `json:"token"`
	RoleID int    `json:"role_id"`
}

type UserProfile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
