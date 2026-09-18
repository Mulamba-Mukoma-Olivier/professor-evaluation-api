package results

import (
	"gorm.io/gorm"

	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByProfessor(
	professorID int,
	courseID int,
	academicYear string,
	period string,
) ([]evaluations.Evaluation, error) {

	var evaluationList []evaluations.Evaluation

	result := r.db.
		Preload("Answers").
		Where(
			"professor_id = ? AND course_id = ? AND academic_year = ? AND period = ?",
			professorID,
			courseID,
			academicYear,
			period,
		).
		Find(&evaluationList)

	if result.Error != nil {
		return nil, result.Error
	}

	return evaluationList, nil
}