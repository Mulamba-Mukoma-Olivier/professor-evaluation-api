package evaluations

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// FakeEvaluationRepository permet de tester le service
// sans utiliser la base de données.
type FakeEvaluationRepository struct {
	// Create
	CreatedEvaluation *Evaluation
	CreateError       error

	// GetAll
	Evaluations []Evaluation
	GetAllError error

	// GetByID
	Evaluation   *Evaluation
	GetByIDError error

	// Exists
	ExistsResult bool
	ExistsError  error

	// Professor
	ProfessorExistsResult bool
	ProfessorExistsError  error

	// Course
	CourseExistsResult bool
	CourseExistsError  error

	// Criterion
	CriterionExistsResult bool
	CriterionExistsError  error

	// Criterion active
	CriterionActiveResult bool
	CriterionActiveError  error

	// Delete
	DeleteError error
}

// =========================================================
// FAKE REPOSITORY - CREATE
// =========================================================

func (f *FakeEvaluationRepository) Create(
	evaluation Evaluation,
) (*Evaluation, error) {

	if f.CreateError != nil {
		return nil, f.CreateError
	}

	evaluation.ID = 1

	f.CreatedEvaluation = &evaluation

	return &evaluation, nil
}

// =========================================================
// FAKE REPOSITORY - GET ALL
// =========================================================

func (f *FakeEvaluationRepository) GetAll() ([]Evaluation, error) {

	if f.GetAllError != nil {
		return nil, f.GetAllError
	}

	return f.Evaluations, nil
}

// =========================================================
// FAKE REPOSITORY - GET BY ID
// =========================================================

func (f *FakeEvaluationRepository) GetByID(
	id int,
) (*Evaluation, error) {

	if f.GetByIDError != nil {
		return nil, f.GetByIDError
	}

	return f.Evaluation, nil
}

// =========================================================
// FAKE REPOSITORY - EXISTS
// =========================================================

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

// =========================================================
// FAKE REPOSITORY - PROFESSOR EXISTS
// =========================================================

func (f *FakeEvaluationRepository) ProfessorExists(
	professorID int,
) (bool, error) {

	if f.ProfessorExistsError != nil {
		return false, f.ProfessorExistsError
	}

	return f.ProfessorExistsResult, nil
}

// =========================================================
// FAKE REPOSITORY - COURSE EXISTS
// =========================================================

func (f *FakeEvaluationRepository) CourseExists(
	courseID int,
) (bool, error) {

	if f.CourseExistsError != nil {
		return false, f.CourseExistsError
	}

	return f.CourseExistsResult, nil
}

// =========================================================
// FAKE REPOSITORY - CRITERION EXISTS
// =========================================================

func (f *FakeEvaluationRepository) CriterionExists(
	criterionID int,
) (bool, error) {

	if f.CriterionExistsError != nil {
		return false, f.CriterionExistsError
	}

	return f.CriterionExistsResult, nil
}

// =========================================================
// FAKE REPOSITORY - CRITERION ACTIVE
// =========================================================

func (f *FakeEvaluationRepository) CriterionIsActive(
	criterionID int,
) (bool, error) {

	if f.CriterionActiveError != nil {
		return false, f.CriterionActiveError
	}

	return f.CriterionActiveResult, nil
}

// =========================================================
// FAKE REPOSITORY - DELETE
// =========================================================

func (f *FakeEvaluationRepository) Delete(id int) error {
	return f.DeleteError
}

// =========================================================
// CREATE
// =========================================================

// ---------------------------------------------------------
// Create - Success
// ---------------------------------------------------------

func TestService_Create_Success(t *testing.T) {

	repository := &FakeEvaluationRepository{
		ProfessorExistsResult: true,
		CourseExistsResult:    true,
		CriterionExistsResult: true,
		CriterionActiveResult: true,
	}

	service := NewService(repository)

	request := validEvaluationRequest()

	result, err := service.Create(5, request)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, result.ID)
	assert.Equal(t, 5, result.StudentID)
	assert.Equal(t, 10, result.ProfessorID)
	assert.Equal(t, 20, result.CourseID)

	assert.Equal(t, "2025-2026", result.AcademicYear)
	assert.Equal(t, "S1", result.Period)

	assert.Len(t, result.Answers, 1)

	assert.Equal(t, 1, result.Answers[0].CriterionID)
	assert.Equal(t, 4, result.Answers[0].Score)

	require.NotNil(t, repository.CreatedEvaluation)

	assert.Equal(
		t,
		5,
		repository.CreatedEvaluation.StudentID,
	)
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
	request.ProfessorID = 0

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidProfessorID)
}

