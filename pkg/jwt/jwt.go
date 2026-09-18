package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	SecretKey     string
	TokenDuration time.Duration
}

type Claims struct {
	UserID    int    `json:"user_id"`
	Matricule string `json:"matricule"`
	Role      string `json:"role"`

	jwt.RegisteredClaims
}

func NewManager(
	secretKey string,
	tokenDuration time.Duration,
) *Manager {
	return &Manager{
		SecretKey:     secretKey,
		TokenDuration: tokenDuration,
	}
}

func (m *Manager) GenerateToken(
	userID int,
	matricule string,
	role string,
) (string, error) {

	if m.SecretKey == "" {
		return "", errors.New("JWT secret key is required")
	}

	if m.TokenDuration <= 0 {
		return "", errors.New("JWT token duration must be greater than zero")
	}

	now := time.Now()

	claims := Claims{
		UserID:    userID,
		Matricule: matricule,
		Role:      role,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				now.Add(m.TokenDuration),
			),
			IssuedAt: jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "professor-evaluation-api",
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(m.SecretKey))
}

func (m *Manager) ValidateToken(
	tokenString string,
) (*Claims, error) {

	if tokenString == "" {
		return nil, errors.New("token is required")
	}

	if m.SecretKey == "" {
		return nil, errors.New("JWT secret key is required")
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(m.SecretKey), nil
		},
	)

	if err != nil {
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}