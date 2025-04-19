package migration

import (
	"drive/internal/model"

	"gorm.io/gorm"
)

// CreateFilesTable migration creates the files table
type CreateFilesTable struct{}

// ID returns the migration ID
func (m *CreateFilesTable) ID() string {
	return "003_create_files_table"
}

// Migrate runs the migration
func (m *CreateFilesTable) Migrate(tx *gorm.DB) error {
	return tx.AutoMigrate(&model.File{})
}

// Rollback runs the migration rollback
func (m *CreateFilesTable) Rollback(tx *gorm.DB) error {
	return tx.Migrator().DropTable("files")
}
