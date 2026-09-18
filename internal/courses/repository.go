package courses

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

func (r *Repository) GetAll() ([]Course, error) {
	var courses []Course

	result := r.db.Find(&courses)

	if result.Error != nil {
		return nil, result.Error
	}

	return courses, nil
}

func (r *Repository) GetByID(id int) (*Course, error) {
	var course Course

	result := r.db.First(&course, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("course not found")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &course, nil
}

func (r *Repository) Create(course Course) (*Course, error) {
	result := r.db.Create(&course)

	if result.Error != nil {
		return nil, result.Error
	}

	return &course, nil
}

func (r *Repository) Update(id int, course Course) (*Course, error) {
	var existingCourse Course

	result := r.db.First(&existingCourse, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("course not found")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	existingCourse.Code = course.Code
	existingCourse.Name = course.Name
	existingCourse.Description = course.Description
	existingCourse.Department = course.Department
	existingCourse.AcademicYear = course.AcademicYear

	result = r.db.Save(&existingCourse)

	if result.Error != nil {
		return nil, result.Error
	}

	return &existingCourse, nil
}

func (r *Repository) Delete(id int) error {
	result := r.db.Delete(&Course{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("course not found")
	}

	return nil
}