package evaluations

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type FakeEvaluationService struct {
	CreateResult *Evaluation
	CreateError  error

	Evaluations []Evaluation
	GetAllError error

	Evaluation   *Evaluation
	GetByIDError error

	DeleteError error
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

func (f *FakeEvaluationService) Delete(id int) error {
	return f.DeleteError
}

var _ EvaluationService = (*FakeEvaluationService)(nil)

func validCreateEvaluationBody() string {
	return `{
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
}

func TestHandler_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{}
	handler := NewHandler(service)

	router := gin.New()

	router.POST("/evaluations", func(c *gin.Context) {
		c.Set("user_id", 5)
		handler.Create(c)
	})

	request := httptest.NewRequest(
		http.MethodPost,
		"/evaluations",
		bytes.NewBufferString(validCreateEvaluationBody()),
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Contains(t, response.Body.String(), `"student_id":5`)
	assert.Contains(t, response.Body.String(), `"professor_id":10`)
	assert.Contains(t, response.Body.String(), `"course_id":20`)
}

func TestHandler_Create_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{}
	handler := NewHandler(service)

	router := gin.New()
	router.POST("/evaluations", handler.Create)

	request := httptest.NewRequest(
		http.MethodPost,
		"/evaluations",
		bytes.NewBufferString(validCreateEvaluationBody()),
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.Contains(t, response.Body.String(), "user not authenticated")
}

func TestHandler_Create_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{}
	handler := NewHandler(service)

	router := gin.New()

	router.POST("/evaluations", func(c *gin.Context) {
		c.Set("user_id", "5")
		handler.Create(c)
	})

	request := httptest.NewRequest(
		http.MethodPost,
		"/evaluations",
		bytes.NewBufferString(validCreateEvaluationBody()),
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.Contains(t, response.Body.String(), "invalid user ID")
}

func TestHandler_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{}
	handler := NewHandler(service)

	router := gin.New()

	router.POST("/evaluations", func(c *gin.Context) {
		c.Set("user_id", 5)
		handler.Create(c)
	})

	body := `{
		"professor_id": 10,
		"course_id": 20,
		"answers":
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/evaluations",
		bytes.NewBufferString(body),
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "invalid request")
}

func TestHandler_Create_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{
		CreateError: ErrDuplicateEvaluation,
	}

	handler := NewHandler(service)

	router := gin.New()

	router.POST("/evaluations", func(c *gin.Context) {
		c.Set("user_id", 5)
		handler.Create(c)
	})

	request := httptest.NewRequest(
		http.MethodPost,
		"/evaluations",
		bytes.NewBufferString(validCreateEvaluationBody()),
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(
		t,
		response.Body.String(),
		"student has already evaluated this professor",
	)
}

func TestHandler_Create_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{
		CreateError: errors.New("database error"),
	}

	handler := NewHandler(service)

	router := gin.New()

	router.POST("/evaluations", func(c *gin.Context) {
		c.Set("user_id", 5)
		handler.Create(c)
	})

	request := httptest.NewRequest(
		http.MethodPost,
		"/evaluations",
		bytes.NewBufferString(validCreateEvaluationBody()),
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Contains(t, response.Body.String(), "internal server error")
}

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

func TestHandler_Delete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{}

	handler := NewHandler(service)

	router := gin.New()
	router.DELETE("/evaluations/:id", handler.Delete)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/evaluations/1",
		nil,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "evaluation deleted successfully")
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{}

	handler := NewHandler(service)

	router := gin.New()
	router.DELETE("/evaluations/:id", handler.Delete)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/evaluations/abc",
		nil,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "invalid evaluation ID")
}

func TestHandler_Delete_ZeroID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{}

	handler := NewHandler(service)

	router := gin.New()
	router.DELETE("/evaluations/:id", handler.Delete)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/evaluations/0",
		nil,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "invalid evaluation ID")
}

func TestHandler_Delete_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{
		DeleteError: ErrEvaluationNotFound,
	}

	handler := NewHandler(service)

	router := gin.New()
	router.DELETE("/evaluations/:id", handler.Delete)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/evaluations/999",
		nil,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Contains(t, response.Body.String(), "evaluation not found")
}

func TestHandler_Delete_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &FakeEvaluationService{
		DeleteError: errors.New("database error"),
	}

	handler := NewHandler(service)

	router := gin.New()
	router.DELETE("/evaluations/:id", handler.Delete)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/evaluations/1",
		nil,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Contains(t, response.Body.String(), "internal server error")
}
