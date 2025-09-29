package app

import (
	"fmt"
	"net/http"
	"painaway_test/internal/auth"
	"painaway_test/internal/config"
	"painaway_test/internal/diary"
	logm "painaway_test/internal/log"
	"painaway_test/internal/notifications"
	db "painaway_test/internal/storage"
	"painaway_test/internal/users"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type App struct {
	Config *config.Config
	Logger *zap.Logger
	DB     *gorm.DB
	Hub    *notifications.Hub
	Server *http.Server
}

func New() (*App, error) {

	// Init conifg
	cfg, err := config.LoadConfig("config")
	if err != nil {
		return nil, err
	}

	// Init logger
	logger, _ := NewLogger(cfg)
	defer logger.Sync()
	logger.Info("App initialized", zap.String("env", cfg.Env))

	// Init storage
	dbConn, err := db.Connect(&cfg.DBConfig)
	if err != nil {
		return nil, err
	}

	// Migrations
	if err := db.AutoMigrate(dbConn); err != nil {
		return nil, err
	}

	// Init Hub notifications
	hub := notifications.NewHub()

	// Init router
	router := buildRouter(cfg, logger, dbConn, hub)

	address := fmt.Sprintf(":%v", cfg.HTTPServerConfig.ServerPort)
	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	return &App{
		Config: cfg,
		Logger: logger,
		DB:     dbConn,
		Hub:    hub,
		Server: srv,
	}, nil
}

func NewLogger(env *config.Config) (*zap.Logger, error) {
	var cfg zap.Config
	if env.Env == "prod" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.DisableStacktrace = true
	}

	return cfg.Build()
}

func buildRouter(cfg *config.Config, logger *zap.Logger, dbConn *gorm.DB, hub *notifications.Hub) *gin.Engine {
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
	auth.RegisterRoutes(public, authService, &cfg.JWTConfig, logger)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(auth.AuthMiddleware(&cfg.JWTConfig, logger))
	notifications.RegisterRoutes(protected, notifService, hub, logger)
	diary.RegisterRoutes(protected, diaryService, logger)
	users.RegisterRoutes(protected, userService, logger)

	return router
}
