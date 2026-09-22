package auth

import (
	"errors"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Login(request LoginRequest) (*User, error) {
	// Normaliser l'email
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	// Vérifier les données
	if request.Email == "" || request.Password == "" {
		return nil, errors.New("invalid credentials")
	}

	// Rechercher l'utilisateur
	user, err := s.repository.FindByEmail(request.Email)
	if err != nil {
		// Ne jamais révéler si l'email existe ou non
		return nil, errors.New("invalid credentials")
	}

	// Vérifier le mot de passe
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(request.Password),
	); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *Service) Register(request RegisterRequest) (*User, error) {
	// Normaliser les données
	request.Matricule = strings.TrimSpace(request.Matricule)
	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))
	request.Role = strings.ToUpper(strings.TrimSpace(request.Role))

	// Validation
	if request.Matricule == "" {
		return nil, errors.New("matricule is required")
	}

	if request.Name == "" {
		return nil, errors.New("name is required")
	}

	if request.Email == "" {
		return nil, errors.New("email is required")
	}

	if request.Password == "" {
		return nil, errors.New("password is required")
	}

	if len(request.Password) < 8 {
		return nil, errors.New("password must contain at least 8 characters")
	}

	// Validation email
	if _, err := mail.ParseAddress(request.Email); err != nil {
		return nil, errors.New("invalid email")
	}

	// Rôle par défaut
	if request.Role == "" {
		request.Role = "STUDENT"
	}

	// Vérification du rôle
	switch request.Role {
	case "STUDENT", "PROFESSOR", "ADMIN":
		// rôle valide

	default:
		return nil, errors.New(
			"invalid role. Allowed roles: STUDENT, PROFESSOR, ADMIN",
		)
	}

	// Vérifier le matricule
	existingUser, err := s.repository.FindByMatricule(request.Matricule)

	if err == nil && existingUser != nil {
		return nil, errors.New("matricule already exists")
	}

	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	// Vérifier l'email
	existingUser, err = s.repository.FindByEmail(request.Email)

	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	// Hasher le mot de passe
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Construire l'utilisateur
	user := User{
		Matricule:    request.Matricule,
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
		Role:         request.Role,
	}

	// Enregistrer l'utilisateur
	createdUser, err := s.repository.Create(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
