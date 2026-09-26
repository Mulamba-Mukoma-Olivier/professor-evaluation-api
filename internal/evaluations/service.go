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

	ErrProfessorNotFound = errors.New("professor not found")
	ErrCourseNotFound    = errors.New("course not found")
	ErrCriterionNotFound = errors.New("criterion not found")

	ErrCriterionInactive = errors.New("criterion is inactive")

	ErrDuplicateCriterion = errors.New(
		"criterion cannot appear more than once in an evaluation",
	)

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

	// Vérification de l'existence des données liées.
	ProfessorExists(professorID int) (bool, error)

	CourseExists(courseID int) (bool, error)

	CriterionExists(criterionID int) (bool, error)

	CriterionIsActive(criterionID int) (bool, error)

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
// des données et des principales règles métier.
func (s *Service) Create(
	studentID int,
	request CreateEvaluationRequest,
) (*Evaluation, error) {

	// --------------------------------------------------
	// 1. Validation de l'étudiant
	// --------------------------------------------------
	if studentID <= 0 {
		return nil, ErrInvalidStudentID
	}

	// --------------------------------------------------
	// 2. Validation du professeur
	// --------------------------------------------------
	if request.ProfessorID <= 0 {
		return nil, ErrInvalidProfessorID
	}

	professorExists, err := s.repository.ProfessorExists(request.ProfessorID)
	if err != nil {
		return nil, err
	}

	if !professorExists {
		return nil, ErrProfessorNotFound
	}

	// --------------------------------------------------
	// 3. Validation du cours
	// --------------------------------------------------
	if request.CourseID <= 0 {
		return nil, ErrInvalidCourseID
	}

	courseExists, err := s.repository.CourseExists(request.CourseID)
	if err != nil {
		return nil, err
	}

	if !courseExists {
		return nil, ErrCourseNotFound
	}

	// --------------------------------------------------
	// 4. Validation de l'année académique
	// --------------------------------------------------
	request.AcademicYear = strings.TrimSpace(request.AcademicYear)

	if request.AcademicYear == "" {
		return nil, ErrAcademicYearRequired
	}

	// --------------------------------------------------
	// 5. Validation de la période
	// --------------------------------------------------
	request.Period = strings.TrimSpace(request.Period)

	if request.Period == "" {
		return nil, ErrPeriodRequired
	}

	// --------------------------------------------------
	// 6. Validation des réponses
	// --------------------------------------------------
	if len(request.Answers) == 0 {
		return nil, ErrAnswersRequired
	}

	// Permet d'empêcher le même critère plusieurs fois.
	criteria := make(map[int]struct{}, len(request.Answers))

	for _, answer := range request.Answers {

		// ----------------------------------------------
		// Validation du CriterionID
		// ----------------------------------------------
		if answer.CriterionID <= 0 {
			return nil, ErrInvalidCriterionID
		}

		// ----------------------------------------------
		// Vérification de l'existence du critère
		// ----------------------------------------------
		criterionExists, err := s.repository.CriterionExists(
			answer.CriterionID,
		)

		if err != nil {
			return nil, err
		}

		if !criterionExists {
			return nil, ErrCriterionNotFound
		}

		// ----------------------------------------------
		// Vérification que le critère est actif
		// ----------------------------------------------
		criterionActive, err := s.repository.CriterionIsActive(
			answer.CriterionID,
		)

		if err != nil {
			return nil, err
		}

		if !criterionActive {
			return nil, ErrCriterionInactive
		}

		// ----------------------------------------------
		// Validation de la note
		// ----------------------------------------------
		if answer.Score < 1 || answer.Score > 5 {
			return nil, ErrInvalidScore
		}

		// ----------------------------------------------
		// Détection des critères en double
		// ----------------------------------------------
		if _, exists := criteria[answer.CriterionID]; exists {
			return nil, ErrDuplicateCriterion
		}

		criteria[answer.CriterionID] = struct{}{}
	}

	// --------------------------------------------------
	// 7. Vérification d'une évaluation existante
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
	// 8. Transformation DTO -> Model
	// --------------------------------------------------
	answers := make([]EvaluationAnswer, 0, len(request.Answers))

	for _, answer := range request.Answers {
		answers = append(answers, EvaluationAnswer{
			CriterionID: answer.CriterionID,
			Score:       answer.Score,
		})
	}

	// --------------------------------------------------
	// 9. Création de l'évaluation
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