package eligibility

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type fakeEligibilityService struct {
	createFunc func(eligibility *Eligibility) (*Eligibility, error)
	checkFunc  func(studentID int) (*Eligibility, error)
	updateFunc func(eligibility *Eligibility) (*Eligibility, error)
	deleteFunc func(studentID int) error
}

func (f *fakeEligibilityService) Create(eligibility *Eligibility) (*Eligibility, error) {
	if f.createFunc != nil {
		return f.createFunc(eligibility)
	}

	return nil, nil
}

func (f *fakeEligibilityService) Check(studentID int) (*Eligibility, error) {
	if f.checkFunc != nil {
		return f.checkFunc(studentID)
	}

	return nil, nil
}

func (f *fakeEligibilityService) Update(eligibility *Eligibility) (*Eligibility, error) {
	if f.updateFunc != nil {
		return f.updateFunc(eligibility)
	}

	return nil, nil
}

func (f *fakeEligibilityService) Delete(studentID int) error {
	if f.deleteFunc != nil {
		return f.deleteFunc(studentID)
	}

	return nil
}

func setupEligibilityHandlerTest(service EligibilityService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(service)

	router := gin.New()

	router.POST("/eligibility", handler.Create)
	router.GET("/eligibility/:student_id", handler.Check)
	router.PUT("/eligibility/:student_id", handler.Update)
	router.DELETE("/eligibility/:student_id", handler.Delete)

	return router
}

// --------------------------------------------------
// CHECK
// --------------------------------------------------

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
	assert.Contains(t, recorder.Body.String(), "ENROLLMENT_INVALID")
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
	assert.Contains(t, recorder.Body.String(), "invalid student ID")
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

// --------------------------------------------------
// CREATE
// --------------------------------------------------

func TestHandler_Create_Success(t *testing.T) {
	service := &fakeEligibilityService{
		createFunc: func(eligibility *Eligibility) (*Eligibility, error) {
			assert.Equal(t, 10, eligibility.StudentID)
			assert.True(t, eligibility.Enrollment)
			assert.True(t, eligibility.AcademicFees)

			eligibility.ID = 1
			eligibility.Eligible = true

			return eligibility, nil
		},
	}

	router := setupEligibilityHandlerTest(service)

	body := `{
		"student_id": 10,
		"enrollment": true,
		"academic_fees": true,
		"laboratory_fees": true,
		"access_fees": true
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/eligibility",
		strings.NewReader(body),
	)

	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"id":1`)
}

func TestHandler_Create_InvalidRequest(t *testing.T) {
	service := &fakeEligibilityService{}

	router := setupEligibilityHandlerTest(service)

	body := `{
		"student_id": "abc"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/eligibility",
		strings.NewReader(body),
	)

	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHandler_Create_ServiceError(t *testing.T) {
	service := &fakeEligibilityService{
		createFunc: func(eligibility *Eligibility) (*Eligibility, error) {
			return nil, errors.New("database error")
		},
	}

	router := setupEligibilityHandlerTest(service)

	body := `{
		"student_id": 10,
		"enrollment": true
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/eligibility",
		strings.NewReader(body),
	)

	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

// --------------------------------------------------
// UPDATE
// --------------------------------------------------

func TestHandler_Update_Success(t *testing.T) {
	service := &fakeEligibilityService{
		updateFunc: func(eligibility *Eligibility) (*Eligibility, error) {
			assert.Equal(t, 10, eligibility.StudentID)
			assert.False(t, eligibility.AcademicFees)

			eligibility.Eligible = false
			eligibility.Reasons = []string{
				"ACADEMIC_FEES_NOT_PAID",
			}

			return eligibility, nil
		},
	}

	router := setupEligibilityHandlerTest(service)

	body := `{
		"enrollment": true,
		"academic_fees": false,
		"laboratory_fees": true,
		"access_fees": true
	}`

	request := httptest.NewRequest(
		http.MethodPut,
		"/eligibility/10",
		strings.NewReader(body),
	)

	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"eligible":false`)
	assert.Contains(t, recorder.Body.String(), "ACADEMIC_FEES_NOT_PAID")
}

func TestHandler_Update_InvalidStudentID(t *testing.T) {
	service := &fakeEligibilityService{}

	router := setupEligibilityHandlerTest(service)

	body := `{
		"academic_fees": true
	}`

	request := httptest.NewRequest(
		http.MethodPut,
		"/eligibility/abc",
		strings.NewReader(body),
	)

	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHandler_Update_NotFound(t *testing.T) {
	service := &fakeEligibilityService{
		updateFunc: func(eligibility *Eligibility) (*Eligibility, error) {
			return nil, ErrEligibilityNotFound
		},
	}

	router := setupEligibilityHandlerTest(service)

	body := `{
		"academic_fees": true
	}`

	request := httptest.NewRequest(
		http.MethodPut,
		"/eligibility/999",
		strings.NewReader(body),
	)

	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

// --------------------------------------------------
// DELETE
// --------------------------------------------------

func TestHandler_Delete_Success(t *testing.T) {
	service := &fakeEligibilityService{
		deleteFunc: func(studentID int) error {
			assert.Equal(t, 10, studentID)
			return nil
		},
	}

	router := setupEligibilityHandlerTest(service)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/eligibility/10",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(
		t,
		recorder.Body.String(),
		"eligibility deleted successfully",
	)
}

func TestHandler_Delete_InvalidStudentID(t *testing.T) {
	service := &fakeEligibilityService{}

	router := setupEligibilityHandlerTest(service)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/eligibility/abc",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHandler_Delete_NotFound(t *testing.T) {
	service := &fakeEligibilityService{
		deleteFunc: func(studentID int) error {
			return ErrEligibilityNotFound
		},
	}

	router := setupEligibilityHandlerTest(service)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/eligibility/999",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}