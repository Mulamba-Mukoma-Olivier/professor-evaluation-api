package eligibility

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeEligibilityRepository struct {
	createFunc         func(eligibility *Eligibility) error
	getByStudentIDFunc func(studentID int) (*Eligibility, error)
	updateFunc         func(eligibility *Eligibility) error
	deleteFunc         func(studentID int) error
}

func (f *fakeEligibilityRepository) Create(eligibility *Eligibility) error {
	if f.createFunc != nil {
		return f.createFunc(eligibility)
	}
	return nil
}

func (f *fakeEligibilityRepository) GetByStudentID(studentID int) (*Eligibility, error) {
	if f.getByStudentIDFunc != nil {
		return f.getByStudentIDFunc(studentID)
	}
	return nil, nil
}

func (f *fakeEligibilityRepository) Update(eligibility *Eligibility) error {
	if f.updateFunc != nil {
		return f.updateFunc(eligibility)
	}
	return nil
}

func (f *fakeEligibilityRepository) Delete(studentID int) error {
	if f.deleteFunc != nil {
		return f.deleteFunc(studentID)
	}
	return nil
}

// --------------------------------------------------
// CHECK
// --------------------------------------------------

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

// --------------------------------------------------
// CREATE
// --------------------------------------------------

func TestService_Create_Success(t *testing.T) {
	repository := &fakeEligibilityRepository{
		createFunc: func(eligibility *Eligibility) error {
			eligibility.ID = 1
			return nil
		},
	}

	service := NewService(repository)

	eligibility := &Eligibility{
		StudentID:      10,
		Enrollment:     true,
		AcademicFees:   true,
		LaboratoryFees: true,
		AccessFees:     true,
	}

	result, err := service.Create(eligibility)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.True(t, result.Eligible)
	assert.Empty(t, result.Reasons)
}

func TestService_Create_NotEligible(t *testing.T) {
	repository := &fakeEligibilityRepository{}

	service := NewService(repository)

	eligibility := &Eligibility{
		StudentID:      10,
		Enrollment:     true,
		AcademicFees:   false,
		LaboratoryFees: true,
		AccessFees:     true,
	}

	result, err := service.Create(eligibility)

	assert.NoError(t, err)
	assert.False(t, result.Eligible)
	assert.Equal(t, []string{
		"ACADEMIC_FEES_NOT_PAID",
	}, result.Reasons)
}

func TestService_Create_InvalidStudentID(t *testing.T) {
	repository := &fakeEligibilityRepository{}

	service := NewService(repository)

	result, err := service.Create(&Eligibility{
		StudentID: 0,
	})

	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid student ID")
}

func TestService_Create_NilEligibility(t *testing.T) {
	repository := &fakeEligibilityRepository{}

	service := NewService(repository)

	result, err := service.Create(nil)

	assert.Nil(t, result)
	assert.EqualError(t, err, "eligibility is required")
}

func TestService_Create_RepositoryError(t *testing.T) {
	repository := &fakeEligibilityRepository{
		createFunc: func(eligibility *Eligibility) error {
			return errors.New("database error")
		},
	}

	service := NewService(repository)

	eligibility := &Eligibility{
		StudentID:      10,
		Enrollment:     true,
		AcademicFees:   true,
		LaboratoryFees: true,
		AccessFees:     true,
	}

	result, err := service.Create(eligibility)

	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")
}

// --------------------------------------------------
// UPDATE
// --------------------------------------------------

func TestService_Update_Success(t *testing.T) {
	repository := &fakeEligibilityRepository{
		updateFunc: func(eligibility *Eligibility) error {
			return nil
		},
	}

	service := NewService(repository)

	eligibility := &Eligibility{
		StudentID:      10,
		Enrollment:     true,
		AcademicFees:   false,
		LaboratoryFees: true,
		AccessFees:     true,
	}

	result, err := service.Update(eligibility)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Eligible)
	assert.Equal(t, []string{
		"ACADEMIC_FEES_NOT_PAID",
	}, result.Reasons)
}

func TestService_Update_InvalidStudentID(t *testing.T) {
	repository := &fakeEligibilityRepository{}

	service := NewService(repository)

	result, err := service.Update(&Eligibility{
		StudentID: 0,
	})

	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid student ID")
}

func TestService_Update_NilEligibility(t *testing.T) {
	repository := &fakeEligibilityRepository{}

	service := NewService(repository)

	result, err := service.Update(nil)

	assert.Nil(t, result)
	assert.EqualError(t, err, "eligibility is required")
}

func TestService_Update_NotFound(t *testing.T) {
	repository := &fakeEligibilityRepository{
		updateFunc: func(eligibility *Eligibility) error {
			return ErrEligibilityNotFound
		},
	}

	service := NewService(repository)

	result, err := service.Update(&Eligibility{
		StudentID: 10,
	})

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrEligibilityNotFound)
}

// --------------------------------------------------
// DELETE
// --------------------------------------------------

func TestService_Delete_Success(t *testing.T) {
	repository := &fakeEligibilityRepository{
		deleteFunc: func(studentID int) error {
			return nil
		},
	}

	service := NewService(repository)

	err := service.Delete(10)

	assert.NoError(t, err)
}

func TestService_Delete_InvalidStudentID(t *testing.T) {
	repository := &fakeEligibilityRepository{}

	service := NewService(repository)

	err := service.Delete(0)

	assert.EqualError(t, err, "invalid student ID")
}

func TestService_Delete_NotFound(t *testing.T) {
	repository := &fakeEligibilityRepository{
		deleteFunc: func(studentID int) error {
			return ErrEligibilityNotFound
		},
	}

	service := NewService(repository)

	err := service.Delete(10)

	assert.ErrorIs(t, err, ErrEligibilityNotFound)
}

func TestService_Delete_RepositoryError(t *testing.T) {
	repository := &fakeEligibilityRepository{
		deleteFunc: func(studentID int) error {
			return errors.New("database error")
		},
	}

	service := NewService(repository)

	err := service.Delete(10)

	assert.EqualError(t, err, "database error")
}
