package database

import (
	"context"
	"drive/internal/config"
	"drive/internal/util"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// CustomGormLogger implements gorm logger.Interface
type CustomGormLogger struct {
	logger *util.Logger
}

// LogMode sets the log mode
func (l *CustomGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

// Info logs info messages
func (l *CustomGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Info(msg, zap.Any("data", data))
}

// Warn logs warn messages
func (l *CustomGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Warn(msg, zap.Any("data", data))
}

// Error logs error messages
func (l *CustomGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Error(msg, zap.Any("data", data))
}

// Trace logs SQL queries
func (l *CustomGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if err != nil {
		sql, rows := fc()
		l.logger.Error("SQL query failed",
			zap.Error(err),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("duration", time.Since(begin)),
		)
		return
	}

	sql, rows := fc()
	l.logger.Debug("SQL query executed",
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("duration", time.Since(begin)),
	)
}

// InitDatabase initializes the database connection
func InitDatabase(cfg *config.Config, logger *util.Logger) (*gorm.DB, error) {
	// Build DSN string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	// Create custom logger
	gormLogger := &CustomGormLogger{logger: logger}

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info("Successfully connected to database")

	return db, nil
}
