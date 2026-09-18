package criteria

import "github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"

type Criterion struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	MaxScore    int    `json:"max_score" gorm:"not null"`
	Active      bool   `json:"active"`
	EvaluationAnswers []evaluations.EvaluationAnswer `json:"evaluation_answers" gorm:"foreignKey:CriterionID;references:ID"`
}

type CreateCriterionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	MaxScore    int    `json:"max_score" binding:"required"`
}

type UpdateCriterionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	MaxScore    int    `json:"max_score" binding:"required"`
	Active      bool   `json:"active"`
}