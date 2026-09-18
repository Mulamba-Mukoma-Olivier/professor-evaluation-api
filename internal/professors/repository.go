package professors

import (
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll() ([]Professor, error) {
	var professors []Professor

	result := r.db.Find(&professors)

	if result.Error != nil {
		return nil, result.Error
	}

	return professors, nil
}

func (r *Repository) GetActive() ([]Professor, error) {
	var professors []Professor

	result := r.db.
		Where("active = ?", true).
		Find(&professors)

	if result.Error != nil {
		return nil, result.Error
	}

	return professors, nil
}

func (r *Repository) GetByID(id int) (*Professor, error) {
	var professor Professor

	result := r.db.First(&professor, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("professor not found")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &professor, nil
}

func (r *Repository) Create(professor Professor) (*Professor, error) {
	result := r.db.Create(&professor)

	if result.Error != nil {
		return nil, result.Error
	}

	return &professor, nil
}

func (r *Repository) Update(
	id int,
	professor Professor,
) (*Professor, error) {

	var existingProfessor Professor

	result := r.db.First(&existingProfessor, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("professor not found")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	existingProfessor.Matricule = professor.Matricule
	existingProfessor.FirstName = professor.FirstName
	existingProfessor.LastName = professor.LastName
	existingProfessor.Email = professor.Email
	existingProfessor.Department = professor.Department
	existingProfessor.Grade = professor.Grade
	existingProfessor.Active = professor.Active

	result = r.db.Save(&existingProfessor)

	if result.Error != nil {
		return nil, result.Error
	}

	return &existingProfessor, nil
}

func (r *Repository) Delete(id int) error {
	result := r.db.Delete(&Professor{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("professor not found")
	}

	return nil
}