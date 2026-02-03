package rest

import (
	"freepass-2026/entity"
	"freepass-2026/model"
	"freepass-2026/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) CreateCanteen(c *gin.Context) {
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

	var param model.CreateCanteenParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	resp, err := r.service.CanteenService.CreateCanteen(ownerEntity.UserID, param)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create canteen", err)
		return
	}
	response.Success(c, http.StatusOK, "canteen created successfully", resp)
}
