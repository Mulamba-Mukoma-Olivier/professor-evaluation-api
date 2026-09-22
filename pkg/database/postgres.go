package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// NewPostgresDB crée une connexion GORM vers PostgreSQL.
func NewPostgresDB(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to database: %w",
			err,
		)
	}

	// Récupérer la connexion SQL sous-jacente.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get database instance: %w",
			err,
		)
	}

	// Configuration du pool de connexions.
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	// Vérifier que PostgreSQL est réellement accessible.
	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()

		return nil, fmt.Errorf(
			"failed to ping database: %w",
			err,
		)
	}

	return db, nil
}
