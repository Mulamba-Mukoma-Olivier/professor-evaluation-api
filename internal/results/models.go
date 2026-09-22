package results

type CriterionResult struct {
	CriterionID int     `json:"criterion_id"`
	Average     float64 `json:"average"`
	Responses   int     `json:"responses"`
}

type ProfessorResult struct {
	ProfessorID   int               `json:"professor_id"`
	CourseID      int               `json:"course_id"`
	AcademicYear  string            `json:"academic_year"`
	Period        string            `json:"period"`
	TotalReviews  int               `json:"total_reviews"`
	GlobalAverage float64           `json:"global_average"`
	Criteria      []CriterionResult `json:"criteria"`
}
