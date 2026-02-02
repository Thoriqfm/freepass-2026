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