// ---------------------------------------------------------
// Professor not found
// ---------------------------------------------------------

func TestService_Create_ProfessorNotFound(t *testing.T) {

	repository := &FakeEvaluationRepository{
		ProfessorExistsResult: false,
	}

	service := NewService(repository)

	request := validEvaluationRequest()

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrProfessorNotFound)
}

// ---------------------------------------------------------
// Professor repository error
// ---------------------------------------------------------

func TestService_Create_ProfessorExistsError(t *testing.T) {

	repositoryError := errors.New("database error")

	repository := &FakeEvaluationRepository{
		ProfessorExistsError: repositoryError,
	}

	service := NewService(repository)

	request := validEvaluationRequest()

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repositoryError)
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
// Course not found
// ---------------------------------------------------------

func TestService_Create_CourseNotFound(t *testing.T) {

	repository := &FakeEvaluationRepository{
		ProfessorExistsResult: true,
		CourseExistsResult:    false,
	}

	service := NewService(repository)

	request := validEvaluationRequest()

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrCourseNotFound)
}

// ---------------------------------------------------------
// Course repository error
// ---------------------------------------------------------

func TestService_Create_CourseExistsError(t *testing.T) {

	repositoryError := errors.New("database error")

	repository := &FakeEvaluationRepository{
		ProfessorExistsResult: true,
		CourseExistsError:    repositoryError,
	}

	service := NewService(repository)

	request := validEvaluationRequest()

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repositoryError)
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

	repository := &FakeEvaluationRepository{
		ProfessorExistsResult: true,
		CourseExistsResult:    true,
	}

	service := NewService(repository)

	request := validEvaluationRequest()
	request.Answers[0].CriterionID = 0

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidCriterionID)
}

// ---------------------------------------------------------
// Criterion not found
// ---------------------------------------------------------

func TestService_Create_CriterionNotFound(t *testing.T) {

	repository := &FakeEvaluationRepository{
		ProfessorExistsResult: true,
		CourseExistsResult:    true,
		CriterionExistsResult: false,
	}

	service := NewService(repository)

	request := validEvaluationRequest()

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrCriterionNotFound)
}

// ---------------------------------------------------------
// Criterion exists repository error
// ---------------------------------------------------------

func TestService_Create_CriterionExistsError(t *testing.T) {

	repositoryError := errors.New("database error")

	repository := &FakeEvaluationRepository{
		ProfessorExistsResult: true,
		CourseExistsResult:    true,
		CriterionExistsError: repositoryError,
	}

	service := NewService(repository)

	request := validEvaluationRequest()

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repositoryError)
}

// ---------------------------------------------------------
// Criterion inactive
// ---------------------------------------------------------

func TestService_Create_CriterionInactive(t *testing.T) {

	repository := &FakeEvaluationRepository{
		ProfessorExistsResult: true,
		CourseExistsResult:    true,
		CriterionExistsResult: true,
		CriterionActiveResult: false,
	}

	service := NewService(repository)

	request := validEvaluationRequest()

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrCriterionInactive)
}

// ---------------------------------------------------------
// Criterion active repository error
// ---------------------------------------------------------

func TestService_Create_CriterionActiveError(t *testing.T) {

	repositoryError := errors.New("database error")

	repository := &FakeEvaluationRepository{
		ProfessorExistsResult: true,
		CourseExistsResult:    true,
		CriterionExistsResult: true,
		CriterionActiveError:  repositoryError,
	}

	service := NewService(repository)

	request := validEvaluationRequest()

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repositoryError)
}

// ---------------------------------------------------------
// Score too low
// ---------------------------------------------------------

