package middleware

import (
	"errors"
	"freepass-2026/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *middleware) OnlyAdmin(c *gin.Context) {
	user, err := m.jwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusForbidden, "failed to get user", err)
		c.Abort()
		return
	}

	if user.RoleID != 1 {
		response.Error(c, http.StatusForbidden, "this endpoint cant be access", errors.New("user dont have access"))
		c.Abort()
		return
	}

	c.Next()
}
