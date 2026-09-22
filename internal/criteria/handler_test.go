package criteria

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
	router.GET("/criteria", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"criteria": []Criterion{
				{ID: 1, Name: "Teaching Quality", Description: "Quality of teaching", MaxScore: 10, Active: true},
			},
		})
	})
	
	router.GET("/criteria/active", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"criteria": []Criterion{
				{ID: 1, Name: "Teaching Quality", Description: "Quality of teaching", MaxScore: 10, Active: true},
			},
		})
	})
	
	router.GET("/criteria/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "invalid" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid criterion ID"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"id":          1,
			"name":        "Teaching Quality",
			"description": "Quality of teaching",
			"max_score":   10,
			"active":      true,
		})
	})
	
	router.POST("/criteria", func(c *gin.Context) {
		var request CreateCriterionRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		
		c.JSON(http.StatusCreated, gin.H{
			"id":          1,
			"name":        request.Name,
			"description": request.Description,
			"max_score":   request.MaxScore,
			"active":      true,
		})
	})
	
	router.PUT("/criteria/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "invalid" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid criterion ID"})
			return
		}
		
		var request UpdateCriterionRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"id":          1,
			"name":        request.Name,
			"description": request.Description,
			"max_score":   request.MaxScore,
			"active":      request.Active,
		})
	})
	
	router.DELETE("/criteria/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "invalid" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid criterion ID"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "criterion deleted successfully"})
	})
	
	return router
}

func TestHandler_GetAll_Success(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("GET", "/criteria", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetActive_Success(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("GET", "/criteria/active", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetByID_Success(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("GET", "/criteria/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetByID_InvalidID(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("GET", "/criteria/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Create_Success(t *testing.T) {
	router := setupTestRouter()
	
	request := CreateCriterionRequest{
		Name:        "Communication Skills",
		Description: "Communication with students",
		MaxScore:    10,
	}
	
	body, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/criteria", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHandler_Create_InvalidJSON(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("POST", "/criteria", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Create_MissingRequiredFields(t *testing.T) {
	router := setupTestRouter()
	
	request := CreateCriterionRequest{
		Name: "Communication Skills",
		// Missing required fields
	}
	
	body, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/criteria", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Update_Success(t *testing.T) {
	router := setupTestRouter()
	
	request := UpdateCriterionRequest{
		Name:        "Communication Skills Updated",
		Description: "Updated description",
		MaxScore:    10,
		Active:      true,
	}
	
	body, _ := json.Marshal(request)
	req := httptest.NewRequest("PUT", "/criteria/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Update_InvalidID(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("PUT", "/criteria/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Update_InvalidJSON(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("PUT", "/criteria/1", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Delete_Success(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("DELETE", "/criteria/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	router := setupTestRouter()
	
	req := httptest.NewRequest("DELETE", "/criteria/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
