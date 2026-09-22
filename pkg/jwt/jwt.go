package jwt

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrSecretKeyRequired    = errors.New("JWT secret key is required")
	ErrTokenDurationInvalid = errors.New("JWT token duration must be greater than zero")
	ErrTokenRequired        = errors.New("token is required")
	ErrInvalidToken         = errors.New("invalid token")
	ErrUnexpectedMethod     = errors.New("unexpected signing method")
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

// GenerateToken génère un JWT signé avec HS256.
func (m *Manager) GenerateToken(
	userID int,
	matricule string,
	role string,
) (string, error) {

	if strings.TrimSpace(m.SecretKey) == "" {
		return "", ErrSecretKeyRequired
	}

	if m.TokenDuration <= 0 {
		return "", ErrTokenDurationInvalid
	}

	if userID <= 0 {
		return "", errors.New("invalid user ID")
	}

	matricule = strings.TrimSpace(matricule)
	role = strings.TrimSpace(role)

	if matricule == "" {
		return "", errors.New("matricule is required")
	}

	if role == "" {
		return "", errors.New("role is required")
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

// ValidateToken vérifie et décode un JWT.
func (m *Manager) ValidateToken(
	tokenString string,
) (*Claims, error) {

	if strings.TrimSpace(tokenString) == "" {
		return nil, ErrTokenRequired
	}

	if strings.TrimSpace(m.SecretKey) == "" {
		return nil, ErrSecretKeyRequired
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrUnexpectedMethod
			}

			return []byte(m.SecretKey), nil
		},
	)

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
