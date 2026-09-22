package criteria

import "github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"

type Criterion struct {
	ID                int                            `json:"id" gorm:"primaryKey"`
	Name              string                         `json:"name" gorm:"not null"`
	Description       string                         `json:"description"`
	MaxScore          int                            `json:"max_score" gorm:"not null"`
	Active            bool                           `json:"active"`
	EvaluationAnswers []evaluations.EvaluationAnswer `json:"evaluation_answers,omitempty" gorm:"foreignKey:CriterionID;references:ID"`
}

type CreateCriterionRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=150"`
	Description string `json:"description"`
	MaxScore    int    `json:"max_score" binding:"required,min=1"`
}

type UpdateCriterionRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=150"`
	Description string `json:"description"`
	MaxScore    int    `json:"max_score" binding:"required,min=1"`
	Active      bool   `json:"active"`
}
