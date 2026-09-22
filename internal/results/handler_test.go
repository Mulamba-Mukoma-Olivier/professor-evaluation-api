package results

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FakeResultsService struct {
	result *ProfessorResult
	err    error
}

func (f *FakeResultsService) GetProfessorResult(
	professorID int,
	courseID int,
	academicYear string,
	period string,
) (*ProfessorResult, error) {
	return f.result, f.err
}

func setupResultsRouter(service ResultsService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	handler := NewHandler(service)

	router.GET(
		"/results/professors/:professor_id",
		handler.GetProfessorResult,
	)

	return router
}

func TestGetProfessorResult_Success(t *testing.T) {
	expectedResult := &ProfessorResult{
		ProfessorID:   1,
		CourseID:      10,
		AcademicYear:  "2025-2026",
		Period:        "semester_1",
		TotalReviews:  2,
		GlobalAverage: 4.0,
		Criteria: []CriterionResult{
			{
				CriterionID: 1,
				Average:     3.5,
				Responses:   2,
			},
			{
				CriterionID: 2,
				Average:     4.5,
				Responses:   2,
			},
		},
	}

	service := &FakeResultsService{
		result: expectedResult,
	}

	router := setupResultsRouter(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/results/professors/1?course_id=10&academic_year=2025-2026&period=semester_1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"professor_id":1`)
	assert.Contains(t, recorder.Body.String(), `"course_id":10`)
	assert.Contains(t, recorder.Body.String(), `"global_average":4`)
}

func TestGetProfessorResult_InvalidProfessorID(t *testing.T) {
	service := &FakeResultsService{}
	router := setupResultsRouter(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/results/professors/abc?course_id=10&academic_year=2025-2026&period=semester_1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid professor ID")
}

func TestGetProfessorResult_ZeroProfessorID(t *testing.T) {
	service := &FakeResultsService{}
	router := setupResultsRouter(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/results/professors/0?course_id=10&academic_year=2025-2026&period=semester_1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid professor ID")
}

func TestGetProfessorResult_InvalidCourseID(t *testing.T) {
	service := &FakeResultsService{}
	router := setupResultsRouter(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/results/professors/1?course_id=abc&academic_year=2025-2026&period=semester_1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid course ID")
}

func TestGetProfessorResult_MissingCourseID(t *testing.T) {
	service := &FakeResultsService{}
	router := setupResultsRouter(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/results/professors/1?academic_year=2025-2026&period=semester_1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid course ID")
}

func TestGetProfessorResult_MissingAcademicYear(t *testing.T) {
	service := &FakeResultsService{}
	router := setupResultsRouter(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/results/professors/1?course_id=10&period=semester_1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "academic_year is required")
}

func TestGetProfessorResult_WhitespaceAcademicYear(t *testing.T) {
	service := &FakeResultsService{}
	router := setupResultsRouter(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/results/professors/1?course_id=10&academic_year=%20%20&period=semester_1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "academic_year is required")
}

func TestGetProfessorResult_MissingPeriod(t *testing.T) {
	service := &FakeResultsService{}
	router := setupResultsRouter(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/results/professors/1?course_id=10&academic_year=2025-2026",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "period is required")
}

func TestGetProfessorResult_WhitespacePeriod(t *testing.T) {
	service := &FakeResultsService{}
	router := setupResultsRouter(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/results/professors/1?course_id=10&academic_year=2025-2026&period=%20%20",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "period is required")
}

func TestGetProfessorResult_NoEvaluations(t *testing.T) {
	service := &FakeResultsService{
		err: ErrNoEvaluations,
	}

	router := setupResultsRouter(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/results/professors/1?course_id=10&academic_year=2025-2026&period=semester_1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "no evaluations found")
}

func TestGetProfessorResult_NoAnswers(t *testing.T) {
	service := &FakeResultsService{
		err: ErrNoAnswers,
	}

	router := setupResultsRouter(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/results/professors/1?course_id=10&academic_year=2025-2026&period=semester_1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "evaluations contain no answers")
}

func TestGetProfessorResult_InternalServerError(t *testing.T) {
	service := &FakeResultsService{
		err: errors.New("database connection failed"),
	}

	router := setupResultsRouter(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/results/professors/1?course_id=10&academic_year=2025-2026&period=semester_1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "internal server error")
}
