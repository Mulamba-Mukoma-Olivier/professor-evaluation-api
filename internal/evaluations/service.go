package evaluations

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidStudentID    = errors.New("invalid student ID")
	ErrInvalidProfessorID  = errors.New("invalid professor ID")
	ErrInvalidCourseID     = errors.New("invalid course ID")
	ErrInvalidEvaluationID = errors.New("invalid evaluation ID")

	ErrAcademicYearRequired = errors.New("academic year is required")
	ErrPeriodRequired       = errors.New("evaluation period is required")
	ErrAnswersRequired      = errors.New("evaluation must contain at least one answer")

	ErrInvalidCriterionID = errors.New("invalid criterion ID")
	ErrInvalidScore       = errors.New("score must be between 1 and 5")

	ErrDuplicateCriterion = errors.New("criterion cannot appear more than once in an evaluation")
	ErrDuplicateEvaluation = errors.New(
		"student has already evaluated this professor for this course, academic year and period",
	)
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

// Create crée une nouvelle évaluation après validation
// des données reçues et des règles métier de base.
func (s *Service) Create(
	studentID int,
	request CreateEvaluationRequest,
) (*Evaluation, error) {

	// --------------------------------------------------
	// Validation de l'étudiant
	// --------------------------------------------------
	if studentID <= 0 {
		return nil, ErrInvalidStudentID
	}

	// --------------------------------------------------
	// Validation du professeur
	// --------------------------------------------------
	if request.ProfessorID <= 0 {
		return nil, ErrInvalidProfessorID
	}

	// --------------------------------------------------
	// Validation du cours
	// --------------------------------------------------
	if request.CourseID <= 0 {
		return nil, ErrInvalidCourseID
	}

	// --------------------------------------------------
	// Validation de l'année académique
	// --------------------------------------------------
	request.AcademicYear = strings.TrimSpace(request.AcademicYear)

	if request.AcademicYear == "" {
		return nil, ErrAcademicYearRequired
	}

	// --------------------------------------------------
	// Validation de la période
	// --------------------------------------------------
	request.Period = strings.TrimSpace(request.Period)

	if request.Period == "" {
		return nil, ErrPeriodRequired
	}

	// --------------------------------------------------
	// Validation des réponses
	// --------------------------------------------------
	if len(request.Answers) == 0 {
		return nil, ErrAnswersRequired
	}

	// Permet de détecter deux fois le même critère.
	criteria := make(map[int]struct{}, len(request.Answers))

	for _, answer := range request.Answers {

		// Vérification du critère.
		if answer.CriterionID <= 0 {
			return nil, ErrInvalidCriterionID
		}

		// Vérification de la note.
		if answer.Score < 1 || answer.Score > 5 {
			return nil, ErrInvalidScore
		}

		// Vérification des critères en double.
		if _, exists := criteria[answer.CriterionID]; exists {
			return nil, ErrDuplicateCriterion
		}

		criteria[answer.CriterionID] = struct{}{}
	}

	// --------------------------------------------------
	// Vérification d'une éventuelle évaluation existante
	// --------------------------------------------------
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

	// --------------------------------------------------
	// Transformation DTO -> Model
	// --------------------------------------------------
	answers := make([]EvaluationAnswer, 0, len(request.Answers))

	for _, answer := range request.Answers {
		answers = append(answers, EvaluationAnswer{
			CriterionID: answer.CriterionID,
			Score:       answer.Score,
		})
	}

	// --------------------------------------------------
	// Création de l'évaluation
	// --------------------------------------------------
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

// GetAll retourne toutes les évaluations.
func (s *Service) GetAll() ([]Evaluation, error) {
	return s.repository.GetAll()
}

// GetByID retourne une évaluation par son ID.
func (s *Service) GetByID(id int) (*Evaluation, error) {
	if id <= 0 {
		return nil, ErrInvalidEvaluationID
	}

	return s.repository.GetByID(id)
}

// Delete supprime une évaluation.
func (s *Service) Delete(id int) error {
	if id <= 0 {
		return ErrInvalidEvaluationID
	}

	return s.repository.Delete(id)
}
