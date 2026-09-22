package eligibility

import "errors"

type EligibilityRepository interface {
	GetByStudentID(studentID int) (*Eligibility, error)
}

type Service struct {
	repository EligibilityRepository
}

func NewService(repository EligibilityRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Check(studentID int) (*Eligibility, error) {
	if studentID <= 0 {
		return nil, errors.New("invalid student ID")
	}

	eligibility, err := s.repository.GetByStudentID(studentID)
	if err != nil {
		return nil, err
	}

	var reasons []string

	if !eligibility.Enrollment {
		reasons = append(reasons, "ENROLLMENT_INVALID")
	}

	if !eligibility.AcademicFees {
		reasons = append(reasons, "ACADEMIC_FEES_NOT_PAID")
	}

	if !eligibility.LaboratoryFees {
		reasons = append(reasons, "LABORATORY_FEES_NOT_PAID")
	}

	if !eligibility.AccessFees {
		reasons = append(reasons, "ACCESS_FEES_NOT_PAID")
	}

	eligibility.Reasons = reasons
	eligibility.Eligible = len(reasons) == 0

	return eligibility, nil
}
