package courses

import (
	"errors"
	"strings"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAll() ([]Course, error) {
	return s.repository.GetAll()
}

func (s *Service) GetByID(id int) (*Course, error) {
	if id <= 0 {
		return nil, errors.New("invalid course ID")
	}

	return s.repository.GetByID(id)
}

func (s *Service) Create(request CreateCourseRequest) (*Course, error) {
	request.Code = strings.TrimSpace(request.Code)
	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)
	request.Department = strings.TrimSpace(request.Department)
	request.AcademicYear = strings.TrimSpace(request.AcademicYear)

	if request.Code == "" {
		return nil, errors.New("course code is required")
	}

	if request.Name == "" {
		return nil, errors.New("course name is required")
	}

	if request.Department == "" {
		return nil, errors.New("department is required")
	}

	if request.AcademicYear == "" {
		return nil, errors.New("academic year is required")
	}

	course := Course{
		Code:         request.Code,
		Name:         request.Name,
		Description:  request.Description,
		Department:   request.Department,
		AcademicYear: request.AcademicYear,
	}

	return s.repository.Create(course)
}

func (s *Service) Update(
	id int,
	request UpdateCourseRequest,
) (*Course, error) {
	if id <= 0 {
		return nil, errors.New("invalid course ID")
	}

	request.Code = strings.TrimSpace(request.Code)
	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)
	request.Department = strings.TrimSpace(request.Department)
	request.AcademicYear = strings.TrimSpace(request.AcademicYear)

	if request.Code == "" {
		return nil, errors.New("course code is required")
	}

	if request.Name == "" {
		return nil, errors.New("course name is required")
	}

	if request.Department == "" {
		return nil, errors.New("department is required")
	}

	if request.AcademicYear == "" {
		return nil, errors.New("academic year is required")
	}

	course := Course{
		Code:         request.Code,
		Name:         request.Name,
		Description:  request.Description,
		Department:   request.Department,
		AcademicYear: request.AcademicYear,
	}

	return s.repository.Update(id, course)
}

func (s *Service) Delete(id int) error {
	if id <= 0 {
		return errors.New("invalid course ID")
	}

	return s.repository.Delete(id)
}
