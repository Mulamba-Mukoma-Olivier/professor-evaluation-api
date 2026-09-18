package criteria

import "errors"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAll() ([]Criterion, error) {
	return s.repository.GetAll()
}

func (s *Service) GetActive() ([]Criterion, error) {
	return s.repository.GetActive()
}

func (s *Service) GetByID(id int) (*Criterion, error) {
	if id <= 0 {
		return nil, errors.New("invalid criterion ID")
	}

	return s.repository.GetByID(id)
}

func (s *Service) Create(request CreateCriterionRequest) (*Criterion, error) {
	if request.Name == "" {
		return nil, errors.New("criterion name is required")
	}

	if request.MaxScore <= 0 {
		return nil, errors.New("max score must be greater than zero")
	}

	criterion := Criterion{
		Name:        request.Name,
		Description: request.Description,
		MaxScore:    request.MaxScore,
		Active:      true,
	}

	return s.repository.Create(criterion)
}

func (s *Service) Update(id int, request UpdateCriterionRequest) (*Criterion, error) {
	if id <= 0 {
		return nil, errors.New("invalid criterion ID")
	}

	if request.Name == "" {
		return nil, errors.New("criterion name is required")
	}

	if request.MaxScore <= 0 {
		return nil, errors.New("max score must be greater than zero")
	}

	criterion := Criterion{
		Name:        request.Name,
		Description: request.Description,
		MaxScore:    request.MaxScore,
		Active:      request.Active,
	}

	return s.repository.Update(id, criterion)
}

func (s *Service) Delete(id int) error {
	if id <= 0 {
		return errors.New("invalid criterion ID")
	}

	return s.repository.Delete(id)
}