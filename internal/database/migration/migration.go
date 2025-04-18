package migration

import (
	"drive/internal/util"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Migration interface defines the methods required for a migration
type Migration interface {
	ID() string
	Migrate(*gorm.DB) error
	Rollback(*gorm.DB) error
}

// MigrationRecord represents a migration record in the database
type MigrationRecord struct {
	ID        string    `gorm:"primaryKey"`
	AppliedAt time.Time `gorm:"autoCreateTime"`
}

// Migrator handles database migrations
type Migrator struct {
	db         *gorm.DB
	migrations []Migration
	logger     *util.Logger
}

// NewMigrator creates a new migrator
func NewMigrator(db *gorm.DB, logger *util.Logger) *Migrator {
	return &Migrator{
		db:         db,
		migrations: []Migration{},
		logger:     logger,
	}
}

// AddMigration adds a migration to the migrator
func (m *Migrator) AddMigration(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// Migrate runs all pending migrations
func (m *Migrator) Migrate() error {
	// Create migrations table if it doesn't exist
	if err := m.db.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	var appliedMigrations []MigrationRecord
	if err := m.db.Find(&appliedMigrations).Error; err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Convert to map for easier lookup
	appliedMap := make(map[string]bool)
	for _, migration := range appliedMigrations {
		appliedMap[migration.ID] = true
	}

	// Run pending migrations
	for _, migration := range m.migrations {
		id := migration.ID()
		if !appliedMap[id] {
			m.logger.Info("Running migration", zap.String("migration_id", id))

			// Start transaction
			tx := m.db.Begin()
			if tx.Error != nil {
				return fmt.Errorf("failed to begin transaction: %w", tx.Error)
			}

			// Run migration
			if err := migration.Migrate(tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("migration %s failed: %w", id, err)
			}

			// Record migration
			if err := tx.Create(&MigrationRecord{ID: id}).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to record migration %s: %w", id, err)
			}

			// Commit transaction
			if err := tx.Commit().Error; err != nil {
				return fmt.Errorf("failed to commit migration %s: %w", id, err)
			}

			m.logger.Info("Migration successful", zap.String("migration_id", id))
		}
	}

	return nil
}

// Rollback rolls back the last n migrations
func (m *Migrator) Rollback(n int) error {
	// Get applied migrations in reverse order
	var appliedMigrations []MigrationRecord
	if err := m.db.Order("applied_at DESC").Limit(n).Find(&appliedMigrations).Error; err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Map migrations by ID
	migrationsMap := make(map[string]Migration)
	for _, migration := range m.migrations {
		migrationsMap[migration.ID()] = migration
	}

	// Rollback migrations
	for _, record := range appliedMigrations {
		migration, exists := migrationsMap[record.ID]
		if !exists {
			m.logger.Warn("Migration not found, skipping rollback", zap.String("migration_id", record.ID))
			continue
		}

		m.logger.Info("Rolling back migration", zap.String("migration_id", record.ID))

		// Start transaction
		tx := m.db.Begin()
		if tx.Error != nil {
			return fmt.Errorf("failed to begin transaction: %w", tx.Error)
		}

		// Rollback migration
		if err := migration.Rollback(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("rollback of %s failed: %w", record.ID, err)
		}

		// Remove migration record
		if err := tx.Delete(&MigrationRecord{ID: record.ID}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to remove migration record %s: %w", record.ID, err)
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit rollback of %s: %w", record.ID, err)
		}

		m.logger.Info("Rollback successful", zap.String("migration_id", record.ID))
	}

	return nil
}
