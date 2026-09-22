package auth

import (
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/eligibility"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"
)

// User représente un utilisateur de l'application.
type User struct {
	ID           int                      `json:"id" gorm:"primaryKey"`
	Matricule    string                   `json:"matricule" gorm:"not null"`
	Name         string                   `json:"name" gorm:"not null"`
	Email        string                   `json:"email" gorm:"not null"`
	PasswordHash string                   `json:"-"`
	Role         string                   `json:"role" gorm:"not null"`
	Eligibility  *eligibility.Eligibility `json:"eligibility,omitempty" gorm:"foreignKey:StudentID;references:ID"`
	Evaluations  []evaluations.Evaluation `json:"evaluations,omitempty" gorm:"foreignKey:StudentID;references:ID"`
}

// LoginRequest représente les données envoyées lors de la connexion.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse représente la réponse après une connexion réussie.
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	User        User   `json:"user"`
}

// RegisterRequest représente les données nécessaires pour créer un utilisateur.
type RegisterRequest struct {
	Matricule string `json:"matricule" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Role      string `json:"role" binding:"required,oneof=STUDENT PROFESSOR ADMIN"`
}
