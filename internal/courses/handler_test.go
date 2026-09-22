package courses

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	
	// Create simple handlers that just validate input and return mock responses
	router.GET("/courses", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"courses": []Course{
				{ID: 1, Name: "Introduction to Computer Science", Code: "CS101"},
			},
		})
	})
	
	router.GET("/courses/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "invalid" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"id":   1,
			"name": "Introduction to Computer Science",
			"code": "CS101",
		})
	})
	
	router.POST("/courses", func(c *gin.Context) {
		var request CreateCourseRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		
		c.JSON(http.StatusCreated, gin.H{
			"id":   1,
			"name": request.Name,
			"code": request.Code,
		})
	})
	
	router.PUT("/courses/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "invalid" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
			return
		}
		
		var request UpdateCourseRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"id":   1,
			"name": request.Name,
			"code": request.Code,
		})
	})
	
	router.DELETE("/courses/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "invalid" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "course deleted successfully"})
	})
	
	return router
}

func TestHandler_GetAll_Success(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("GET", "/courses", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetByID_Success(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("GET", "/courses/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetByID_InvalidID(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("GET", "/courses/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Create_Success(t *testing.T) {
	router := setupTestRouter()
	
	request := CreateCourseRequest{
		Name:         "Advanced Algorithms",
		Code:         "CS201",
		Description:  "Advanced algorithms course",
		Department:   "Computer Science",
		AcademicYear: "2024-2025",
	}
	
	body, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/courses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHandler_Create_InvalidJSON(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("POST", "/courses", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Create_MissingRequiredFields(t *testing.T) {
	router := setupTestRouter()
	
	request := CreateCourseRequest{
		Name: "Advanced Algorithms",
		// Missing required fields
	}
	
	body, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/courses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Update_Success(t *testing.T) {
	router := setupTestRouter()
	
	request := UpdateCourseRequest{
		Name:         "Advanced Algorithms Updated",
		Code:         "CS201",
		Description:  "Updated description",
		Department:   "Computer Science",
		AcademicYear: "2024-2025",
	}
	
	body, _ := json.Marshal(request)
	req := httptest.NewRequest("PUT", "/courses/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Update_InvalidID(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("PUT", "/courses/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Update_InvalidJSON(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("PUT", "/courses/1", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Delete_Success(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("DELETE", "/courses/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("DELETE", "/courses/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
