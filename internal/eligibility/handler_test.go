package eligibility

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type fakeEligibilityService struct {
	checkFunc func(studentID int) (*Eligibility, error)
}

func (f *fakeEligibilityService) Check(studentID int) (*Eligibility, error) {
	return f.checkFunc(studentID)
}

func setupEligibilityHandlerTest(service EligibilityService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(service)

	router := gin.New()

	router.GET(
		"/eligibility/:student_id",
		handler.Check,
	)

	return router
}

func TestHandler_Check_Success_Eligible(t *testing.T) {
	service := &fakeEligibilityService{
		checkFunc: func(studentID int) (*Eligibility, error) {
			assert.Equal(t, 10, studentID)

			return &Eligibility{
				ID:             1,
				StudentID:      10,
				Enrollment:     true,
				AcademicFees:   true,
				LaboratoryFees: true,
				AccessFees:     true,
				Eligible:       true,
				Reasons:        []string{},
			}, nil
		},
	}

	router := setupEligibilityHandlerTest(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/eligibility/10",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"eligible":true`)
}

func TestHandler_Check_Success_NotEligible(t *testing.T) {
	service := &fakeEligibilityService{
		checkFunc: func(studentID int) (*Eligibility, error) {
			return &Eligibility{
				StudentID:    studentID,
				Enrollment:   false,
				AcademicFees: true,
				Eligible:     false,
				Reasons: []string{
					"ENROLLMENT_INVALID",
				},
			}, nil
		},
	}

	router := setupEligibilityHandlerTest(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/eligibility/10",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(
		t,
		recorder.Body.String(),
		"ENROLLMENT_INVALID",
	)
	assert.Contains(t, recorder.Body.String(), `"eligible":false`)
}

func TestHandler_Check_InvalidStudentID(t *testing.T) {
	service := &fakeEligibilityService{}

	router := setupEligibilityHandlerTest(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/eligibility/abc",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(
		t,
		recorder.Body.String(),
		"invalid student ID",
	)
}

func TestHandler_Check_ZeroStudentID(t *testing.T) {
	service := &fakeEligibilityService{}

	router := setupEligibilityHandlerTest(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/eligibility/0",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHandler_Check_NotFound(t *testing.T) {
	service := &fakeEligibilityService{
		checkFunc: func(studentID int) (*Eligibility, error) {
			return nil, ErrEligibilityNotFound
		},
	}

	router := setupEligibilityHandlerTest(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/eligibility/999",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Contains(
		t,
		recorder.Body.String(),
		"student eligibility not found",
	)
}

func TestHandler_Check_InternalError(t *testing.T) {
	service := &fakeEligibilityService{
		checkFunc: func(studentID int) (*Eligibility, error) {
			return nil, errors.New("database error")
		},
	}

	router := setupEligibilityHandlerTest(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/eligibility/10",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(
		t,
		recorder.Body.String(),
		"internal server error",
	)
}
