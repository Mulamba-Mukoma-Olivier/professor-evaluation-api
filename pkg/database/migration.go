package database

import (
	"gorm.io/gorm"

	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/auth"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/courses"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/criteria"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/eligibility"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/professors"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&auth.User{},
		&courses.Course{},
		&criteria.Criterion{},
		&eligibility.Eligibility{},
		&professors.Professor{},
		&evaluations.Evaluation{},
		&evaluations.EvaluationAnswer{},
	)
}