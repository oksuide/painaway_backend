package users

import (
	"net/http"

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
	userID := c.MustGet("userID").(uint)
	user, err := h.Service.GetProfile(userID)
	if err != nil {
		h.Logger.Error("failed to get user profile", zap.Uint("userID", userID), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.Logger.Info("user profile retrieved", zap.Uint("userID", userID))
	c.JSON(http.StatusOK, gin.H{
		"id":            user.ID,
		"username":      user.Username,
		"email":         user.Email,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"father_name":   user.FatherName,
		"sex":           user.Sex,
		"date_of_birth": user.DateOfBirth.Format("02.01.2006"),
		"groups":        user.Groups,
	})
}
