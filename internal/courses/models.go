package courses

import "github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"

type Course struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	Code         string `json:"code" gorm:"uniqueIndex;not null"`
	Name         string `json:"name" gorm:"not null"`
	Description  string `json:"description"`
	Department   string `json:"department"`
	AcademicYear string `json:"academic_year"`
	Evaluations []evaluations.Evaluation `json:"evaluations" gorm:"foreignKey:CourseID;references:ID"`
}

type CreateCourseRequest struct {
	Code         string `json:"code" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	Department   string `json:"department" binding:"required"`
	AcademicYear string `json:"academic_year" binding:"required"`
}

type UpdateCourseRequest struct {
	Code         string `json:"code" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	Department   string `json:"department" binding:"required"`
	AcademicYear string `json:"academic_year" binding:"required"`
}