package eligibility

type Eligibility struct {
	ID             int      `json:"id" gorm:"primaryKey"`
	StudentID      int      `json:"student_id" gorm:"not null;uniqueIndex"`
	Enrollment     bool     `json:"enrollment" gorm:"not null"`
	AcademicFees   bool     `json:"academic_fees" gorm:"not null"`
	LaboratoryFees bool     `json:"laboratory_fees" gorm:"not null"`
	AccessFees     bool     `json:"access_fees" gorm:"not null"`
	Eligible       bool     `json:"eligible" gorm:"not null"`
	Reasons        []string `json:"reasons" gorm:"-"`
}

type CreateEligibilityRequest struct {
	StudentID      int  `json:"student_id" binding:"required"`
	Enrollment     bool `json:"enrollment"`
	AcademicFees   bool `json:"academic_fees"`
	LaboratoryFees bool `json:"laboratory_fees"`
	AccessFees     bool `json:"access_fees"`
}

type UpdateEligibilityRequest struct {
	Enrollment     bool `json:"enrollment"`
	AcademicFees   bool `json:"academic_fees"`
	LaboratoryFees bool `json:"laboratory_fees"`
	AccessFees     bool `json:"access_fees"`
}