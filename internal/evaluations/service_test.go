package evaluations

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// FakeEvaluationRepository permet de tester le service
// sans utiliser PostgreSQL.
type FakeEvaluationRepository struct {
	CreatedEvaluation *Evaluation
	CreateError       error

	Evaluations       []Evaluation
	GetAllError       error

	Evaluation        *Evaluation
	GetByIDError      error

	ExistsResult      bool
	ExistsError       error

	DeleteError       error
}

func (f *FakeEvaluationRepository) Create(evaluation Evaluation) (*Evaluation, error) {
	if f.CreateError != nil {
		return nil, f.CreateError
	}

	evaluation.ID = 1

	if evaluation.SubmittedAt.IsZero() {
		evaluation.SubmittedAt = time.Now()
	}

	f.CreatedEvaluation = &evaluation

	return &evaluation, nil
}

func (f *FakeEvaluationRepository) GetAll() ([]Evaluation, error) {
	if f.GetAllError != nil {
		return nil, f.GetAllError
	}

	return f.Evaluations, nil
}

func (f *FakeEvaluationRepository) GetByID(id int) (*Evaluation, error) {
	if f.GetByIDError != nil {
		return nil, f.GetByIDError
	}

	return f.Evaluation, nil
}

func (f *FakeEvaluationRepository) Exists(
	studentID int,
	professorID int,
	courseID int,
	academicYear string,
	period string,
) (bool, error) {

	if f.ExistsError != nil {
		return false, f.ExistsError
	}

	return f.ExistsResult, nil
}

func (f *FakeEvaluationRepository) Delete(id int) error {
	return f.DeleteError
}

// ---------------------------------------------------------
// Create
// ---------------------------------------------------------

func TestService_Create_Success(t *testing.T) {
	repository := &FakeEvaluationRepository{}

	service := NewService(repository)

	request := CreateEvaluationRequest{
		ProfessorID:  10,
		CourseID:     20,
		AcademicYear: "2025-2026",
		Period:       "S1",
		Answers: []Answer{
			{
				CriterionID: 1,
				Score:       4,
			},
			{
				CriterionID: 2,
				Score:       5,
			},
		},
	}

	result, err := service.Create(5, request)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, result.ID)
	assert.Equal(t, 5, result.StudentID)
	assert.Equal(t, 10, result.ProfessorID)
	assert.Equal(t, 20, result.CourseID)
	assert.Equal(t, "2025-2026", result.AcademicYear)
	assert.Equal(t, "S1", result.Period)
	assert.Len(t, result.Answers, 2)

	assert.Equal(t, 1, result.Answers[0].CriterionID)
	assert.Equal(t, 4, result.Answers[0].Score)

	assert.Equal(t, 2, result.Answers[1].CriterionID)
	assert.Equal(t, 5, result.Answers[1].Score)

	require.NotNil(t, repository.CreatedEvaluation)
	assert.Equal(t, 5, repository.CreatedEvaluation.StudentID)
}

// ---------------------------------------------------------
// Invalid student ID
// ---------------------------------------------------------

func TestService_Create_InvalidStudentID(t *testing.T) {
	repository := &FakeEvaluationRepository{}
	service := NewService(repository)

	request := validEvaluationRequest()

	result, err := service.Create(0, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidStudentID)
}

// ---------------------------------------------------------
// Invalid professor ID
// ---------------------------------------------------------

func TestService_Create_InvalidProfessorID(t *testing.T) {
	repository := &FakeEvaluationRepository{}
	service := NewService(repository)

	request := validEvaluationRequest()

	result, err := service.Create(5, request)

	request.ProfessorID = 0

	result, err = service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidProfessorID)
}

// ---------------------------------------------------------
// Invalid course ID
// ---------------------------------------------------------

func TestService_Create_InvalidCourseID(t *testing.T) {
	repository := &FakeEvaluationRepository{}
	service := NewService(repository)

	request := validEvaluationRequest()
	request.CourseID = 0

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidCourseID)
}

// ---------------------------------------------------------
// Academic year required
// ---------------------------------------------------------

func TestService_Create_AcademicYearRequired(t *testing.T) {
	repository := &FakeEvaluationRepository{}
	service := NewService(repository)

	request := validEvaluationRequest()
	request.AcademicYear = ""

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrAcademicYearRequired)
}

// ---------------------------------------------------------
// Period required
// ---------------------------------------------------------

func TestService_Create_PeriodRequired(t *testing.T) {
	repository := &FakeEvaluationRepository{}
	service := NewService(repository)

	request := validEvaluationRequest()
	request.Period = ""

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrPeriodRequired)
}

// ---------------------------------------------------------
// No answers
// ---------------------------------------------------------

func TestService_Create_NoAnswers(t *testing.T) {
	repository := &FakeEvaluationRepository{}
	service := NewService(repository)

	request := validEvaluationRequest()
	request.Answers = nil

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrAnswersRequired)
}

