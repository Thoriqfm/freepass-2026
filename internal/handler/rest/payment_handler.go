package rest

import (
	"freepass-2026/entity"
	"freepass-2026/model"
	"freepass-2026/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) CreatePayment(c *gin.Context) {
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

	var param model.CreatePaymentParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	resp, err := r.service.PaymentService.CreatePayment(userEntity.UserID, param)
	if err != nil {
		switch err.Error() {
		case "order not found":
			response.Error(c, http.StatusNotFound, "order not found", err)
			return
		case "access denied: order does not belong to this user":
			response.Error(c, http.StatusForbidden, "access denied", err)
			return
		case "order already paid":
			response.Error(c, http.StatusConflict, "order already paid", err)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to create payment", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "payment created successfully", resp)
}
