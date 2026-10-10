package database

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewDB ouvre PostgreSQL lorsque DB_DRIVER=postgres, sinon utilise SQLite.
func NewDB() (*gorm.DB, error) {
	if os.Getenv("DB_DRIVER") == "postgres" {
		return NewPostgresDB()
	}
	return NewSQLiteDB()
}

// NewPostgresDB ouvre PostgreSQL avec une URL fournie dans DATABASE_URL.
func NewPostgresDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is required when DB_DRIVER=postgres")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	return configureConnectionPool(db, 5)
}

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

	return configureConnectionPool(db, 1)
}

func configureConnectionPool(db *gorm.DB, maxOpenConns int) (*gorm.DB, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get database instance: %w",
			err,
		)
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()

		return nil, fmt.Errorf(
			"failed to ping database: %w",
			err,
		)
	}

	return db, nil
}
