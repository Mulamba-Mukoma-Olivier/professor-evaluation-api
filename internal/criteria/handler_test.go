package criteria

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeCriterionService struct {
	getAllFunc     func() ([]Criterion, error)
	getActiveFunc  func() ([]Criterion, error)
	getByIDFunc    func(int) (*Criterion, error)
	createFunc     func(CreateCriterionRequest) (*Criterion, error)
	updateFunc     func(int, UpdateCriterionRequest) (*Criterion, error)
	deleteFunc     func(int) error
}

func (f *fakeCriterionService) GetAll() ([]Criterion, error) {
	return f.getAllFunc()
}

func (f *fakeCriterionService) GetActive() ([]Criterion, error) {
	return f.getActiveFunc()
}

func (f *fakeCriterionService) GetByID(id int) (*Criterion, error) {
	return f.getByIDFunc(id)
}

func (f *fakeCriterionService) Create(request CreateCriterionRequest) (*Criterion, error) {
	return f.createFunc(request)
}

func (f *fakeCriterionService) Update(
	id int,
	request UpdateCriterionRequest,
) (*Criterion, error) {
	return f.updateFunc(id, request)
}

func (f *fakeCriterionService) Delete(id int) error {
	return f.deleteFunc(id)
}

func setupHandlerTest() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	service := &fakeCriterionService{
		getAllFunc: func() ([]Criterion, error) {
			return []Criterion{}, nil
		},
		getActiveFunc: func() ([]Criterion, error) {
			return []Criterion{}, nil
		},
		getByIDFunc: func(id int) (*Criterion, error) {
			return nil, ErrCriterionNotFound
		},
		createFunc: func(request CreateCriterionRequest) (*Criterion, error) {
			return nil, nil
		},
		updateFunc: func(
			id int,
			request UpdateCriterionRequest,
		) (*Criterion, error) {
			return nil, nil
		},
		deleteFunc: func(id int) error {
			return nil
		},
	}

	handler := NewHandler(service)

	router.GET("/criteria", handler.GetAll)
	router.GET("/criteria/active", handler.GetActive)
	router.GET("/criteria/:id", handler.GetByID)
	router.POST("/criteria", handler.Create)
	router.PUT("/criteria/:id", handler.Update)
	router.DELETE("/criteria/:id", handler.Delete)

	return router
}

func performRequest(
	router *gin.Engine,
	method string,
	url string,
	body interface{},
) *httptest.ResponseRecorder {

	var requestBody []byte

	if body != nil {
		requestBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(
		method,
		url,
		bytes.NewBuffer(requestBody),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	return recorder
}

// ---------------------------------------------------------
// GetAll
// ---------------------------------------------------------

func TestHandler_GetAll_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &fakeCriterionService{
		getAllFunc: func() ([]Criterion, error) {
			return []Criterion{
				{
					ID:          1,
					Name:        "Clarté",
					Description: "Clarté des explications",
					MaxScore:    5,
					Active:      true,
				},
			}, nil
		},
	}

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/criteria", handler.GetAll)

	recorder := performRequest(
		router,
		http.MethodGet,
		"/criteria",
		nil,
	)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)

	require.NoError(t, err)
	assert.NotNil(t, response["criteria"])
}

