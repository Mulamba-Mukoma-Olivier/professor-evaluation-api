package professors

import "github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"

type Professor struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	Matricule    string `json:"matricule" gorm:"uniqueIndex;not null"`
	FirstName    string `json:"first_name" gorm:"not null"`
	LastName     string `json:"last_name" gorm:"not null"`
	Email        string `json:"email"`
	Department   string `json:"department" gorm:"not null"`
	Grade        string `json:"grade"`
	Active       bool   `json:"active"`
	Status       string `json:"status" gorm:"default:'active'"` // active, inactive, on_leave, retired
	Evaluations []evaluations.Evaluation `json:"evaluations" gorm:"foreignKey:ProfessorID;references:ID"`
}

type CreateProfessorRequest struct {
	Matricule  string `json:"matricule" binding:"required"`
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Email      string `json:"email"`
	Department string `json:"department" binding:"required"`
	Grade      string `json:"grade"`
	Status     string `json:"status" gorm:"default:'active'"`
}

type UpdateProfessorRequest struct {
	Matricule  string `json:"matricule" binding:"required"`
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Email      string `json:"email"`
	Department string `json:"department" binding:"required"`
	Grade      string `json:"grade"`
	Active     bool   `json:"active"`
	Status     string `json:"status"`
}