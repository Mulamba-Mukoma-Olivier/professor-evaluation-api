package courses

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type fakeCourseService struct {
	getAllFunc  func() ([]Course, error)
	getByIDFunc func(id int) (*Course, error)
	createFunc  func(request CreateCourseRequest) (*Course, error)
	updateFunc  func(id int, request UpdateCourseRequest) (*Course, error)
	deleteFunc  func(id int) error
}

func (f *fakeCourseService) GetAll() ([]Course, error) {
	if f.getAllFunc != nil {
		return f.getAllFunc()
	}

	return nil, nil
}

func (f *fakeCourseService) GetByID(id int) (*Course, error) {
	if f.getByIDFunc != nil {
		return f.getByIDFunc(id)
	}

	return nil, nil
}

func (f *fakeCourseService) Create(
	request CreateCourseRequest,
) (*Course, error) {
	if f.createFunc != nil {
		return f.createFunc(request)
	}

	return nil, nil
}

func (f *fakeCourseService) Update(
	id int,
	request UpdateCourseRequest,
) (*Course, error) {
	if f.updateFunc != nil {
		return f.updateFunc(id, request)
	}

	return nil, nil
}

func (f *fakeCourseService) Delete(id int) error {
	if f.deleteFunc != nil {
		return f.deleteFunc(id)
	}

	return nil
}

func setupHandlerTest(service CourseService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(service)

	router := gin.New()

	router.GET("/courses", handler.GetAll)
	router.GET("/courses/:id", handler.GetByID)
	router.POST("/courses", handler.Create)
	router.PUT("/courses/:id", handler.Update)
	router.DELETE("/courses/:id", handler.Delete)

	return router
}

func performRequest(
	router *gin.Engine,
	method string,
	path string,
	body string,
) *httptest.ResponseRecorder {
	request := httptest.NewRequest(
		method,
		path,
		strings.NewReader(body),
	)

	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	return recorder
}

// --------------------------------------------------
// GET ALL
// --------------------------------------------------

func TestHandler_GetAll_Success(t *testing.T) {
	service := &fakeCourseService{
		getAllFunc: func() ([]Course, error) {
			return []Course{
				{
					ID:           1,
					Code:         "INF301",
					Name:         "Programmation",
					Description:  "Cours de programmation",
					Department:   "Informatique",
					AcademicYear: "2025-2026",
				},
				{
					ID:           2,
					Code:         "INF302",
					Name:         "Bases de données",
					Description:  "Cours de bases de données",
					Department:   "Informatique",
					AcademicYear: "2025-2026",
				},
			}, nil
		},
	}

	router := setupHandlerTest(service)

	response := performRequest(
		router,
		http.MethodGet,
		"/courses",
		"",
	)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "INF301")
	assert.Contains(t, response.Body.String(), "INF302")
}

func TestHandler_GetAll_InternalError(t *testing.T) {
	service := &fakeCourseService{
		getAllFunc: func() ([]Course, error) {
			return nil, errors.New("database error")
		},
	}

	router := setupHandlerTest(service)

	response := performRequest(
		router,
		http.MethodGet,
		"/courses",
		"",
	)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Contains(t, response.Body.String(), "internal server error")
}

// --------------------------------------------------
// GET BY ID
// --------------------------------------------------

func TestHandler_GetByID_Success(t *testing.T) {
	service := &fakeCourseService{
		getByIDFunc: func(id int) (*Course, error) {
			assert.Equal(t, 1, id)

			return &Course{
				ID:           1,
				Code:         "INF301",
				Name:         "Programmation",
				Description:  "Cours de programmation",
				Department:   "Informatique",
				AcademicYear: "2025-2026",
			}, nil
		},
	}

	router := setupHandlerTest(service)

	response := performRequest(
		router,
		http.MethodGet,
		"/courses/1",
		"",
	)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "INF301")
	assert.Contains(t, response.Body.String(), "Programmation")
}

func TestHandler_GetByID_InvalidID(t *testing.T) {
	service := &fakeCourseService{}

	router := setupHandlerTest(service)

	response := performRequest(
		router,
		http.MethodGet,
		"/courses/abc",
		"",
	)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "invalid course ID")
}

func TestHandler_GetByID_NotFound(t *testing.T) {
	service := &fakeCourseService{
		getByIDFunc: func(id int) (*Course, error) {
			return nil, ErrCourseNotFound
		},
	}

	router := setupHandlerTest(service)

	response := performRequest(
		router,
		http.MethodGet,
		"/courses/999",
		"",
	)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Contains(t, response.Body.String(), "course not found")
}

func TestHandler_GetByID_InternalError(t *testing.T) {
	service := &fakeCourseService{
		getByIDFunc: func(id int) (*Course, error) {
			return nil, errors.New("database error")
		},
	}

	router := setupHandlerTest(service)

	response := performRequest(
		router,
		http.MethodGet,
		"/courses/1",
		"",
	)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Contains(t, response.Body.String(), "internal server error")
}

// --------------------------------------------------
// CREATE
// --------------------------------------------------

