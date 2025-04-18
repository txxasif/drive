package migration

import (
	"drive/internal/util"

	"gorm.io/gorm"
)

// RegisterMigrations registers all migrations
func RegisterMigrations(db *gorm.DB, logger *util.Logger) *Migrator {
	migrator := NewMigrator(db, logger)

	// Register migrations in order
	migrator.AddMigration(&CreateUsersTable{})
	// Add more migrations here as your application grows
	// migrator.AddMigration(&CreateFoldersTable{})
	// migrator.AddMigration(&CreateFilesTable{})

	return migrator
}
