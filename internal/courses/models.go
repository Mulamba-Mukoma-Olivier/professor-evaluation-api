package courses

import "github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"

type Course struct {
	ID           int                         `json:"id" gorm:"primaryKey"`
	Code         string                      `json:"code" gorm:"uniqueIndex;not null"`
	Name         string                      `json:"name" gorm:"not null"`
	Description  string                      `json:"description"`
	Department   string                      `json:"department"`
	AcademicYear string                      `json:"academic_year"`
	Evaluations  []evaluations.Evaluation    `json:"evaluations,omitempty" gorm:"foreignKey:CourseID;references:ID"`
}

type CreateCourseRequest struct {
	Code         string `json:"code" binding:"required,min=2,max=20"`
	Name         string `json:"name" binding:"required,min=2,max=150"`
	Description  string `json:"description"`
	Department   string `json:"department" binding:"required,min=2,max=100"`
	AcademicYear string `json:"academic_year" binding:"required,min=4,max=20"`
}

type UpdateCourseRequest struct {
	Code         string `json:"code" binding:"required,min=2,max=20"`
	Name         string `json:"name" binding:"required,min=2,max=150"`
	Description  string `json:"description"`
	Department   string `json:"department" binding:"required,min=2,max=100"`
	AcademicYear string `json:"academic_year" binding:"required,min=4,max=20"`
}
