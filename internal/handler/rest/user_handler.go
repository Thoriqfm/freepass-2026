package rest

import (
	"freepass-2026/entity"
	"freepass-2026/model"
	"freepass-2026/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) RegisterHandler(c *gin.Context) {
	var param model.UserRegisterParam

	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	err = r.service.UserService.Register(param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to register user", err)
		return
	}
	response.Success(c, http.StatusOK, "success to register user", nil)
}

func (r *Rest) LoginUser(c *gin.Context) {
	var param model.UserLoginParam

	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	reps, err := r.service.UserService.Login(param)
	if err != nil {
		if err.Error() == "email or password is wrong" {
			response.Error(c, http.StatusUnauthorized, "email or password is wrong", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to login user", err)
		return
	}

	response.Success(c, http.StatusOK, "user logged in successfully", reps)

}

func (r *Rest) GetUserProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	userEntity, ok := user.(*entity.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "failed to get user", nil)
		return
	}

	reps, err := r.service.UserService.GetUserProfile(userEntity.UserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get user profile", err)
		return
	}

	response.Success(c, http.StatusOK, "user profile retrieved successfully", reps)
}

func (r *Rest) UpdateProfile(c *gin.Context) {
	var param model.UpdateUserProfile

	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

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

	reps, err := r.service.UserService.UpdateUserProfile(userEntity.UserID, param)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update user profile", err)
		return
	}

	response.Success(c, http.StatusOK, "user profile updated successfully", reps)
}

/*
* ADMIN FEATURES
 */

func (r *Rest) LoginAdmin(c *gin.Context) {
	var param model.UserLoginParam

	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	reps, err := r.service.UserService.LoginAdmin(param)
	if err != nil {
		switch err.Error() {
		case "email or password is wrong":
			response.Error(c, http.StatusUnauthorized, "email or password is wrong", err)
			return
		case "access denied: admin only":
			response.Error(c, http.StatusForbidden, "access denied: admin only", err)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to login admin", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "admin logged in successfully", reps)
}

/*
* OWNER FEATURES
 */

func (r *Rest) LoginOwner(c *gin.Context) {
	var param model.UserLoginParam

	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	reps, err := r.service.UserService.LoginOwner(param)
	if err != nil {
		switch err.Error() {
		case "email or password is wrong":
			response.Error(c, http.StatusUnauthorized, "email or password is wrong", err)
			return
		case "access denied: owner only":
			response.Error(c, http.StatusForbidden, "access denied: owner only", err)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to login owner", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "owner logged in successfully", reps)
}
