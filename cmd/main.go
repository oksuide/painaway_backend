package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"painaway_test/internal/auth"
	"painaway_test/internal/config"
	"painaway_test/internal/diary"
	logm "painaway_test/internal/log"
	"painaway_test/internal/notifications"
	db "painaway_test/internal/storage"
	"painaway_test/internal/users"
	"syscall"
	"time"

	_ "painaway_test/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// @title PainAway API
// @version 1.0
// @description API for PainAway
// @host localhost:8080
// @BasePath /
func main() {
	// Init conifg
	cfg, err := config.LoadConfig("config")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Init logger
	logger, _ := NewLogger(cfg)
	defer logger.Sync()

	logger.Info("Starting pain-away", zap.String("env", cfg.Env))
	logger.Debug("Debug messages are enabled")

	// Init storage
	dbConn, err := db.Connect(&cfg.DBConfig)
	if err != nil {
		logger.Fatal("Database connection failed", zap.Error(err))
	}
	logger.Info("Database connection established")

	// Migrations
	if err := db.AutoMigrate(dbConn); err != nil {
		logger.Fatal("Database migration failed", zap.Error(err))
	}
	logger.Info("Database migrations completed")
	// Init hub for notifications
	hub := notifications.NewHub()

	// Start server
	address := fmt.Sprintf(":%v", cfg.HTTPServerConfig.ServerPort)
	srv := &http.Server{
		Addr:    address,
		Handler: buildHandler(logger, cfg, dbConn, hub),
	}
	go func() {
		logger.Info("Server is running", zap.String("addr", address))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	switch cfg.Env {
	case "prod":
		return zap.NewProduction()
	default:
		return zap.NewDevelopment()
	}
}

func buildHandler(logger *zap.Logger, cfg *config.Config, dbConn *gorm.DB, hub *notifications.Hub) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logm.LoggerMiddleware(logger))

	// Repositories
	userRepo := users.NewRepository(dbConn)
	diaryRepo := diary.NewRepository(dbConn)
	notifRepo := notifications.NewRepository(dbConn)

	// Services
	authService := auth.NewService(userRepo)
	notifService := notifications.NewService(notifRepo, hub)
	diaryService := diary.NewService(diaryRepo, notifService, logger)
	userService := users.NewService(userRepo)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	public := router.Group("/api")
	{
		auth.RegisterRoutes(public, authService, &cfg.JWTConfig, logger)
	}

	// Protected routes
	protected := router.Group("/api")
	protected.Use(auth.AuthMiddleware(&cfg.JWTConfig))
	{
		notifications.RegisterRoutes(protected, notifService, hub, logger)
		diary.RegisterRoutes(protected, diaryService, logger)
		users.RegisterRoutes(protected, userService, logger)
	}

	return router
}
