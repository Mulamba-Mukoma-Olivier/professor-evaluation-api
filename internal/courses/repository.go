package courses

import (
	"errors"

	"gorm.io/gorm"
)

var ErrCourseNotFound = errors.New("course not found")

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

	if err := r.db.Find(&courses).Error; err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *Repository) GetByID(id int) (*Course, error) {
	var course Course

	result := r.db.First(&course, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrCourseNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &course, nil
}

func (r *Repository) Create(course Course) (*Course, error) {
	if err := r.db.Create(&course).Error; err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *Repository) Update(id int, course Course) (*Course, error) {
	var existingCourse Course

	result := r.db.First(&existingCourse, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrCourseNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	existingCourse.Code = course.Code
	existingCourse.Name = course.Name
	existingCourse.Description = course.Description
	existingCourse.Department = course.Department
	existingCourse.AcademicYear = course.AcademicYear

	if err := r.db.Save(&existingCourse).Error; err != nil {
		return nil, err
	}

	return &existingCourse, nil
}

func (r *Repository) Delete(id int) error {
	result := r.db.Delete(&Course{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrCourseNotFound
	}

	return nil
}
