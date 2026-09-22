package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gorm.io/gorm"
)

// Migrate applique les migrations SQL de la base de données.
func Migrate(db *gorm.DB, config Config) error {
	log.Println("Starting database migration...")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.SSLMode,
	)

	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	defer func() {
		sourceErr, databaseErr := m.Close()

		if sourceErr != nil {
			log.Printf("migration source close error: %v", sourceErr)
		}

		if databaseErr != nil {
			log.Printf("migration database close error: %v", databaseErr)
		}
	}()

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("Database is already up to date")
	} else {
		log.Println("Database migrations applied successfully")
	}

	return nil
}
