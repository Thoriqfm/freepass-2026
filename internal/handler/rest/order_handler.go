package rest

import (
	"freepass-2026/entity"
	"freepass-2026/model"
	"freepass-2026/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) CreateOrder(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	userEntity, ok := user.(*entity.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "failed to get user session", nil)
		return
	}

	var param model.CreateOrderParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	resp, err := r.service.OrderService.CreateOrder(userEntity.UserID, param)
	if err != nil {
		switch err.Error() {
		case "canteen not found":
			response.Error(c, http.StatusNotFound, "canteen not found", err)
			return
		case "canteen is closed":
			response.Error(c, http.StatusForbidden, "canteen is closed", err)
			return
		case "menu not found":
			response.Error(c, http.StatusNotFound, "menu not found", err)
			return
		case "menu does not belong to this canteen":
			response.Error(c, http.StatusBadRequest, "menu does not belong to this canteen", err)
			return
		case "menu is not available":
			response.Error(c, http.StatusBadRequest, "menu is not available", err)
			return
		default:
			if err.Error()[:21] == "insufficient stock for" {
				response.Error(c, http.StatusConflict, "insufficient stock", err)
				return
			}
			response.Error(c, http.StatusInternalServerError, "failed to create order", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "order created successfully", resp)
}
