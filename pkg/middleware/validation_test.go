package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{"valid email", "test@example.com", true},
		{"valid email with subdomain", "user@mail.example.com", true},
		{"invalid email - no @", "invalid-email", false},
		{"invalid email - no domain", "test@", false},
		{"invalid email - no local part", "@example.com", false},
		{"invalid email - spaces", "test @example.com", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateEmail(tt.email)
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestValidateMatricule(t *testing.T) {
	tests := []struct {
		name      string
		matricule string
		valid     bool
	}{
		{"valid matricule", "123456", true},
		{"valid matricule with letters", "AB12345", true},
		{"valid matricule max length", "12345678901234567890", true},
		{"invalid matricule - too short", "12345", false},
		{"invalid matricule - too long", "123456789012345678901", false},
		{"invalid matricule - special chars", "123-456", false},
		{"invalid matricule - spaces", "123 456", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateMatricule(tt.matricule)
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal input", "John Doe", "john doe"},
		{"SQL injection - single quote", "John' OR '1'='1", "john or 1=1"},
		{"SQL injection - semicolon", "John; DROP TABLE users", "john  table users"},
		{"SQL injection - comment", "John--", "john"},
		{"SQL injection - exec", "John; EXEC xp_cmdshell", "john  cmdshell"},
		{"mixed dangerous patterns", "John'; DROP TABLE users; --", "john  table users "},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeInput(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
