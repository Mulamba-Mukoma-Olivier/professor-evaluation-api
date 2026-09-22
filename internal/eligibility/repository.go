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

func (r *Repository) GetByStudentID(studentID int) (*Eligibility, error) {
	var eligibility Eligibility

	result := r.db.
		Where("student_id = ?", studentID).
		First(&eligibility)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrEligibilityNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &eligibility, nil
}
