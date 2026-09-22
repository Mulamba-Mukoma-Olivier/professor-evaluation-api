package eligibility

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeEligibilityRepository struct {
	getByStudentIDFunc func(studentID int) (*Eligibility, error)
}

func (f *fakeEligibilityRepository) GetByStudentID(studentID int) (*Eligibility, error) {
	return f.getByStudentIDFunc(studentID)
}

func TestService_Check_Success_Eligible(t *testing.T) {
	repository := &fakeEligibilityRepository{
		getByStudentIDFunc: func(studentID int) (*Eligibility, error) {
			return &Eligibility{
				ID:             1,
				StudentID:      studentID,
				Enrollment:     true,
				AcademicFees:   true,
				LaboratoryFees: true,
				AccessFees:     true,
			}, nil
		},
	}

	service := NewService(repository)

	result, err := service.Check(10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Eligible)
	assert.Empty(t, result.Reasons)
}

func TestService_Check_NotEligible_Enrollment(t *testing.T) {
	repository := &fakeEligibilityRepository{
		getByStudentIDFunc: func(studentID int) (*Eligibility, error) {
			return &Eligibility{
				StudentID:      studentID,
				Enrollment:     false,
				AcademicFees:   true,
				LaboratoryFees: true,
				AccessFees:     true,
			}, nil
		},
	}

	service := NewService(repository)

	result, err := service.Check(10)

	assert.NoError(t, err)
	assert.False(t, result.Eligible)
	assert.Equal(t, []string{"ENROLLMENT_INVALID"}, result.Reasons)
}

func TestService_Check_NotEligible_AcademicFees(t *testing.T) {
	repository := &fakeEligibilityRepository{
		getByStudentIDFunc: func(studentID int) (*Eligibility, error) {
			return &Eligibility{
				StudentID:      studentID,
				Enrollment:     true,
				AcademicFees:   false,
				LaboratoryFees: true,
				AccessFees:     true,
			}, nil
		},
	}

	service := NewService(repository)

	result, err := service.Check(10)

	assert.NoError(t, err)
	assert.False(t, result.Eligible)
	assert.Equal(t, []string{"ACADEMIC_FEES_NOT_PAID"}, result.Reasons)
}

func TestService_Check_NotEligible_LaboratoryFees(t *testing.T) {
	repository := &fakeEligibilityRepository{
		getByStudentIDFunc: func(studentID int) (*Eligibility, error) {
			return &Eligibility{
				StudentID:      studentID,
				Enrollment:     true,
				AcademicFees:   true,
				LaboratoryFees: false,
				AccessFees:     true,
			}, nil
		},
	}

	service := NewService(repository)

	result, err := service.Check(10)

	assert.NoError(t, err)
	assert.False(t, result.Eligible)
	assert.Equal(t, []string{"LABORATORY_FEES_NOT_PAID"}, result.Reasons)
}

func TestService_Check_NotEligible_AccessFees(t *testing.T) {
	repository := &fakeEligibilityRepository{
		getByStudentIDFunc: func(studentID int) (*Eligibility, error) {
			return &Eligibility{
				StudentID:      studentID,
				Enrollment:     true,
				AcademicFees:   true,
				LaboratoryFees: true,
				AccessFees:     false,
			}, nil
		},
	}

	service := NewService(repository)

	result, err := service.Check(10)

	assert.NoError(t, err)
	assert.False(t, result.Eligible)
	assert.Equal(t, []string{"ACCESS_FEES_NOT_PAID"}, result.Reasons)
}

func TestService_Check_NotEligible_MultipleReasons(t *testing.T) {
	repository := &fakeEligibilityRepository{
		getByStudentIDFunc: func(studentID int) (*Eligibility, error) {
			return &Eligibility{
				StudentID:      studentID,
				Enrollment:     false,
				AcademicFees:   false,
				LaboratoryFees: false,
				AccessFees:     false,
			}, nil
		},
	}

	service := NewService(repository)

	result, err := service.Check(10)

	assert.NoError(t, err)
	assert.False(t, result.Eligible)

	assert.Equal(t, []string{
		"ENROLLMENT_INVALID",
		"ACADEMIC_FEES_NOT_PAID",
		"LABORATORY_FEES_NOT_PAID",
		"ACCESS_FEES_NOT_PAID",
	}, result.Reasons)
}

func TestService_Check_InvalidStudentID(t *testing.T) {
	repository := &fakeEligibilityRepository{}

	service := NewService(repository)

	result, err := service.Check(0)

	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid student ID")
}

func TestService_Check_RepositoryError(t *testing.T) {
	repository := &fakeEligibilityRepository{
		getByStudentIDFunc: func(studentID int) (*Eligibility, error) {
			return nil, errors.New("database error")
		},
	}

	service := NewService(repository)

	result, err := service.Check(10)

	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")
}

func TestService_Check_NotFound(t *testing.T) {
	repository := &fakeEligibilityRepository{
		getByStudentIDFunc: func(studentID int) (*Eligibility, error) {
			return nil, ErrEligibilityNotFound
		},
	}

	service := NewService(repository)

	result, err := service.Check(10)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrEligibilityNotFound)
}
