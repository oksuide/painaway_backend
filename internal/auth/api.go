package auth

import (
	"net/http"
	"painaway_test/internal/config"
	"painaway_test/internal/utils"
	"painaway_test/models"
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
	var req struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		FatherName  string `json:"father_name"`
		Sex         string `json:"sex"`
		DateOfBirth string `json:"date_of_birth"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Warn("Failed to bind register request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dob time.Time
	if req.DateOfBirth != "" {
		dob, _ = time.Parse("2006-01-02", req.DateOfBirth)
	}

	user := &models.User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		FatherName:  req.FatherName,
		Sex:         req.Sex,
		DateOfBirth: dob,
		Groups:      "Patient", // по умолчанию
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}

	if err := h.Service.Register(user); err != nil {
		h.Logger.Error("Failed to register user", zap.Error(err), zap.String("email", user.Email), zap.String("username", user.Username))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Logger.Info("User registered successfully", zap.String("email", user.Email), zap.String("username", user.Username))
	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Warn("Failed to bind login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.Login(req.Username, req.Password)
	if err != nil {
		h.Logger.Warn("Invalid login attempt", zap.String("username", req.Username))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, _ := utils.GenerateAccessToken(*h.JWTConfig, user.ID, user.Email, user.Groups)
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
