package rest

import (
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
