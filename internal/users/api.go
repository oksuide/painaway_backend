package users

import (
	"net/http"
	"painaway_test/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	Service *Service
	Logger  *zap.Logger
}

func RegisterRoutes(rg *gin.RouterGroup, service *Service, logger *zap.Logger) {
	h := &Handler{Service: service, Logger: logger}
	rg.GET("/auth/profile", h.GetProfile)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.NewErrorRepsonse(c, http.StatusUnauthorized, "unauthorized", h.Logger)
		return
	}
	user, err := h.Service.GetProfile(userID.(uint))
	if err != nil {
		h.Logger.Error("failed to get user profile",
			zap.Uint("userID", userID.(uint)),
			zap.Error(err))

		response.NewErrorRepsonse(c, http.StatusInternalServerError, "failed to fetch user profile", h.Logger)
		return
	}
	h.Logger.Info("user profile retrieved", zap.Uint("userID", userID.(uint)))

	respData := gin.H{
		"id":            user.ID,
		"username":      user.Username,
		"email":         user.Email,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"father_name":   user.FatherName,
		"sex":           user.Sex,
		"date_of_birth": user.DateOfBirth.Format("02.01.2006"),
		"groups":        user.Groups,
	}

	c.JSON(http.StatusOK, respData)
}
