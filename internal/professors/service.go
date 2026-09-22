package professors

import (
	"errors"
	"strings"
)

var (
	ErrInvalidProfessorID       = errors.New("invalid professor ID")
	ErrProfessorMatriculeRequired = errors.New("professor matricule is required")
	ErrFirstNameRequired        = errors.New("first name is required")
	ErrLastNameRequired         = errors.New("last name is required")
	ErrDepartmentRequired       = errors.New("department is required")
	ErrInvalidProfessorStatus   = errors.New("invalid professor status")
)

type ProfessorRepository interface {
	GetAll() ([]Professor, error)
	GetAllPaginated(page, pageSize int) ([]Professor, int, error)
	GetActive() ([]Professor, error)
	GetByStatus(status string) ([]Professor, error)
	GetByID(id int) (*Professor, error)
	Create(professor Professor) (*Professor, error)
	Update(id int, professor Professor) (*Professor, error)
	Delete(id int) error
}

type Service struct {
	repository ProfessorRepository
}

func NewService(repository ProfessorRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAll() ([]Professor, error) {
	return s.repository.GetAll()
}

func (s *Service) GetAllPaginated(page, pageSize int) ([]Professor, int, error) {
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	return s.repository.GetAllPaginated(page, pageSize)
}

func (s *Service) GetActive() ([]Professor, error) {
	return s.repository.GetActive()
}

func (s *Service) GetByStatus(status string) ([]Professor, error) {
	status = strings.TrimSpace(status)

	if status == "" {
		return nil, ErrInvalidProfessorStatus
	}

	if !isValidStatus(status) {
		return nil, ErrInvalidProfessorStatus
	}

	return s.repository.GetByStatus(status)
}

func (s *Service) GetByID(id int) (*Professor, error) {
	if id <= 0 {
		return nil, ErrInvalidProfessorID
	}

	return s.repository.GetByID(id)
}

func (s *Service) Create(
	request CreateProfessorRequest,
) (*Professor, error) {

	request.Matricule = strings.TrimSpace(request.Matricule)
	request.FirstName = strings.TrimSpace(request.FirstName)
	request.LastName = strings.TrimSpace(request.LastName)
	request.Email = strings.TrimSpace(request.Email)
	request.Department = strings.TrimSpace(request.Department)
	request.Grade = strings.TrimSpace(request.Grade)
	request.Status = strings.TrimSpace(request.Status)

	if request.Matricule == "" {
		return nil, ErrProfessorMatriculeRequired
	}

	if request.FirstName == "" {
		return nil, ErrFirstNameRequired
	}

	if request.LastName == "" {
		return nil, ErrLastNameRequired
	}

	if request.Department == "" {
		return nil, ErrDepartmentRequired
	}

	status := request.Status

	if status == "" {
		status = "active"
	}

	if !isValidStatus(status) {
		return nil, ErrInvalidProfessorStatus
	}

	professor := Professor{
		Matricule:  request.Matricule,
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		Email:      request.Email,
		Department: request.Department,
		Grade:      request.Grade,
		Active:     status == "active",
		Status:     status,
	}

	return s.repository.Create(professor)
}

func (s *Service) Update(
	id int,
	request UpdateProfessorRequest,
) (*Professor, error) {

	if id <= 0 {
		return nil, ErrInvalidProfessorID
	}

	request.Matricule = strings.TrimSpace(request.Matricule)
	request.FirstName = strings.TrimSpace(request.FirstName)
	request.LastName = strings.TrimSpace(request.LastName)
	request.Email = strings.TrimSpace(request.Email)
	request.Department = strings.TrimSpace(request.Department)
	request.Grade = strings.TrimSpace(request.Grade)
	request.Status = strings.TrimSpace(request.Status)

	if request.Matricule == "" {
		return nil, ErrProfessorMatriculeRequired
	}

	if request.FirstName == "" {
		return nil, ErrFirstNameRequired
	}

	if request.LastName == "" {
		return nil, ErrLastNameRequired
	}

	if request.Department == "" {
		return nil, ErrDepartmentRequired
	}

	status := request.Status

	if status == "" {
		if request.Active {
			status = "active"
		} else {
			status = "inactive"
		}
	}

	if !isValidStatus(status) {
		return nil, ErrInvalidProfessorStatus
	}

	professor := Professor{
		Matricule:  request.Matricule,
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		Email:      request.Email,
		Department: request.Department,
		Grade:      request.Grade,
		Active:     request.Active,
		Status:     status,
	}

	return s.repository.Update(id, professor)
}

func (s *Service) Delete(id int) error {
	if id <= 0 {
		return ErrInvalidProfessorID
	}

	return s.repository.Delete(id)
}

func isValidStatus(status string) bool {
	switch status {
	case "active", "inactive", "on_leave", "retired":
		return true
	default:
		return false
	}
}