package professors

import "errors"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAll() ([]Professor, error) {
	return s.repository.GetAll()
}

func (s *Service) GetActive() ([]Professor, error) {
	return s.repository.GetActive()
}

func (s *Service) GetByID(id int) (*Professor, error) {

	if id <= 0 {
		return nil, errors.New("invalid professor ID")
	}

	return s.repository.GetByID(id)
}

func (s *Service) Create(
	request CreateProfessorRequest,
) (*Professor, error) {

	if request.Matricule == "" {
		return nil, errors.New(
			"professor matricule is required",
		)
	}

	if request.FirstName == "" {
		return nil, errors.New(
			"first name is required",
		)
	}

	if request.LastName == "" {
		return nil, errors.New(
			"last name is required",
		)
	}

	if request.Department == "" {
		return nil, errors.New(
			"department is required",
		)
	}

	professor := Professor{
		Matricule:  request.Matricule,
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		Email:      request.Email,
		Department: request.Department,
		Grade:      request.Grade,
		Active:     true,
	}

	return s.repository.Create(professor)
}

func (s *Service) Update(
	id int,
	request UpdateProfessorRequest,
) (*Professor, error) {

	if id <= 0 {
		return nil, errors.New("invalid professor ID")
	}

	if request.Matricule == "" {
		return nil, errors.New(
			"professor matricule is required",
		)
	}

	if request.FirstName == "" {
		return nil, errors.New(
			"first name is required",
		)
	}

	if request.LastName == "" {
		return nil, errors.New(
			"last name is required",
		)
	}

	if request.Department == "" {
		return nil, errors.New(
			"department is required",
		)
	}

	professor := Professor{
		Matricule:  request.Matricule,
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		Email:      request.Email,
		Department: request.Department,
		Grade:      request.Grade,
		Active:     request.Active,
	}

	return s.repository.Update(id, professor)
}

func (s *Service) Delete(id int) error {

	if id <= 0 {
		return errors.New("invalid professor ID")
	}

	return s.repository.Delete(id)
}