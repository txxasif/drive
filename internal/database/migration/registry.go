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
	migrator.AddMigration(&CreateFoldersTable{})
	migrator.AddMigration(&CreateFilesTable{})
	migrator.AddMigration(&CreateSharesTable{})

	return migrator
}
