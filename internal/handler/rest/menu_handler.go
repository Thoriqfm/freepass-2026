package rest

import (
	"freepass-2026/entity"
	"freepass-2026/model"
	"freepass-2026/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) CreateMenu(c *gin.Context) {
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

	var param model.CreateMenuParam
	err = c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	err = r.service.MenuService.CreateMenu(ownerEntity.UserID, canteenUUID, param)
	if err != nil {
		switch err.Error() {
		case "canteen not found":
			response.Error(c, http.StatusNotFound, "canteen not found", err)
			return
		case "access denied: canteen not owned by this owner":
			response.Error(c, http.StatusForbidden, "canteen not owned by this owner", err)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to create menu", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "menu created successfully", nil)

}
