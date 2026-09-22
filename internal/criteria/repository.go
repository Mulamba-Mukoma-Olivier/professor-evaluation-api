package criteria

import (
	"errors"

	"gorm.io/gorm"
)

var ErrCriterionNotFound = errors.New("criterion not found")

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

	if err := r.db.Find(&criteria).Error; err != nil {
		return nil, err
	}

	return criteria, nil
}

func (r *Repository) GetActive() ([]Criterion, error) {
	var criteria []Criterion

	if err := r.db.
		Where("active = ?", true).
		Find(&criteria).Error; err != nil {
		return nil, err
	}

	return criteria, nil
}

func (r *Repository) GetByID(id int) (*Criterion, error) {
	var criterion Criterion

	result := r.db.First(&criterion, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrCriterionNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &criterion, nil
}

func (r *Repository) Create(criterion Criterion) (*Criterion, error) {
	if err := r.db.Create(&criterion).Error; err != nil {
		return nil, err
	}

	return &criterion, nil
}

func (r *Repository) Update(id int, criterion Criterion) (*Criterion, error) {
	var existingCriterion Criterion

	result := r.db.First(&existingCriterion, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrCriterionNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	existingCriterion.Name = criterion.Name
	existingCriterion.Description = criterion.Description
	existingCriterion.MaxScore = criterion.MaxScore
	existingCriterion.Active = criterion.Active

	if err := r.db.Save(&existingCriterion).Error; err != nil {
		return nil, err
	}

	return &existingCriterion, nil
}

func (r *Repository) Delete(id int) error {
	result := r.db.Delete(&Criterion{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrCriterionNotFound
	}

	return nil
}
