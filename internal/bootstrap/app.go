package bootstrap

import (
	"drive/internal/config"
	"drive/internal/database"
	"drive/internal/handler"
	"drive/internal/repository"
	"drive/internal/routes"
	"drive/internal/service"
	"drive/internal/util"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Config   *config.Config
	Database *gorm.DB
	Router   http.Handler
	Logger   *util.Logger
}

func NewApp(cfg *config.Config) (*App, error) {
	logger := util.NewLogger(cfg.Logging.Level)
	logger.Info("Initializing application")

	db, err := database.InitDatabase(cfg, logger)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	logger.Info("Database connection established")

	if err := database.RunMigrations(db, logger); err != nil {
		logger.Error("Failed to run migrations", zap.Error(err))
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	logger.Info("Database migrations completed")

	jwtService := util.NewJwtService(util.ServiceConfig{
		SecretKey:     cfg.JWT.Secret,
		AccessExpiry:  cfg.JWT.AccessExpiresIn,
		RefreshExpiry: cfg.JWT.RefreshExpiresIn,
	})

	repo := repository.NewRepositories(db)
	services := service.NewServices(*repo, jwtService, logger, cfg)
	handler := handler.NewHandler(services)
	routes := routes.SetupRoutes(handler, services.Auth)

	logger.Info("Application initialized successfully")

	return &App{
		Config:   cfg,
		Database: db,
		Router:   routes,
		Logger:   logger,
	}, nil
}
