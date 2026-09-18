package validator

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Validate valide une structure à l'aide des tags
// du package go-playground/validator.
func Validate(data any) error {
	if data == nil {
		return errors.New("data is required")
	}

	return validate.Struct(data)
}

// ValidateRequired vérifie qu'une chaîne
// contient une valeur non vide.
func ValidateRequired(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("field is required")
	}

	return nil
}

// ValidatePassword vérifie la longueur minimale
// du mot de passe.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New(
			"password must contain at least 8 characters",
		)
	}

	return nil
}

// ValidateScore vérifie qu'une note est comprise
// entre 1 et 5.
func ValidateScore(score int) error {
	if score < 1 || score > 5 {
		return errors.New(
			"score must be between 1 and 5",
		)
	}

	return nil
}