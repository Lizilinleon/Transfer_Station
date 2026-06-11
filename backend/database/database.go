package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"ai-chat-platform/backend/config"
	"ai-chat-platform/backend/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// New opens the SQLite database, migrates tables, and prepares default records.
func New(cfg config.Config) (*gorm.DB, error) {
	if err := os.MkdirAll(filepath.Dir(cfg.SQLitePath), 0o755); err != nil {
		return nil, fmt.Errorf("create sqlite directory: %w", err)
	}

	// Use a quieter GORM logger that still reports slow SQL and warnings.
	dbLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             500 * time.Millisecond,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(sqlite.Open(cfg.SQLitePath), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	// AutoMigrate keeps the local SQLite schema aligned with model definitions.
	if err := db.AutoMigrate(
		&model.User{},
		&model.AIModel{},
		&model.ProviderChannel{},
		&model.ModelAbility{},
		&model.APIKey{},
		&model.UsageLog{},
		&model.SystemOption{},
		&model.ChatSession{},
		&model.ChatMessage{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate tables: %w", err)
	}

	if err := seedDefaults(db, cfg); err != nil {
		return nil, fmt.Errorf("seed defaults: %w", err)
	}

	return db, nil
}
