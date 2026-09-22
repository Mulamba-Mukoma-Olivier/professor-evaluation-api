package evaluations

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// FakeEvaluationService permet de tester le handler
// sans utiliser PostgreSQL ni le repository.
type FakeEvaluationService struct {
	CreateResult    *Evaluation
	CreateError     error

	Evaluations     []Evaluation
	GetAllError     error

	Evaluation      *Evaluation
	GetByIDError    error
}

func (f *FakeEvaluationService) Create(
	studentID int,
	request CreateEvaluationRequest,
) (*Evaluation, error) {
	if f.CreateError != nil {
		return nil, f.CreateError
	}

	if f.CreateResult != nil {
		return f.CreateResult, nil
	}

	return &Evaluation{
		ID:           1,
		StudentID:    studentID,
		ProfessorID:  request.ProfessorID,
		CourseID:     request.CourseID,
		AcademicYear: request.AcademicYear,
		Period:       request.Period,
		Answers: []EvaluationAnswer{
			{
				ID:          1,
				CriterionID: request.Answers[0].CriterionID,
				Score:       request.Answers[0].Score,
			},
		},
	}, nil
}

func (f *FakeEvaluationService) GetAll() ([]Evaluation, error) {
	if f.GetAllError != nil {
		return nil, f.GetAllError
	}

	return f.Evaluations, nil
}

func (f *FakeEvaluationService) GetByID(id int) (*Evaluation, error) {
	if f.GetByIDError != nil {
		return nil, f.GetByIDError
	}

	return f.Evaluation, nil
}

// ---------------------------------------------------------
// Create - succès
// ---------------------------------------------------------

func TestHandler_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{}

	handler := NewHandler(service)

	router := gin.New()
	router.POST("/students/:student_id/evaluations", handler.Create)

	body := `{
		"professor_id": 10,
		"course_id": 20,
		"academic_year": "2025-2026",
		"period": "S1",
		"answers": [
			{
				"criterion_id": 1,
				"score": 4
			}
		]
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/students/5/evaluations",
		bytes.NewBufferString(body),
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Code)

	assert.Contains(t, response.Body.String(), `"student_id":5`)
	assert.Contains(t, response.Body.String(), `"professor_id":10`)
	assert.Contains(t, response.Body.String(), `"course_id":20`)
}

// ---------------------------------------------------------
// Create - student ID invalide
// ---------------------------------------------------------

func TestHandler_Create_InvalidStudentID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{}
	handler := NewHandler(service)

	router := gin.New()
	router.POST("/students/:student_id/evaluations", handler.Create)

	body := `{
		"professor_id": 10,
		"course_id": 20,
		"academic_year": "2025-2026",
		"period": "S1",
		"answers": [
			{
				"criterion_id": 1,
				"score": 4
			}
		]
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/students/abc/evaluations",
		bytes.NewBufferString(body),
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "invalid student ID")
}

// ---------------------------------------------------------
// Create - JSON invalide
// ---------------------------------------------------------

func TestHandler_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{}
	handler := NewHandler(service)

	router := gin.New()
	router.POST("/students/:student_id/evaluations", handler.Create)

	body := `{
		"professor_id": 10,
		"course_id": 20,
		"answers":
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/students/5/evaluations",
		bytes.NewBufferString(body),
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "invalid request")
}

// ---------------------------------------------------------
// Create - erreur métier
// ---------------------------------------------------------

func TestHandler_Create_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{
		CreateError: errors.New("student has already evaluated this professor for this course"),
	}

	handler := NewHandler(service)

	router := gin.New()
	router.POST("/students/:student_id/evaluations", handler.Create)

	body := `{
		"professor_id": 10,
		"course_id": 20,
		"academic_year": "2025-2026",
		"period": "S1",
		"answers": [
			{
				"criterion_id": 1,
				"score": 4
			}
		]
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/students/5/evaluations",
		bytes.NewBufferString(body),
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(
		t,
		response.Body.String(),
		"student has already evaluated this professor for this course",
	)
}

