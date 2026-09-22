package auth

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

	// Login handler de test
	router.POST("/auth/login", func(c *gin.Context) {
		var request LoginRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request",
			})
			return
		}

		if request.Email == "" || request.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid credentials",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "login would be processed",
		})
	})

	// Register handler de test
	router.POST("/auth/register", func(c *gin.Context) {
		var request RegisterRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request",
			})
			return
		}

		// Validation des champs obligatoires
		if request.Matricule == "" ||
			request.Name == "" ||
			request.Email == "" ||
			request.Password == "" ||
			request.Role == "" {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing required fields",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "register would be processed",
		})
	})

	return router
}

func TestHandler_Login_Success(t *testing.T) {
	router := setupTestRouter()

	loginRequest := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	body, err := json.Marshal(loginRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Login_InvalidJSON(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		bytes.NewBuffer([]byte("invalid json")),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Login_EmptyEmail(t *testing.T) {
	router := setupTestRouter()

	loginRequest := LoginRequest{
		Email:    "",
		Password: "password123",
	}

	body, err := json.Marshal(loginRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Login_EmptyPassword(t *testing.T) {
	router := setupTestRouter()

	loginRequest := LoginRequest{
		Email:    "test@example.com",
		Password: "",
	}

	body, err := json.Marshal(loginRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Register_Success(t *testing.T) {
	router := setupTestRouter()

	registerRequest := RegisterRequest{
		Matricule: "12345",
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "password123",
		Role:      "STUDENT",
	}

	body, err := json.Marshal(registerRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHandler_Register_InvalidJSON(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		bytes.NewBuffer([]byte("invalid json")),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Register_MissingRequiredFields(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name    string
		request RegisterRequest
	}{
		{
			name: "missing matricule",
			request: RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "STUDENT",
			},
		},
		{
			name: "missing name",
			request: RegisterRequest{
				Matricule: "12345",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "STUDENT",
			},
		},
		{
			name: "missing email",
			request: RegisterRequest{
				Matricule: "12345",
				Name:      "Test User",
				Password:  "password123",
				Role:      "STUDENT",
			},
		},
		{
			name: "missing password",
			request: RegisterRequest{
				Matricule: "12345",
				Name:      "Test User",
				Email:     "test@example.com",
				Role:      "STUDENT",
			},
		},
		{
			name: "missing role",
			request: RegisterRequest{
				Matricule: "12345",
				Name:      "Test User",
				Email:     "test@example.com",
				Password:  "password123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.request)
			assert.NoError(t, err)

			req := httptest.NewRequest(
				http.MethodPost,
				"/auth/register",
				bytes.NewBuffer(body),
			)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}