func TestHandler_GetAll_Error(t *testing.T) {
	service := &fakeCriterionService{
		getAllFunc: func() ([]Criterion, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/criteria", handler.GetAll)

	recorder := performRequest(
		router,
		http.MethodGet,
		"/criteria",
		nil,
	)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "failed to retrieve criteria")
}

// ---------------------------------------------------------
// GetActive
// ---------------------------------------------------------

func TestHandler_GetActive_Success(t *testing.T) {
	service := &fakeCriterionService{
		getActiveFunc: func() ([]Criterion, error) {
			return []Criterion{
				{
					ID:       1,
					Name:     "Clarté",
					MaxScore: 5,
					Active:   true,
				},
			}, nil
		},
	}

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/criteria/active", handler.GetActive)

	recorder := performRequest(
		router,
		http.MethodGet,
		"/criteria/active",
		nil,
	)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Clarté")
}

func TestHandler_GetActive_Error(t *testing.T) {
	service := &fakeCriterionService{
		getActiveFunc: func() ([]Criterion, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/criteria/active", handler.GetActive)

	recorder := performRequest(
		router,
		http.MethodGet,
		"/criteria/active",
		nil,
	)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(
		t,
		recorder.Body.String(),
		"failed to retrieve active criteria",
	)
}

// ---------------------------------------------------------
// GetByID
// ---------------------------------------------------------

func TestHandler_GetByID_Success(t *testing.T) {
	service := &fakeCriterionService{
		getByIDFunc: func(id int) (*Criterion, error) {
			return &Criterion{
				ID:          id,
				Name:        "Clarté",
				Description: "Clarté des explications",
				MaxScore:    5,
				Active:      true,
			}, nil
		},
	}

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/criteria/:id", handler.GetByID)

	recorder := performRequest(
		router,
		http.MethodGet,
		"/criteria/1",
		nil,
	)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Clarté")
}

func TestHandler_GetByID_InvalidID(t *testing.T) {
	service := &fakeCriterionService{}

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/criteria/:id", handler.GetByID)

	recorder := performRequest(
		router,
		http.MethodGet,
		"/criteria/abc",
		nil,
	)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid criterion ID")
}

func TestHandler_GetByID_ZeroID(t *testing.T) {
	service := &fakeCriterionService{}

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/criteria/:id", handler.GetByID)

	recorder := performRequest(
		router,
		http.MethodGet,
		"/criteria/0",
		nil,
	)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid criterion ID")
}

func TestHandler_GetByID_NotFound(t *testing.T) {
	service := &fakeCriterionService{
		getByIDFunc: func(id int) (*Criterion, error) {
			return nil, ErrCriterionNotFound
		},
	}

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/criteria/:id", handler.GetByID)

	recorder := performRequest(
		router,
		http.MethodGet,
		"/criteria/99",
		nil,
	)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "criterion not found")
}

func TestHandler_GetByID_InternalError(t *testing.T) {
	service := &fakeCriterionService{
		getByIDFunc: func(id int) (*Criterion, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewHandler(service)

	router := gin.New()
	router.GET("/criteria/:id", handler.GetByID)

	recorder := performRequest(
		router,
		http.MethodGet,
		"/criteria/1",
		nil,
	)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(
		t,
		recorder.Body.String(),
		"failed to retrieve criterion",
	)
}

// ---------------------------------------------------------
// Create
// ---------------------------------------------------------

func TestHandler_Create_Success(t *testing.T) {
	service := &fakeCriterionService{
		createFunc: func(request CreateCriterionRequest) (*Criterion, error) {
			return &Criterion{
				ID:          1,
				Name:        request.Name,
				Description: request.Description,
				MaxScore:    request.MaxScore,
				Active:      true,
			}, nil
		},
	}

	handler := NewHandler(service)

	router := gin.New()
	router.POST("/criteria", handler.Create)

	body := CreateCriterionRequest{
		Name:        "Clarté",
		Description: "Clarté des explications",
		MaxScore:    5,
	}

	recorder := performRequest(
		router,
		http.MethodPost,
		"/criteria",
		body,
	)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Clarté")
}

func TestHandler_Create_InvalidRequest(t *testing.T) {
	service := &fakeCriterionService{}

	handler := NewHandler(service)

	router := gin.New()
	router.POST("/criteria", handler.Create)

	recorder := performRequest(
		router,
		http.MethodPost,
		"/criteria",
		map[string]interface{}{
			"name": "Clarté",
		},
	)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid request")
}

func TestHandler_Create_ServiceError(t *testing.T) {
	service := &fakeCriterionService{
		createFunc: func(request CreateCriterionRequest) (*Criterion, error) {
			return nil, ErrCriterionNameRequired
		},
	}

	handler := NewHandler(service)

	router := gin.New()
	router.POST("/criteria", handler.Create)

	body := CreateCriterionRequest{
		Name:     "Clarté",
		MaxScore: 5,
	}

	recorder := performRequest(
		router,
		http.MethodPost,
		"/criteria",
		body,
	)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "criterion name is required")
}

// ---------------------------------------------------------
// Update
// ---------------------------------------------------------

func TestHandler_Update_Success(t *testing.T) {
	service := &fakeCriterionService{
		updateFunc: func(
			id int,
			request UpdateCriterionRequest,
		) (*Criterion, error) {

			return &Criterion{
				ID:          id,
				Name:        request.Name,
				Description: request.Description,
				MaxScore:    request.MaxScore,
				Active:      request.Active,
			}, nil
		},
	}

	handler := NewHandler(service)

	router := gin.New()
	router.PUT("/criteria/:id", handler.Update)

	body := UpdateCriterionRequest{
		Name:        "Clarté",
		Description: "Bonne clarté",
		MaxScore:    5,
		Active:      true,
	}

	recorder := performRequest(
		router,
		http.MethodPut,
		"/criteria/1",
		body,
	)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Clarté")
}

func TestHandler_Update_InvalidID(t *testing.T) {
	service := &fakeCriterionService{}

	handler := NewHandler(service)

	router := gin.New()
	router.PUT("/criteria/:id", handler.Update)

	body := UpdateCriterionRequest{
		Name:     "Clarté",
		MaxScore: 5,
	}

	recorder := performRequest(
		router,
		http.MethodPut,
		"/criteria/abc",
		body,
	)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid criterion ID")
}

func TestHandler_Update_ZeroID(t *testing.T) {
	service := &fakeCriterionService{}

	handler := NewHandler(service)

	router := gin.New()
	router.PUT("/criteria/:id", handler.Update)

	body := UpdateCriterionRequest{
		Name:     "Clarté",
		MaxScore: 5,
	}

	recorder := performRequest(
		router,
		http.MethodPut,
		"/criteria/0",
		body,
	)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid criterion ID")
}

func TestHandler_Update_InvalidRequest(t *testing.T) {
	service := &fakeCriterionService{}

	handler := NewHandler(service)

	router := gin.New()
	router.PUT("/criteria/:id", handler.Update)

	recorder := performRequest(
		router,
		http.MethodPut,
		"/criteria/1",
		map[string]interface{}{
			"name": "Clarté",
		},
	)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid request")
}

func TestHandler_Update_NotFound(t *testing.T) {
	service := &fakeCriterionService{
		updateFunc: func(
			id int,
			request UpdateCriterionRequest,
		) (*Criterion, error) {
			return nil, ErrCriterionNotFound
		},
	}

	handler := NewHandler(service)

	router := gin.New()
	router.PUT("/criteria/:id", handler.Update)

	body := UpdateCriterionRequest{
		Name:     "Clarté",
		MaxScore: 5,
		Active:   true,
	}

	recorder := performRequest(
		router,
		http.MethodPut,
		"/criteria/99",
		body,
	)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "criterion not found")
}

func TestHandler_Update_ServiceError(t *testing.T) {
	service := &fakeCriterionService{
		updateFunc: func(
			id int,
			request UpdateCriterionRequest,
		) (*Criterion, error) {
			return nil, ErrMaxScoreInvalid
		},
	}

	handler := NewHandler(service)

	router := gin.New()
	router.PUT("/criteria/:id", handler.Update)

	body := UpdateCriterionRequest{
		Name:     "Clarté",
		MaxScore: 0,
		Active:   true,
	}

	recorder := performRequest(
		router,
		http.MethodPut,
		"/criteria/1",
		body,
	)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "max score must be greater than zero")
}

// ---------------------------------------------------------
// Delete
// ---------------------------------------------------------

func TestHandler_Delete_Success(t *testing.T) {
	service := &fakeCriterionService{
		deleteFunc: func(id int) error {
			return nil
		},
	}

	handler := NewHandler(service)

	router := gin.New()
	router.DELETE("/criteria/:id", handler.Delete)

	recorder := performRequest(
		router,
		http.MethodDelete,
		"/criteria/1",
		nil,
	)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(
		t,
		recorder.Body.String(),
		"criterion deleted successfully",
	)
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	service := &fakeCriterionService{}

	handler := NewHandler(service)

	router := gin.New()
	router.DELETE("/criteria/:id", handler.Delete)

	recorder := performRequest(
		router,
		http.MethodDelete,
		"/criteria/abc",
		nil,
	)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid criterion ID")
}

func TestHandler_Delete_ZeroID(t *testing.T) {
	service := &fakeCriterionService{}

	handler := NewHandler(service)

	router := gin.New()
	router.DELETE("/criteria/:id", handler.Delete)

	recorder := performRequest(
		router,
		http.MethodDelete,
		"/criteria/0",
		nil,
	)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid criterion ID")
}

func TestHandler_Delete_NotFound(t *testing.T) {
	service := &fakeCriterionService{
		deleteFunc: func(id int) error {
			return ErrCriterionNotFound
		},
	}

	handler := NewHandler(service)

	router := gin.New()
	router.DELETE("/criteria/:id", handler.Delete)

	recorder := performRequest(
		router,
		http.MethodDelete,
		"/criteria/99",
		nil,
	)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "criterion not found")
}

func TestHandler_Delete_InternalError(t *testing.T) {
	service := &fakeCriterionService{
		deleteFunc: func(id int) error {
			return errors.New("database error")
		},
	}

	handler := NewHandler(service)

	router := gin.New()
	router.DELETE("/criteria/:id", handler.Delete)

	recorder := performRequest(
		router,
		http.MethodDelete,
		"/criteria/1",
		nil,
	)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(
		t,
		recorder.Body.String(),
		"failed to delete criterion",
	)
}
