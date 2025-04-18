package migration

import (
	"drive/internal/model"

	"gorm.io/gorm"
)

// CreateUsersTable migration creates the users table
type CreateUsersTable struct{}

// ID returns the migration ID
func (m *CreateUsersTable) ID() string {
	return "001_create_users_table"
}

// Migrate runs the migration
func (m *CreateUsersTable) Migrate(tx *gorm.DB) error {
	return tx.AutoMigrate(&model.User{})
}

// Rollback runs the migration rollback
func (m *CreateUsersTable) Rollback(tx *gorm.DB) error {
	return tx.Migrator().DropTable("users")
}
