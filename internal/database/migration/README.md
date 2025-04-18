# Database Migration System

This package provides a database migration system for the application.

## Overview

The migration system allows for:
- Versioned migrations with sequential IDs
- Transaction support for safe migrations
- Rollback capability
- Tracking of applied migrations

## How It Works

1. Migrations are defined as Go structs that implement the `Migration` interface
2. Each migration has an ID, a `Migrate` method, and a `Rollback` method
3. Migrations are registered in `registry.go`
4. The migrator applies migrations in sequence, skipping those already applied
5. Migration status is tracked in a database table called `migration_records`

## Creating a New Migration

1. Create a new file in this directory with a name like `NNN_description.go` where `NNN` is a sequential number
2. Implement the `Migration` interface
3. Register the migration in `registry.go`

Example:

```go
// 002_create_folders_table.go
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
```

Then register it in `registry.go`:

```go
// Register migrations in order
migrator.AddMigration(&CreateUsersTable{})
migrator.AddMigration(&CreateFoldersTable{})
```

## Running Migrations

Migrations are automatically run during application startup via the `bootstrap.NewApp()` function.

You can also run migrations separately using the `cmd/migrate` tool:

```bash
# Run migrations
go run cmd/migrate/main.go

# Rollback the last migration
go run cmd/migrate/main.go -rollback

# Rollback the last N migrations
go run cmd/migrate/main.go -rollback -n 3
``` 