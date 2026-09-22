package professors

import "github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"

type Professor struct {
	ID           int                    `json:"id" gorm:"primaryKey"`
	Matricule    string                 `json:"matricule" gorm:"uniqueIndex;not null"`
	FirstName    string                 `json:"first_name" gorm:"not null"`
	LastName     string                 `json:"last_name" gorm:"not null"`
	Email        string                 `json:"email"`
	Department   string                 `json:"department" gorm:"not null"`
	Grade        string                 `json:"grade"`
	Active       bool                   `json:"active"`
	Status       string                 `json:"status" gorm:"default:'active'"`
	Evaluations  []evaluations.Evaluation `json:"evaluations,omitempty" gorm:"foreignKey:ProfessorID;references:ID"`
}

type CreateProfessorRequest struct {
	Matricule  string `json:"matricule" binding:"required,min=2,max=50"`
	FirstName  string `json:"first_name" binding:"required,min=2,max=100"`
	LastName   string `json:"last_name" binding:"required,min=2,max=100"`
	Email      string `json:"email" binding:"omitempty,email"`
	Department string `json:"department" binding:"required,min=2,max=150"`
	Grade      string `json:"grade" binding:"omitempty,max=100"`
	Status     string `json:"status" binding:"omitempty,oneof=active inactive on_leave retired"`
}

type UpdateProfessorRequest struct {
	Matricule  string `json:"matricule" binding:"required,min=2,max=50"`
	FirstName  string `json:"first_name" binding:"required,min=2,max=100"`
	LastName   string `json:"last_name" binding:"required,min=2,max=100"`
	Email      string `json:"email" binding:"omitempty,email"`
	Department string `json:"department" binding:"required,min=2,max=150"`
	Grade      string `json:"grade" binding:"omitempty,max=100"`
	Active     bool   `json:"active"`
	Status     string `json:"status" binding:"omitempty,oneof=active inactive on_leave retired"`
}
