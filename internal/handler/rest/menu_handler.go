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

func (r *Rest) GetMenusByCanteen(c *gin.Context) {
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

	resp, err := r.service.MenuService.GetMenusByCanteen(ownerEntity.UserID, canteenUUID)
	if err != nil {
		switch err.Error() {
		case "canteen not found":
			response.Error(c, http.StatusNotFound, "canteen not found", err)
			return
		case "access denied: canteen not owned by this owner":
			response.Error(c, http.StatusForbidden, "canteen not owned by this owner", err)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get menus", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "menus retrieved successfully", resp)

}

func (r *Rest) GetMenuByMenuID(c *gin.Context) {
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

	menuID := c.Param("menu_id")
	menuUUID, err := uuid.Parse(menuID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid menu ID", err)
		return
	}

	resp, err := r.service.MenuService.GetMenuByMenuID(ownerEntity.UserID, menuUUID)
	if err != nil {
		switch err.Error() {
		case "menu not found":
			response.Error(c, http.StatusNotFound, "menu not found", err)
			return
		case "access denied: menu not owned by this owner":
			response.Error(c, http.StatusForbidden, "menu not owned by this owner", err)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get menu", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "menu detail retrieved successfully", resp)
}

func (r *Rest) UpdateMenu(c *gin.Context) {
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

	menuID := c.Param("menu_id")
	menuUUID, err := uuid.Parse(menuID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid menu ID", err)
		return
	}

	var param model.UpdateMenuParam
	err = c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	resp, err := r.service.MenuService.UpdateMenu(ownerEntity.UserID, menuUUID, param)
	if err != nil {
		switch err.Error() {
		case "menu not found":
			response.Error(c, http.StatusNotFound, "menu not found", err)
			return
		case "canteen not found":
			response.Error(c, http.StatusNotFound, "canteen not found", err)
			return
		case "access denied: canteen not owned by this owner":
			response.Error(c, http.StatusForbidden, "canteen not owned by this owner", err)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to update menu", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "menu updated successfully", resp)
}

func (r *Rest) DeleteMenu(c *gin.Context) {
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

	menuID := c.Param("menu_id")
	menuUUID, err := uuid.Parse(menuID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid menu ID", err)
		return
	}

	resp, err := r.service.MenuService.DeleteMenu(ownerEntity.UserID, menuUUID)
	if err != nil {
		switch err.Error() {
		case "menu not found":
			response.Error(c, http.StatusNotFound, "menu not found", err)
			return
		case "canteen not found":
			response.Error(c, http.StatusNotFound, "canteen not found", err)
			return
		case "access denied: canteen not owned by this owner":
			response.Error(c, http.StatusForbidden, "canteen not owned by this owner", err)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to delete menu", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "menu deleted successfully", resp)
}
