package evaluations

import (
	"errors"

	"gorm.io/gorm"
)

var ErrEvaluationNotFound = errors.New("evaluation not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(evaluation Evaluation) (*Evaluation, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Création de l'évaluation.
		if err := tx.Create(&evaluation).Error; err != nil {
			return err
		}

		// Création des réponses.
		for i := range evaluation.Answers {
			evaluation.Answers[i].EvaluationID = evaluation.ID

			if err := tx.Create(&evaluation.Answers[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Recharge l'évaluation avec ses réponses.
	if err := r.db.
		Preload("Answers").
		First(&evaluation, evaluation.ID).Error; err != nil {
		return nil, err
	}

	return &evaluation, nil
}

func (r *Repository) GetAll() ([]Evaluation, error) {
	var evaluations []Evaluation

	result := r.db.
		Preload("Answers").
		Find(&evaluations)

	if result.Error != nil {
		return nil, result.Error
	}

	return evaluations, nil
}

func (r *Repository) GetByID(id int) (*Evaluation, error) {
	var evaluation Evaluation

	result := r.db.
		Preload("Answers").
		First(&evaluation, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrEvaluationNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &evaluation, nil
}

func (r *Repository) Exists(
	studentID int,
	professorID int,
	courseID int,
	academicYear string,
	period string,
) (bool, error) {
	var count int64

	result := r.db.
		Model(&Evaluation{}).
		Where(
			"student_id = ? AND professor_id = ? AND course_id = ? AND academic_year = ? AND period = ?",
			studentID,
			professorID,
			courseID,
			academicYear,
			period,
		).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func (r *Repository) Delete(id int) error {
	result := r.db.Delete(&Evaluation{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrEvaluationNotFound
	}

	return nil
}
