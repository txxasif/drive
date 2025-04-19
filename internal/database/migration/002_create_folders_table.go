package migration

import (
	"drive/internal/model"

	"gorm.io/gorm"
)

// CreateFoldersTable migration creates the folders table
type CreateFoldersTable struct{}

// ID returns the migration ID
func (m *CreateFoldersTable) ID() string {
	return "002_create_folders_table"
}

// Migrate runs the migration
func (m *CreateFoldersTable) Migrate(tx *gorm.DB) error {
	return tx.AutoMigrate(&model.Folder{})
}

// Rollback runs the migration rollback
func (m *CreateFoldersTable) Rollback(tx *gorm.DB) error {
	return tx.Migrator().DropTable("folders")
}
