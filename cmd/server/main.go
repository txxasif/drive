package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"drive/internal/bootstrap"
	"drive/internal/config"
	"drive/internal/util"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		os.Exit(1)
	}

	// Initialize logger
	logger := util.NewLogger(cfg.Logging.Level)
	logger.Info("Configuration loaded successfully")

	// Create app instance
	app, err := bootstrap.NewApp(cfg)
	if err != nil {
		logger.Fatal("Failed to create app", util.WithError(err))
	}
	logger.Info("App instance created successfully")

	// Close database connection when the application exits
	defer func() {
		sqlDB, err := app.Database.DB()
		if err != nil {
			logger.Error("Failed to get underlying *sql.DB", util.WithError(err))
			return
		}
		if err := sqlDB.Close(); err != nil {
			logger.Error("Failed to close database connection", util.WithError(err))
		}
	}()

	// Create server
	srv := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: app.Router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", util.WithPath(cfg.Server.Address))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", util.WithError(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", util.WithError(err))
	}

	logger.Info("Server exited properly")
}
