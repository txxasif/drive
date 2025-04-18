package main

import (
	"drive/internal/config"
	"drive/internal/database/migration"
	"drive/internal/util"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	rollback := flag.Bool("rollback", false, "Rollback the last migration")
	rollbackN := flag.Int("n", 1, "Number of migrations to rollback")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Create logger
	logger := util.NewLogger(zapcore.InfoLevel)

	// Setup database connection
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		os.Exit(1)
	}

	// Create migrator
	migrator := migration.RegisterMigrations(db, logger)

	// Run migrations or rollback
	if *rollback {
		logger.Info("Rolling back migrations", zap.Int("count", *rollbackN))
		if err := migrator.Rollback(*rollbackN); err != nil {
			logger.Error("Failed to rollback migrations", zap.Error(err))
			os.Exit(1)
		}
		logger.Info("Migration rollback completed successfully")
	} else {
		logger.Info("Running migrations")
		if err := migrator.Migrate(); err != nil {
			logger.Error("Failed to run migrations", zap.Error(err))
			os.Exit(1)
		}
		logger.Info("Migrations completed successfully")
	}

	os.Exit(0)
}
