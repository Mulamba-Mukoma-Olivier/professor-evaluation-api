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
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))
	if request.Email == "" {
		return nil, errors.New("invalid credentials")
	}

	user, err := s.repository.FindByEmail(request.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(request.Password),
	)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *Service) Register(request RegisterRequest) (*User, error) {
	// check existing user
	if user, err := s.repository.FindByMatricule(request.Matricule); err == nil && user != nil {
		return nil, errors.New("user already exists")
	}
	if user, err := s.repository.FindByEmail(request.Email); err == nil && user != nil {
		return nil, errors.New("user already exists")
	}

	// normalize and validate email
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))
	if _, err := mail.ParseAddress(request.Email); err != nil {
		return nil, errors.New("invalid email")
	}

	// default role
	if strings.TrimSpace(request.Role) == "" {
		request.Role = "STUDENT"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := User{
		Matricule:    request.Matricule,
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(hash),
		Role:         request.Role,
	}

	created, err := s.repository.Create(user)
	if err != nil {
		return nil, err
	}

	return created, nil
}
