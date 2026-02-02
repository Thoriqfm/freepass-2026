package rest

import (
	"fmt"
	"freepass-2026/internal/service"
	"os"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	router  *gin.Engine
	service *service.Service
}

func NewRest(service *service.Service) *Rest {
	return &Rest{
		router:  gin.Default(),
		service: service,
	}
}

func (r *Rest) MountEndPoint() {
	baseURL := r.router.Group("/api/bcc")

	auth := baseURL.Group("/auth")
	auth.POST("/register", r.RegisterHandler)
}

func (r *Rest) Run() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")

	r.router.Run(fmt.Sprintf("%s:%s", addr, port))

}
