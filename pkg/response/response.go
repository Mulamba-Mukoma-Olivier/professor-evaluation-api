package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse représente le format standard des réponses de l'API.
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

// Success retourne une réponse HTTP réussie.
func Success(
	c *gin.Context,
	status int,
	message string,
	data any,
) {
	c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error retourne une réponse HTTP en erreur.
func Error(
	c *gin.Context,
	status int,
	message string,
	err any,
) {
	c.JSON(status, APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// BadRequest retourne une erreur 400.
func BadRequest(c *gin.Context, message string, err any) {
	Error(c, http.StatusBadRequest, message, err)
}

// Unauthorized retourne une erreur 401.
func Unauthorized(c *gin.Context, message string, err any) {
	Error(c, http.StatusUnauthorized, message, err)
}

// Forbidden retourne une erreur 403.
func Forbidden(c *gin.Context, message string, err any) {
	Error(c, http.StatusForbidden, message, err)
}

// NotFound retourne une erreur 404.
func NotFound(c *gin.Context, message string, err any) {
	Error(c, http.StatusNotFound, message, err)
}

// Conflict retourne une erreur 409.
func Conflict(c *gin.Context, message string, err any) {
	Error(c, http.StatusConflict, message, err)
}

// UnprocessableEntity retourne une erreur 422.
func UnprocessableEntity(c *gin.Context, message string, err any) {
	Error(c, http.StatusUnprocessableEntity, message, err)
}

// InternalServerError retourne une erreur 500.
func InternalServerError(c *gin.Context, message string, err any) {
	Error(c, http.StatusInternalServerError, message, err)
}