func TestService_Create_ScoreTooLow(t *testing.T) {

	repository := &FakeEvaluationRepository{
		ProfessorExistsResult: true,
		CourseExistsResult:    true,
		CriterionExistsResult: true,
		CriterionActiveResult: true,
	}

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

	repository := &FakeEvaluationRepository{
		ProfessorExistsResult: true,
		CourseExistsResult:    true,
		CriterionExistsResult: true,
		CriterionActiveResult: true,
	}

	service := NewService(repository)

	request := validEvaluationRequest()
	request.Answers[0].Score = 6

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidScore)
}

// ---------------------------------------------------------
// Duplicate criterion
// ---------------------------------------------------------

func TestService_Create_DuplicateCriterion(t *testing.T) {

	repository := &FakeEvaluationRepository{
		ProfessorExistsResult: true,
		CourseExistsResult:    true,
		CriterionExistsResult: true,
		CriterionActiveResult: true,
	}

	service := NewService(repository)

	request := validEvaluationRequest()

	request.Answers = []Answer{
		{
			CriterionID: 1,
			Score:       4,
		},
		{
			CriterionID: 1,
			Score:       5,
		},
	}

	result, err := service.Create(5, request)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrDuplicateCriterion)
}

// ---------------------------------------------------------
// Duplicate evaluation
// ---------------------------------------------------------

func TestService_Create_DuplicateEvaluation(t *testing.T) {

	repository := &FakeEvaluationRepository{
		ProfessorExistsResult: true,
		CourseExistsResult:    true,
		CriterionExistsResult: true,
		CriterionActiveResult: true,
		ExistsResult:          true,
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
		ProfessorExistsResult: true,
		CourseExistsResult:    true,
		CriterionExistsResult: true,
		CriterionActiveResult: true,
		ExistsError:           repositoryError,
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
		ProfessorExistsResult: true,
		CourseExistsResult:    true,
		CriterionExistsResult: true,
		CriterionActiveResult: true,
		CreateError:           repositoryError,
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

	repository := &FakeEvaluationRepository{
		ProfessorExistsResult: true,
		CourseExistsResult:    true,
		CriterionExistsResult: true,
		CriterionActiveResult: true,
	}

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

// =========================================================
// GET ALL
// =========================================================

// ---------------------------------------------------------
// GetAll - Success
// ---------------------------------------------------------

func TestService_GetAll_Success(t *testing.T) {

	expected := []Evaluation{
		{
			ID:         1,
			StudentID:  5,
			ProfessorID: 10,
			CourseID:   20,
		},
		{
			ID:         2,
			StudentID:  6,
			ProfessorID: 11,
			CourseID:   21,
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
// GetAll - Repository error
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

// =========================================================
// GET BY ID
// =========================================================

// ---------------------------------------------------------
// GetByID - Success
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
// GetByID - Invalid ID
// ---------------------------------------------------------

func TestService_GetByID_InvalidID(t *testing.T) {

	repository := &FakeEvaluationRepository{}
	service := NewService(repository)

	result, err := service.GetByID(0)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidEvaluationID)
}

// ---------------------------------------------------------
// GetByID - Repository error
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

// =========================================================
// DELETE
// =========================================================

// ---------------------------------------------------------
// Delete - Success
// ---------------------------------------------------------

func TestService_Delete_Success(t *testing.T) {

	repository := &FakeEvaluationRepository{}
	service := NewService(repository)

	err := service.Delete(1)

	assert.NoError(t, err)
}

// ---------------------------------------------------------
// Delete - Invalid ID
// ---------------------------------------------------------

func TestService_Delete_InvalidID(t *testing.T) {

	repository := &FakeEvaluationRepository{}
	service := NewService(repository)

	err := service.Delete(0)

	assert.ErrorIs(t, err, ErrInvalidEvaluationID)
}

// ---------------------------------------------------------
// Delete - Repository error
// ---------------------------------------------------------

func TestService_Delete_Error(t *testing.T) {

	repositoryError := errors.New("database error")

	repository := &FakeEvaluationRepository{
		DeleteError: repositoryError,
	}

	service := NewService(repository)

	err := service.Delete(1)

	assert.ErrorIs(t, err, repositoryError)
}

// =========================================================
// HELPER
// =========================================================

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