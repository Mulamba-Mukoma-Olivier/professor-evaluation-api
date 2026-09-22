package evaluations

import "time"

// Evaluation représente une évaluation effectuée par un étudiant
// sur un professeur pour un cours donné.
type Evaluation struct {
	ID           int                `json:"id" gorm:"primaryKey"`
	StudentID    int                `json:"student_id" gorm:"not null"`
	ProfessorID  int                `json:"professor_id" gorm:"not null"`
	CourseID     int                `json:"course_id" gorm:"not null"`
	AcademicYear string             `json:"academic_year" gorm:"not null"`
	Period       string             `json:"period" gorm:"not null"`
	Answers      []EvaluationAnswer `json:"answers" gorm:"foreignKey:EvaluationID"`
	SubmittedAt  time.Time          `json:"submitted_at"`
}

// EvaluationAnswer représente la note donnée à un critère.
type EvaluationAnswer struct {
	ID           int `json:"id" gorm:"primaryKey"`
	EvaluationID int `json:"evaluation_id" gorm:"not null"`
	CriterionID  int `json:"criterion_id" gorm:"not null"`
	Score        int `json:"score" gorm:"not null"`
}

// Answer représente une réponse envoyée lors de la création
// d'une évaluation.
type Answer struct {
	CriterionID int `json:"criterion_id" binding:"required"`
	Score       int `json:"score" binding:"required,min=1,max=5"`
}

// CreateEvaluationRequest représente les données nécessaires
// pour créer une évaluation.
type CreateEvaluationRequest struct {
	ProfessorID  int      `json:"professor_id" binding:"required"`
	CourseID     int      `json:"course_id" binding:"required"`
	AcademicYear string   `json:"academic_year" binding:"required"`
	Period       string   `json:"period" binding:"required"`
	Answers      []Answer `json:"answers" binding:"required,min=1"`
}
