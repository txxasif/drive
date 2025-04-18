package database

import (
	"drive/internal/database/migration"
	"drive/internal/util"

	"gorm.io/gorm"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB, logger *util.Logger) error {
	logger.Info("Running database migrations...")

	migrator := migration.RegisterMigrations(db, logger)
	if err := migrator.Migrate(); err != nil {
		return err
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

// RollbackLastMigration rolls back the last migration
func RollbackLastMigration(db *gorm.DB, logger *util.Logger) error {
	logger.Info("Rolling back last migration...")

	migrator := migration.RegisterMigrations(db, logger)
	if err := migrator.Rollback(1); err != nil {
		return err
	}

	logger.Info("Rollback completed successfully")
	return nil
}
