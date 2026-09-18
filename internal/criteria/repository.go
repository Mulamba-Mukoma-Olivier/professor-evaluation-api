package criteria

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

func (r *Repository) GetAll() ([]Criterion, error) {
	var criteria []Criterion

	result := r.db.Find(&criteria)

	if result.Error != nil {
		return nil, result.Error
	}

	return criteria, nil
}

func (r *Repository) GetActive() ([]Criterion, error) {
	var criteria []Criterion

	result := r.db.
		Where("active = ?", true).
		Find(&criteria)

	if result.Error != nil {
		return nil, result.Error
	}

	return criteria, nil
}

func (r *Repository) GetByID(id int) (*Criterion, error) {
	var criterion Criterion

	result := r.db.First(&criterion, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("criterion not found")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &criterion, nil
}

func (r *Repository) Create(criterion Criterion) (*Criterion, error) {
	result := r.db.Create(&criterion)

	if result.Error != nil {
		return nil, result.Error
	}

	return &criterion, nil
}

func (r *Repository) Update(id int, criterion Criterion) (*Criterion, error) {
	var existingCriterion Criterion

	result := r.db.First(&existingCriterion, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("criterion not found")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	existingCriterion.Name = criterion.Name
	existingCriterion.Description = criterion.Description
	existingCriterion.MaxScore = criterion.MaxScore
	existingCriterion.Active = criterion.Active

	result = r.db.Save(&existingCriterion)

	if result.Error != nil {
		return nil, result.Error
	}

	return &existingCriterion, nil
}

func (r *Repository) Delete(id int) error {
	result := r.db.Delete(&Criterion{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("criterion not found")
	}

	return nil
}