// ---------------------------------------------------------
// Invalid criterion ID
// ---------------------------------------------------------

func TestService_Create_InvalidCriterionID(t *testing.T) {
	repository := &FakeEvaluationRepository{}
	service := NewService(repository)

	request := validEvaluationRequest()
	request.Answers[0].CriterionID = 0

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidCriterionID)
}

// ---------------------------------------------------------
// Score too low
// ---------------------------------------------------------

func TestService_Create_ScoreTooLow(t *testing.T) {
	repository := &FakeEvaluationRepository{}
	service := NewService(repository)

	request := validEvaluationRequest()
	request.Answers[0].Score = 0

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidScore)
}

// ---------------------------------------------------------
// Score too high
// ---------------------------------------------------------

func TestService_Create_ScoreTooHigh(t *testing.T) {
	repository := &FakeEvaluationRepository{}
	service := NewService(repository)

	request := validEvaluationRequest()
	request.Answers[0].Score = 6

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidScore)
}

// ---------------------------------------------------------
// Duplicate evaluation
// ---------------------------------------------------------

func TestService_Create_DuplicateEvaluation(t *testing.T) {
	repository := &FakeEvaluationRepository{
		ExistsResult: true,
	}

	service := NewService(repository)

	request := validEvaluationRequest()

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrDuplicateEvaluation)
}

// ---------------------------------------------------------
// Repository Exists error
// ---------------------------------------------------------

func TestService_Create_ExistsError(t *testing.T) {
	repositoryError := errors.New("database error")

	repository := &FakeEvaluationRepository{
		ExistsError: repositoryError,
	}

	service := NewService(repository)

	request := validEvaluationRequest()

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repositoryError)
}

// ---------------------------------------------------------
// Repository Create error
// ---------------------------------------------------------

func TestService_Create_CreateError(t *testing.T) {
	repositoryError := errors.New("database error")

	repository := &FakeEvaluationRepository{
		CreateError: repositoryError,
	}

	service := NewService(repository)

	request := validEvaluationRequest()

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repositoryError)
}

// ---------------------------------------------------------
// Trim AcademicYear et Period
// ---------------------------------------------------------

func TestService_Create_TrimsFields(t *testing.T) {
	repository := &FakeEvaluationRepository{}
	service := NewService(repository)

	request := validEvaluationRequest()
	request.AcademicYear = " 2025-2026 "
	request.Period = " S1 "

	result, err := service.Create(5, request)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "2025-2026", result.AcademicYear)
	assert.Equal(t, "S1", result.Period)
}

// ---------------------------------------------------------
// GetAll
// ---------------------------------------------------------

func TestService_GetAll_Success(t *testing.T) {
	expected := []Evaluation{
		{
			ID:          1,
			StudentID:   5,
			ProfessorID: 10,
			CourseID:    20,
		},
		{
			ID:          2,
			StudentID:   6,
			ProfessorID: 11,
			CourseID:    21,
		},
	}

	repository := &FakeEvaluationRepository{
		Evaluations: expected,
	}

	service := NewService(repository)

	result, err := service.GetAll()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

// ---------------------------------------------------------
// GetAll repository error
// ---------------------------------------------------------

func TestService_GetAll_Error(t *testing.T) {
	repositoryError := errors.New("database error")

	repository := &FakeEvaluationRepository{
		GetAllError: repositoryError,
	}

	service := NewService(repository)

	result, err := service.GetAll()

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repositoryError)
}

// ---------------------------------------------------------
// GetByID success
// ---------------------------------------------------------

func TestService_GetByID_Success(t *testing.T) {
	expected := &Evaluation{
		ID:         1,
		StudentID:  5,
		ProfessorID: 10,
		CourseID:   20,
	}

	repository := &FakeEvaluationRepository{
		Evaluation: expected,
	}

	service := NewService(repository)

	result, err := service.GetByID(1)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, expected, result)
}

// ---------------------------------------------------------
// GetByID invalid ID
// ---------------------------------------------------------

func TestService_GetByID_InvalidID(t *testing.T) {
	repository := &FakeEvaluationRepository{}
	service := NewService(repository)

	result, err := service.GetByID(0)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidEvaluationID)
}

// ---------------------------------------------------------
// GetByID repository error
// ---------------------------------------------------------

func TestService_GetByID_Error(t *testing.T) {
	repositoryError := errors.New("database error")

	repository := &FakeEvaluationRepository{
		GetByIDError: repositoryError,
	}

	service := NewService(repository)

	result, err := service.GetByID(1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repositoryError)
}

// ---------------------------------------------------------
// Helper
// ---------------------------------------------------------

func validEvaluationRequest() CreateEvaluationRequest {
	return CreateEvaluationRequest{
		ProfessorID:  10,
		CourseID:     20,
		AcademicYear: "2025-2026",
		Period:       "S1",
		Answers: []Answer{
			{
				CriterionID: 1,
				Score:       4,
			},
		},
	}
}
