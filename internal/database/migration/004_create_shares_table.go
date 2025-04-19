package migration

import (
	"drive/internal/model"

	"gorm.io/gorm"
)

// CreateSharesTable migration creates the shares table
type CreateSharesTable struct{}

// ID returns the migration ID
func (m *CreateSharesTable) ID() string {
	return "004_create_shares_table"
}

// Migrate runs the migration
func (m *CreateSharesTable) Migrate(tx *gorm.DB) error {
	return tx.AutoMigrate(&model.Share{})
}

// Rollback runs the migration rollback
func (m *CreateSharesTable) Rollback(tx *gorm.DB) error {
	return tx.Migrator().DropTable("shares")
}
