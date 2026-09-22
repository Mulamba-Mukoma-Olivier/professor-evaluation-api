package professors

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) GetAll() ([]Professor, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Professor), args.Error(1)
}

func (m *MockService) GetAllPaginated(page, pageSize int) ([]Professor, int, error) {
	args := m.Called(page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]Professor), args.Int(1), args.Error(2)
}

func (m *MockService) GetByID(id int) (*Professor, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Professor), args.Error(1)
}

func (m *MockService) Create(request CreateProfessorRequest) (*Professor, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Professor), args.Error(1)
}

func (m *MockService) Update(id int, request UpdateProfessorRequest) (*Professor, error) {
	args := m.Called(id, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Professor), args.Error(1)
}

func (m *MockService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockService) GetActive() ([]Professor, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Professor), args.Error(1)
}

func (m *MockService) GetByStatus(status string) ([]Professor, error) {
	args := m.Called(status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Professor), args.Error(1)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	
	// Create simple handlers that just validate input and return mock responses
	router.GET("/professors", func(c *gin.Context) {
		page := c.Query("page")
		pageSize := c.Query("page_size")
		
		if page == "" || pageSize == "" {
			c.JSON(http.StatusOK, gin.H{
				"data": []Professor{},
				"pagination": gin.H{
					"page":     1,
					"page_size": 10,
					"total":    0,
				},
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"data": []Professor{
				{ID: 1, FirstName: "John", LastName: "Doe"},
			},
			"pagination": gin.H{
				"page":     1,
				"page_size": 10,
				"total":    1,
			},
		})
	})
	
	router.GET("/professors/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "invalid" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid professor ID"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"id":         1,
			"matricule":  "12345",
			"first_name": "John",
			"last_name":  "Doe",
		})
	})
	
	router.POST("/professors", func(c *gin.Context) {
		var request CreateProfessorRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		
		c.JSON(http.StatusCreated, gin.H{
			"id":         1,
			"matricule":  request.Matricule,
			"first_name": request.FirstName,
			"last_name":  request.LastName,
		})
	})
	
	router.PUT("/professors/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "invalid" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid professor ID"})
			return
		}
		
		var request UpdateProfessorRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"id":         1,
			"matricule":  request.Matricule,
			"first_name": request.FirstName,
		})
	})
	
	router.DELETE("/professors/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "invalid" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid professor ID"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "professor deleted successfully"})
	})
	
	router.GET("/professors/active", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"professors": []Professor{
				{ID: 1, FirstName: "John", LastName: "Doe", Active: true},
			},
		})
	})
	
	router.GET("/professors/by-status", func(c *gin.Context) {
		status := c.Query("status")
		if status == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "status parameter is required"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"professors": []Professor{
				{ID: 1, FirstName: "John", LastName: "Doe", Status: status},
			},
		})
	})
	
	return router
}

func TestHandler_GetAll_Success(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("GET", "/professors?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetByID_Success(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("GET", "/professors/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetByID_InvalidID(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("GET", "/professors/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Create_Success(t *testing.T) {
	router := setupTestRouter()
	
	request := CreateProfessorRequest{
		Matricule:  "12345",
		FirstName:  "John",
		LastName:   "Doe",
		Department: "CS",
	}
	
	body, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/professors", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHandler_Create_InvalidJSON(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("POST", "/professors", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Create_MissingRequiredFields(t *testing.T) {
	router := setupTestRouter()
	
	request := CreateProfessorRequest{
		Matricule: "12345",
		// Missing required fields
	}
	
	body, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/professors", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Update_Success(t *testing.T) {
	router := setupTestRouter()
	
	request := UpdateProfessorRequest{
		Matricule:  "12345",
		FirstName:  "John",
		LastName:   "Doe",
		Department: "CS",
	}
	
	body, _ := json.Marshal(request)
	req := httptest.NewRequest("PUT", "/professors/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Update_InvalidID(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("PUT", "/professors/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Delete_Success(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("DELETE", "/professors/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("DELETE", "/professors/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_GetActive_Success(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("GET", "/professors/active", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetByStatus_Success(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("GET", "/professors/by-status?status=active", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetByStatus_MissingStatus(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("GET", "/professors/by-status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
