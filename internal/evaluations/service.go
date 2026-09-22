package evaluations

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidStudentID   = errors.New("invalid student ID")
	ErrInvalidProfessorID = errors.New("invalid professor ID")
	ErrInvalidCourseID    = errors.New("invalid course ID")
	ErrInvalidEvaluationID = errors.New("invalid evaluation ID")
	ErrAcademicYearRequired = errors.New("academic year is required")
	ErrPeriodRequired       = errors.New("evaluation period is required")
	ErrAnswersRequired      = errors.New("evaluation must contain at least one answer")
	ErrInvalidCriterionID   = errors.New("invalid criterion ID")
	ErrInvalidScore         = errors.New("score must be between 1 and 5")
	ErrDuplicateEvaluation = errors.New("student has already evaluated this professor for this course")
)

type EvaluationRepository interface {
	Create(evaluation Evaluation) (*Evaluation, error)
	GetAll() ([]Evaluation, error)
	GetByID(id int) (*Evaluation, error)
	Exists(
		studentID int,
		professorID int,
		courseID int,
		academicYear string,
		period string,
	) (bool, error)
	Delete(id int) error
}

type Service struct {
	repository EvaluationRepository
}

func NewService(repository EvaluationRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	studentID int,
	request CreateEvaluationRequest,
) (*Evaluation, error) {

	if studentID <= 0 {
		return nil, ErrInvalidStudentID
	}

	if request.ProfessorID <= 0 {
		return nil, ErrInvalidProfessorID
	}

	if request.CourseID <= 0 {
		return nil, ErrInvalidCourseID
	}

	request.AcademicYear = strings.TrimSpace(request.AcademicYear)
	if request.AcademicYear == "" {
		return nil, ErrAcademicYearRequired
	}

	request.Period = strings.TrimSpace(request.Period)
	if request.Period == "" {
		return nil, ErrPeriodRequired
	}

	if len(request.Answers) == 0 {
		return nil, ErrAnswersRequired
	}

	for _, answer := range request.Answers {
		if answer.CriterionID <= 0 {
			return nil, ErrInvalidCriterionID
		}

		if answer.Score < 1 || answer.Score > 5 {
			return nil, ErrInvalidScore
		}
	}

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
		return nil, ErrDuplicateEvaluation
	}

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
		return nil, ErrInvalidEvaluationID
	}

	return s.repository.GetByID(id)
}
