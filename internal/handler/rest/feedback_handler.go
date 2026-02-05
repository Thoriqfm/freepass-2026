package rest

import (
	"freepass-2026/entity"
	"freepass-2026/model"
	"freepass-2026/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) CreateFeedback(c *gin.Context) {
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

	var param model.CreateFeedbackParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	resp, err := r.service.FeedbackService.CreateFeedback(userEntity.UserID, param)
	if err != nil {
		switch err.Error() {
		case "order not found":
			response.Error(c, http.StatusNotFound, "order not found", err)
			return
		case "access denied: order does not belong to this user":
			response.Error(c, http.StatusForbidden, "access denied", err)
			return
		case "feedback can only be given for completed orders":
			response.Error(c, http.StatusConflict, "order not completed", err)
			return
		case "feedback can only be given for paid orders":
			response.Error(c, http.StatusConflict, "order not paid", err)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to create feedback", err)
			return
		}
	}

	response.Success(c, http.StatusOK, "feedback created successfully", resp)
}

func (r *Rest) DeleteFeedback(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	owner := user.(*entity.User)
	feedbackIDParam := c.Param("feedback_id")
	feedbackID, err := uuid.Parse(feedbackIDParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid feedback ID", err)
		return
	}

	resp, err := r.service.FeedbackService.DeleteFeedback(owner.UserID, feedbackID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to remove feedback", err)
		return
	}

	response.Success(c, http.StatusOK, "feedback removed successfully", resp)
}
