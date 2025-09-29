package auth

import (
	"net/http"
	"painaway_test/internal/config"
	"painaway_test/internal/response"
	"painaway_test/internal/utils"
	"painaway_test/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	Service   *Service
	JWTConfig *config.JWTConfig
	Logger    *zap.Logger
}

func RegisterRoutes(rg *gin.RouterGroup, service *Service, jwtCfg *config.JWTConfig, logger *zap.Logger) {
	h := &Handler{Service: service, JWTConfig: jwtCfg, Logger: logger}

	rg.POST("auth/login", h.Login)
	rg.POST("auth/register", h.Register)
}

func (h *Handler) Register(c *gin.Context) {
	var input struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		FatherName  string `json:"father_name"`
		Sex         string `json:"sex"`
		DateOfBirth string `json:"date_of_birth"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid request body", h.Logger)
		return
	}
	if len(input.Password) < 6 {
		response.NewErrorResponse(c, http.StatusBadRequest, "password must be at least 6 characters long", h.Logger)
		return
	}

	dob, err := time.Parse(time.DateOnly, input.DateOfBirth)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "date_of_birth must be in YYYY-MM-DD format", h.Logger)
		return
	}

	user := &models.User{
		Username:    strings.TrimSpace(input.Username),
		Email:       input.Email,
		Password:    input.Password,
		FirstName:   strings.TrimSpace(input.FirstName),
		LastName:    strings.TrimSpace(input.LastName),
		FatherName:  strings.TrimSpace(input.FatherName),
		Sex:         input.Sex,
		DateOfBirth: dob,
		Groups:      "Patient",
	}

	if err := h.Service.Register(user); err != nil {
		h.Logger.Error("Failed to register user", zap.Error(err), zap.String("email", user.Email), zap.String("username", user.Username))
		response.NewErrorResponse(c, http.StatusInternalServerError, "registration failed", h.Logger)
		return
	}

	token, err := utils.GenerateAccessToken(*h.JWTConfig, user.ID, user.Groups)
	if err != nil {
		h.Logger.Error("failed to generate access token", zap.Error(err), zap.Uint("userID", user.ID))
		response.NewErrorResponse(c, http.StatusInternalServerError, "failed to generate token", h.Logger)
		return
	}
	h.Logger.Info("User registered successfully", zap.String("email", user.Email), zap.String("username", user.Username))
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"groups":   user.Groups,
		},
	})
}

func (h *Handler) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid request body", h.Logger)
		return
	}

	user, err := h.Service.Login(input.Username, input.Password)
	if err != nil {
		response.NewErrorResponse(c, http.StatusUnauthorized, "invalid credentials", h.Logger)
		return
	}

	token, _ := utils.GenerateAccessToken(*h.JWTConfig, user.ID, user.Groups)
	h.Logger.Info("User logged in successfully", zap.String("username", user.Username), zap.Uint("user_id", user.ID))

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"groups":   user.Groups,
		},
	})
}
