package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"

	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/auth"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/courses"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/criteria"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/eligibility"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/professors"
)

func RunMigrations(cfg Config) error {
	log.Println("Starting database migrations...")

	// Build connection string for migrate
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	// Create migrate instance
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}

	return nil
}

// Migrate is kept for backward compatibility but uses AutoMigrate
// This is deprecated in favor of RunMigrations with golang-migrate
func Migrate(db *gorm.DB) error {
	log.Println("Starting database migration (AutoMigrate)...")
	
	err := db.AutoMigrate(
		&auth.User{},
		&courses.Course{},
		&criteria.Criterion{},
		&eligibility.Eligibility{},
		&professors.Professor{},
		&evaluations.Evaluation{},
		&evaluations.EvaluationAnswer{},
	)
	
	if err != nil {
		return err
	}
	
	log.Println("Database migration completed successfully")
	return nil
}