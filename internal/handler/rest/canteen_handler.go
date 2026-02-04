package rest

import (
	"freepass-2026/entity"
	"freepass-2026/model"
	"freepass-2026/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (r *Rest) GetAllCanteens(c *gin.Context) {
	resp, err := r.service.CanteenService.GetAllCanteens()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get canteens", err)
		return
	}
	response.Success(c, http.StatusOK, "canteens retrieved successfully", resp)
}

func (r *Rest) UpdateCanteenStatus(c *gin.Context) {
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

	canteenID := c.Param("canteen_id")
	canteenUUID, err := uuid.Parse(canteenID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid canteen ID", err)
		return
	}

	var param model.UpdateCantenStatusParam
	err = c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	resp, err := r.service.CanteenService.UpdateCanteenStatus(ownerEntity.UserID, canteenUUID, param)
	if err != nil {
		switch err.Error() {
		case "canteen not found":
			response.Error(c, http.StatusNotFound, "canteen not found", err)
			return
		case "access denied: canteen not owned by this owner":
			response.Error(c, http.StatusForbidden, "canteen not owned by this owner", err)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to update canteen status", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "canteen status updated successfully", resp)

}
