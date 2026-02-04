package rest

import (
	"fmt"
	"freepass-2026/internal/service"
	"freepass-2026/pkg/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	router     *gin.Engine
	service    *service.Service
	middleware middleware.Interface
}

func NewRest(service *service.Service, middleware middleware.Interface) *Rest {
	return &Rest{
		router:     gin.Default(),
		service:    service,
		middleware: middleware,
	}
}

func (r *Rest) MountEndPoint() {
	baseURL := r.router.Group("/api/bcc")

	auth := baseURL.Group("/auth")
	auth.POST("/register", r.RegisterHandler)
	auth.POST("/login", r.LoginUser)
	auth.POST("/login/admin", r.LoginAdmin)
	auth.POST("/login/owner", r.LoginOwner)

	user := baseURL.Group("/user")
	user.Use(r.middleware.AuthenticateUser)
	user.GET("/profile", r.GetUserProfile)
	user.PUT("/profile/update", r.UpdateProfile)
	user.GET("/canteen", r.GetAllCanteens)
	user.GET("/canteen/:canteen_id/menus", r.GetAllAvailableMenusByCanteen)
	user.POST("/order/create", r.CreateOrder)
	user.GET("/orders", r.GetUserOrders)
	user.GET("/order/:order_id", r.GetOrderDetail)

	admin := baseURL.Group("/admin")
	admin.Use(r.middleware.AuthenticateUser, r.middleware.OnlyAdmin)
	admin.POST("/create-canteen-owner", r.CreateCanteenOwner)
	admin.PUT("/canteen-owner/:owner_id/update", r.UpdateCanteenOwnerProfile)
	admin.DELETE("/delete-user/:user_id", r.DeleteUser)

	owner := baseURL.Group("/owner")
	owner.Use(r.middleware.AuthenticateUser, r.middleware.OnlyOwner)
	owner.POST("/canteen/create", r.CreateCanteen)
	owner.POST("/canteen/:canteen_id/create-menu", r.CreateMenu)
	owner.PUT("/canteen/:canteen_id/update-status", r.UpdateCanteenStatus)
	owner.GET("/canteen/:canteen_id/menus", r.GetMenusByCanteen)
	owner.GET("/menu/:menu_id", r.GetMenuByMenuID)
	owner.PUT("/menu/:menu_id/update", r.UpdateMenu)
	owner.DELETE("/menu/:menu_id/delete", r.DeleteMenu)
	owner.PUT("/order/:order_id/update-status", r.UpdateOrderStatus)

}

func (r *Rest) Run() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")

	r.router.Run(fmt.Sprintf("%s:%s", addr, port))

}