// ---------------------------------------------------------
// Create - ressource inexistante
// ---------------------------------------------------------

func TestHandler_Create_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{
		CreateError: ErrEvaluationNotFound,
	}

	handler := NewHandler(service)

	router := gin.New()
	router.POST("/students/:student_id/evaluations", handler.Create)

	body := `{
		"professor_id": 10,
		"course_id": 20,
		"academic_year": "2025-2026",
		"period": "S1",
		"answers": [
			{
				"criterion_id": 1,
				"score": 4
			}
		]
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/students/5/evaluations",
		bytes.NewBufferString(body),
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Contains(t, response.Body.String(), "required resource not found")
}

// ---------------------------------------------------------
// GetAll - succès
// ---------------------------------------------------------

func TestHandler_GetAll_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

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

	service := &FakeEvaluationService{
		Evaluations: expected,
	}

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/evaluations", handler.GetAll)

	request := httptest.NewRequest(
		http.MethodGet,
		"/evaluations",
		nil,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	assert.Contains(t, response.Body.String(), `"evaluations"`)
	assert.Contains(t, response.Body.String(), `"student_id":5`)
	assert.Contains(t, response.Body.String(), `"student_id":6`)
}

// ---------------------------------------------------------
// GetAll - erreur serveur
// ---------------------------------------------------------

func TestHandler_GetAll_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{
		GetAllError: errors.New("database error"),
	}

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/evaluations", handler.GetAll)

	request := httptest.NewRequest(
		http.MethodGet,
		"/evaluations",
		nil,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Contains(t, response.Body.String(), "internal server error")
}

// ---------------------------------------------------------
// GetByID - succès
// ---------------------------------------------------------

func TestHandler_GetByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expected := &Evaluation{
		ID:          1,
		StudentID:   5,
		ProfessorID: 10,
		CourseID:    20,
	}

	service := &FakeEvaluationService{
		Evaluation: expected,
	}

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/evaluations/:id", handler.GetByID)

	request := httptest.NewRequest(
		http.MethodGet,
		"/evaluations/1",
		nil,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	assert.Contains(t, response.Body.String(), `"id":1`)
	assert.Contains(t, response.Body.String(), `"student_id":5`)
}

// ---------------------------------------------------------
// GetByID - ID invalide
// ---------------------------------------------------------

func TestHandler_GetByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{}
	handler := NewHandler(service)

	router := gin.New()
	router.GET("/evaluations/:id", handler.GetByID)

	request := httptest.NewRequest(
		http.MethodGet,
		"/evaluations/abc",
		nil,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "invalid evaluation ID")
}

// ---------------------------------------------------------
// GetByID - ID égal à zéro
// ---------------------------------------------------------

func TestHandler_GetByID_ZeroID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{}
	handler := NewHandler(service)

	router := gin.New()
	router.GET("/evaluations/:id", handler.GetByID)

	request := httptest.NewRequest(
		http.MethodGet,
		"/evaluations/0",
		nil,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "invalid evaluation ID")
}

// ---------------------------------------------------------
// GetByID - évaluation inexistante
// ---------------------------------------------------------

func TestHandler_GetByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{
		GetByIDError: ErrEvaluationNotFound,
	}

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/evaluations/:id", handler.GetByID)

	request := httptest.NewRequest(
		http.MethodGet,
		"/evaluations/999",
		nil,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Contains(t, response.Body.String(), "evaluation not found")
}

// ---------------------------------------------------------
// GetByID - erreur serveur
// ---------------------------------------------------------

func TestHandler_GetByID_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{
		GetByIDError: errors.New("database error"),
	}

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/evaluations/:id", handler.GetByID)

	request := httptest.NewRequest(
		http.MethodGet,
		"/evaluations/1",
		nil,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Contains(t, response.Body.String(), "internal server error")
}

// Vérifie que le fake implémente bien EvaluationService.
var _ EvaluationService = (*FakeEvaluationService)(nil)

// Évite une erreur si require est supprimé plus tard
var _ = require.NoError
