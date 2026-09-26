package eligibility

import (
	"errors"

	"gorm.io/gorm"
)

var ErrEligibilityNotFound = errors.New("student eligibility not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create crée une nouvelle éligibilité.
func (r *Repository) Create(eligibility *Eligibility) error {
	if err := r.db.Create(eligibility).Error; err != nil {
		return err
	}

	return nil
}

// GetByStudentID récupère l'éligibilité d'un étudiant.
func (r *Repository) GetByStudentID(studentID int) (*Eligibility, error) {
	var eligibility Eligibility

	result := r.db.
		Where("student_id = ?", studentID).
		First(&eligibility)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrEligibilityNotFound
		}

		return nil, result.Error
	}

	return &eligibility, nil
}

// Update modifie l'éligibilité d'un étudiant.
func (r *Repository) Update(eligibility *Eligibility) error {
	result := r.db.
		Model(&Eligibility{}).
		Where("student_id = ?", eligibility.StudentID).
		Updates(map[string]any{
			"enrollment":      eligibility.Enrollment,
			"academic_fees":   eligibility.AcademicFees,
			"laboratory_fees": eligibility.LaboratoryFees,
			"access_fees":     eligibility.AccessFees,
			"eligible":        eligibility.Eligible,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrEligibilityNotFound
	}

	return nil
}

// Delete supprime l'éligibilité d'un étudiant.
func (r *Repository) Delete(studentID int) error {
	result := r.db.
		Where("student_id = ?", studentID).
		Delete(&Eligibility{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrEligibilityNotFound
	}

	return nil
}