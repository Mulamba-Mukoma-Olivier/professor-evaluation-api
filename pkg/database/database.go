package database

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewSQLiteDB crée une connexion GORM vers SQLite.
func NewSQLiteDB() (*gorm.DB, error) {
	db, err := gorm.Open(
		sqlite.Open("professor_evaluation.db"),
		&gorm.Config{},
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to SQLite: %w",
			err,
		)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get database instance: %w",
			err,
		)
	}

	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()

		return nil, fmt.Errorf(
			"failed to ping SQLite: %w",
			err,
		)
	}

	return db, nil
}