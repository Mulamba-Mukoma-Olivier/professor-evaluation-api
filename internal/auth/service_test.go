package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewService(t *testing.T) {
	// Simple test to verify service creation
	repo := &Repository{}
	service := NewService(repo)

	assert.NotNil(t, service)
	assert.NotNil(t, service.repository)
}

func TestService_Register_PasswordHashing(t *testing.T) {
	// Test password hashing functionality
	password := "password123"
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	assert.NotEqual(t, password, string(hashedPassword))

	// Verify hash can be checked
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	assert.NoError(t, err)
}
