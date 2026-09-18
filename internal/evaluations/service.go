package evaluations

import (
	"errors"
	"time"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	studentID int,
	request CreateEvaluationRequest,
) (*Evaluation, error) {

	if studentID <= 0 {
		return nil, errors.New("invalid student ID")
	}

	if request.ProfessorID <= 0 {
		return nil, errors.New("invalid professor ID")
	}

	if request.CourseID <= 0 {
		return nil, errors.New("invalid course ID")
	}

	if request.AcademicYear == "" {
		return nil, errors.New("academic year is required")
	}

	if request.Period == "" {
		return nil, errors.New("evaluation period is required")
	}

	if len(request.Answers) == 0 {
		return nil, errors.New("evaluation must contain at least one answer")
	}

	// Vérification des scores
	for _, answer := range request.Answers {
		if answer.CriterionID <= 0 {
			return nil, errors.New("invalid criterion ID")
		}

		if answer.Score < 1 || answer.Score > 5 {
			return nil, errors.New("score must be between 1 and 5")
		}
	}

	// Empêcher une double évaluation
	exists, err := s.repository.Exists(
		studentID,
		request.ProfessorID,
		request.CourseID,
		request.AcademicYear,
		request.Period,
	)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New(
			"student has already evaluated this professor for this course",
		)
	}

	// Conversion des réponses du DTO vers les modèles GORM
	answers := make([]EvaluationAnswer, 0, len(request.Answers))

	for _, answer := range request.Answers {
		answers = append(answers, EvaluationAnswer{
			CriterionID: answer.CriterionID,
			Score:       answer.Score,
		})
	}

	evaluation := Evaluation{
		StudentID:    studentID,
		ProfessorID:  request.ProfessorID,
		CourseID:     request.CourseID,
		AcademicYear: request.AcademicYear,
		Period:       request.Period,
		Answers:      answers,
		SubmittedAt:  time.Now(),
	}

	return s.repository.Create(evaluation)
}

func (s *Service) GetAll() ([]Evaluation, error) {
	return s.repository.GetAll()
}

func (s *Service) GetByID(id int) (*Evaluation, error) {
	if id <= 0 {
		return nil, errors.New("invalid evaluation ID")
	}

	return s.repository.GetByID(id)
}