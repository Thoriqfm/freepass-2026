package rest

import (
	"freepass-2026/entity"
	"freepass-2026/model"
	"freepass-2026/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (r *Rest) GetUserOrders(c *gin.Context) {
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

	resp, err := r.service.OrderService.GetUserOrders(userEntity.UserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get orders", err)
		return
	}

	response.Success(c, http.StatusOK, "orders retrieved successfully", resp)
}

func (r *Rest) GetOrderDetail(c *gin.Context) {
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

	orderID := c.Param("order_id")
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order ID", err)
		return
	}

	resp, err := r.service.OrderService.GetOrderDetail(userEntity.UserID, orderUUID)
	if err != nil {
		switch err.Error() {
		case "order not found":
			response.Error(c, http.StatusNotFound, "order not found", err)
			return
		case "access denied: order does not belong to this user":
			response.Error(c, http.StatusForbidden, "access denied", err)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get order detail", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "order detail retrieved successfully", resp)
}

/*
* Owner handler
 */

func (r *Rest) UpdateOrderStatus(c *gin.Context) {
	owner, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	ownerEntity, ok := owner.(*entity.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "failed to get owner session", nil)
		return
	}

	orderID := c.Param("order_id")
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order ID", err)
		return
	}

	var param struct {
		OrderStatus string `json:"order_status" binding:"required"`
	}

	err = c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	resp, err := r.service.OrderService.UpdateOrderStatus(ownerEntity.UserID, orderUUID, param.OrderStatus)
	if err != nil {
		switch err.Error() {
		case "order not found":
			response.Error(c, http.StatusNotFound, "order not found", err)
			return
		case "canteen not found":
			response.Error(c, http.StatusNotFound, "canteen not found", err)
			return
		case "access denied: canteen not owned by this owner":
			response.Error(c, http.StatusForbidden, "access denied", err)
			return
		case "cannot update order status: payment not yet received":
			response.Error(c, http.StatusConflict, "payment not yet received", err)
			return
		case "invalid order status":
			response.Error(c, http.StatusBadRequest, "invalid order status", err)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to update order status", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "order status updated successfully", resp)
}

// GET /owner/orders
func (r *Rest) GetOwnerAllOrders(c *gin.Context) {
	// Get owner from context
	ownerEntity, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	owner := ownerEntity.(*entity.User)

	// Bind query parameters
	var queryParam model.OwnerOrderQueryParam
	if err := c.ShouldBindQuery(&queryParam); err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind query parameters", err)
		return
	}

	// Call service
	resp, err := r.service.OrderService.GetOwnerAllOrders(owner.UserID, queryParam)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get all orders", err)
		return
	}

	response.Success(c, http.StatusOK, "all orders retrieved successfully", resp)
}
