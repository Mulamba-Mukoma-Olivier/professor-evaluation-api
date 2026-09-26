package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// Migrate crée ou met à jour les tables avec GORM.
func Migrate(db *gorm.DB, models ...interface{}) error {
	log.Println("Starting database migration...")

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf(
			"failed to migrate database: %w",
			err,
		)
	}

	log.Println("Database migration completed successfully")

	return nil
}