func TestHandler_Create_Success(t *testing.T) {
	service := &fakeCourseService{
		createFunc: func(request CreateCourseRequest) (*Course, error) {
			assert.Equal(t, "INF301", request.Code)
			assert.Equal(t, "Programmation", request.Name)

			return &Course{
				ID:           1,
				Code:         request.Code,
				Name:         request.Name,
				Description:  request.Description,
				Department:   request.Department,
				AcademicYear: request.AcademicYear,
			}, nil
		},
	}

	router := setupHandlerTest(service)

	body := `{
		"code": "INF301",
		"name": "Programmation",
		"description": "Cours de programmation",
		"department": "Informatique",
		"academic_year": "2025-2026"
	}`

	response := performRequest(
		router,
		http.MethodPost,
		"/courses",
		body,
	)

	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Contains(t, response.Body.String(), "INF301")
}

func TestHandler_Create_InvalidRequest(t *testing.T) {
	service := &fakeCourseService{}

	router := setupHandlerTest(service)

	body := `{
		"code": "",
		"name": "",
		"department": "",
		"academic_year": ""
	}`

	response := performRequest(
		router,
		http.MethodPost,
		"/courses",
		body,
	)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "invalid request")
}

func TestHandler_Create_InvalidJSON(t *testing.T) {
	service := &fakeCourseService{}

	router := setupHandlerTest(service)

	body := `{
		"code": "INF301",
		"name":
	}`

	response := performRequest(
		router,
		http.MethodPost,
		"/courses",
		body,
	)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "invalid request")
}

// --------------------------------------------------
// UPDATE
// --------------------------------------------------

func TestHandler_Update_Success(t *testing.T) {
	service := &fakeCourseService{
		updateFunc: func(
			id int,
			request UpdateCourseRequest,
		) (*Course, error) {
			assert.Equal(t, 1, id)
			assert.Equal(t, "INF301", request.Code)
			assert.Equal(t, "Programmation avancée", request.Name)

			return &Course{
				ID:           id,
				Code:         request.Code,
				Name:         request.Name,
				Description:  request.Description,
				Department:   request.Department,
				AcademicYear: request.AcademicYear,
			}, nil
		},
	}

	router := setupHandlerTest(service)

	body := `{
		"code": "INF301",
		"name": "Programmation avancée",
		"description": "Nouvelle description",
		"department": "Informatique",
		"academic_year": "2025-2026"
	}`

	response := performRequest(
		router,
		http.MethodPut,
		"/courses/1",
		body,
	)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Programmation avancée")
}

func TestHandler_Update_InvalidID(t *testing.T) {
	service := &fakeCourseService{}

	router := setupHandlerTest(service)

	body := `{
		"code": "INF301",
		"name": "Programmation",
		"department": "Informatique",
		"academic_year": "2025-2026"
	}`

	response := performRequest(
		router,
		http.MethodPut,
		"/courses/abc",
		body,
	)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "invalid course ID")
}

func TestHandler_Update_InvalidRequest(t *testing.T) {
	service := &fakeCourseService{}

	router := setupHandlerTest(service)

	body := `{
		"code": "",
		"name": "",
		"department": "",
		"academic_year": ""
	}`

	response := performRequest(
		router,
		http.MethodPut,
		"/courses/1",
		body,
	)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "invalid request")
}

func TestHandler_Update_NotFound(t *testing.T) {
	service := &fakeCourseService{
		updateFunc: func(
			id int,
			request UpdateCourseRequest,
		) (*Course, error) {
			return nil, ErrCourseNotFound
		},
	}

	router := setupHandlerTest(service)

	body := `{
		"code": "INF301",
		"name": "Programmation",
		"department": "Informatique",
		"academic_year": "2025-2026"
	}`

	response := performRequest(
		router,
		http.MethodPut,
		"/courses/999",
		body,
	)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Contains(t, response.Body.String(), "course not found")
}

// --------------------------------------------------
// DELETE
// --------------------------------------------------

func TestHandler_Delete_Success(t *testing.T) {
	service := &fakeCourseService{
		deleteFunc: func(id int) error {
			assert.Equal(t, 1, id)
			return nil
		},
	}

	router := setupHandlerTest(service)

	response := performRequest(
		router,
		http.MethodDelete,
		"/courses/1",
		"",
	)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(
		t,
		response.Body.String(),
		"course deleted successfully",
	)
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	service := &fakeCourseService{}

	router := setupHandlerTest(service)

	response := performRequest(
		router,
		http.MethodDelete,
		"/courses/abc",
		"",
	)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "invalid course ID")
}

func TestHandler_Delete_NotFound(t *testing.T) {
	service := &fakeCourseService{
		deleteFunc: func(id int) error {
			return ErrCourseNotFound
		},
	}

	router := setupHandlerTest(service)

	response := performRequest(
		router,
		http.MethodDelete,
		"/courses/999",
		"",
	)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Contains(t, response.Body.String(), "course not found")
}

func TestHandler_Delete_InternalError(t *testing.T) {
	service := &fakeCourseService{
		deleteFunc: func(id int) error {
			return errors.New("database error")
		},
	}

	router := setupHandlerTest(service)

	response := performRequest(
		router,
		http.MethodDelete,
		"/courses/1",
		"",
	)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Contains(t, response.Body.String(), "internal server error")
}
