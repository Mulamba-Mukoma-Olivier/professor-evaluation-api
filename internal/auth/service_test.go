package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestNewService(t *testing.T) {
	repo := &Repository{}

	service := NewService(repo)

	require.NotNil(t, service)
	assert.Equal(t, repo, service.repository)
}

func TestPasswordHashing(t *testing.T) {
	password := "password123"

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	// Le mot de passe en clair ne doit jamais être égal au hash.
	assert.NotEqual(t, password, string(hashedPassword))

	// Le hash doit permettre de vérifier le mot de passe original.
	err = bcrypt.CompareHashAndPassword(
		hashedPassword,
		[]byte(password),
	)

	assert.NoError(t, err)
}

func TestPasswordHashing_WrongPassword(t *testing.T) {
	password := "password123"
	wrongPassword := "wrongpassword"

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	require.NoError(t, err)

	err = bcrypt.CompareHashAndPassword(
		hashedPassword,
		[]byte(wrongPassword),
	)

	assert.Error(t, err)
}
