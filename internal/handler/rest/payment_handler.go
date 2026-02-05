package rest

import (
	"freepass-2026/entity"
	"freepass-2026/model"
	"freepass-2026/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) CreatePayment(c *gin.Context) {
	userEntity, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	user := userEntity.(*entity.User)

	orderIDStr := c.Param("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order ID", err)
		return
	}

	var param model.CreatePaymentParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	// Set OrderID dari URL parameter, bukan dari body
	param.OrderID = orderID

	resp, err := r.service.PaymentService.CreatePayment(user.UserID, param)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create payment", err)
		return
	}

	response.Success(c, http.StatusOK, "payment created successfully", resp)
}
