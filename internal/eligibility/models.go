package eligibility

type Eligibility struct {
	ID               int      `json:"id" gorm:"primaryKey"`
	StudentID        int      `json:"student_id" gorm:"not null;uniqueIndex"`
	Enrollment       bool     `json:"enrollment" gorm:"not null"`
	AcademicFees     bool     `json:"academic_fees" gorm:"not null"`
	LaboratoryFees   bool     `json:"laboratory_fees" gorm:"not null"`
	AccessFees       bool     `json:"access_fees" gorm:"not null"`
	Eligible         bool     `json:"eligible" gorm:"not null"`
	Reasons          []string `json:"reasons" gorm:"-"`
}