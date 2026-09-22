package criteria

import (
	"errors"
	"strings"
)

var (
	ErrInvalidCriterionID       = errors.New("invalid criterion ID")
	ErrCriterionNameRequired    = errors.New("criterion name is required")
	ErrMaxScoreInvalid          = errors.New("max score must be greater than zero")
)

type CriterionRepository interface {
	GetAll() ([]Criterion, error)
	GetActive() ([]Criterion, error)
	GetByID(id int) (*Criterion, error)
	Create(criterion Criterion) (*Criterion, error)
	Update(id int, criterion Criterion) (*Criterion, error)
	Delete(id int) error
}

type Service struct {
	repository CriterionRepository
}

func NewService(repository CriterionRepository) *Service {
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
		return nil, ErrInvalidCriterionID
	}

	return s.repository.GetByID(id)
}

func (s *Service) Create(request CreateCriterionRequest) (*Criterion, error) {
	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)

	if request.Name == "" {
		return nil, ErrCriterionNameRequired
	}

	if request.MaxScore <= 0 {
		return nil, ErrMaxScoreInvalid
	}

	criterion := Criterion{
		Name:        request.Name,
		Description: request.Description,
		MaxScore:    request.MaxScore,
		Active:      true,
	}

	return s.repository.Create(criterion)
}

func (s *Service) Update(
	id int,
	request UpdateCriterionRequest,
) (*Criterion, error) {
	if id <= 0 {
		return nil, ErrInvalidCriterionID
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)

	if request.Name == "" {
		return nil, ErrCriterionNameRequired
	}

	if request.MaxScore <= 0 {
		return nil, ErrMaxScoreInvalid
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
		return ErrInvalidCriterionID
	}

	return s.repository.Delete(id)
}
