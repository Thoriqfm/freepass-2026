package middleware

import (
	"freepass-2026/internal/service"
	"freepass-2026/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Interface interface {
	OnlyAdmin(c *gin.Context)
	AuthenticateUser(c *gin.Context)
	OnlyOwner(c *gin.Context)
}

type middleware struct {
	service *service.Service
	jwtAuth jwt.Interface
}

func Init(service *service.Service, jwtAuth jwt.Interface) Interface {
	return &middleware{
		service: service,
		jwtAuth: jwtAuth,
	}
}
