package results

import (
	"errors"
	"testing"

	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FakeResultsRepository struct {
	evaluations []evaluations.Evaluation
	err         error
}

func (f *FakeResultsRepository) GetByProfessor(
	professorID int,
	courseID int,
	academicYear string,
	period string,
) ([]evaluations.Evaluation, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.evaluations, nil
}

func TestService_GetProfessorResult_Success(t *testing.T) {
	repository := &FakeResultsRepository{
		evaluations: []evaluations.Evaluation{
			{
				ID:          1,
				StudentID:   100,
				ProfessorID: 1,
				CourseID:    10,
				AcademicYear: "2025-2026",
				Period:      "semester_1",
				Answers: []evaluations.EvaluationAnswer{
					{
						CriterionID: 1,
						Score:       4,
					},
					{
						CriterionID: 2,
						Score:       5,
					},
				},
			},
			{
				ID:          2,
				StudentID:   101,
				ProfessorID: 1,
				CourseID:    10,
				AcademicYear: "2025-2026",
				Period:      "semester_1",
				Answers: []evaluations.EvaluationAnswer{
					{
						CriterionID: 1,
						Score:       3,
					},
					{
						CriterionID: 2,
						Score:       4,
					},
				},
			},
		},
	}

	service := NewService(repository)

	result, err := service.GetProfessorResult(
		1,
		10,
		"2025-2026",
		"semester_1",
	)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, result.ProfessorID)
	assert.Equal(t, 10, result.CourseID)
	assert.Equal(t, "2025-2026", result.AcademicYear)
	assert.Equal(t, "semester_1", result.Period)

	assert.Equal(t, 2, result.TotalReviews)
	assert.Equal(t, 4.0, result.GlobalAverage)

	require.Len(t, result.Criteria, 2)

	// Les critères doivent être triés par ID.
	assert.Equal(t, 1, result.Criteria[0].CriterionID)
	assert.Equal(t, 3.5, result.Criteria[0].Average)
	assert.Equal(t, 2, result.Criteria[0].Responses)

	assert.Equal(t, 2, result.Criteria[1].CriterionID)
	assert.Equal(t, 4.5, result.Criteria[1].Average)
	assert.Equal(t, 2, result.Criteria[1].Responses)
}

func TestService_GetProfessorResult_InvalidProfessorID(t *testing.T) {
	repository := &FakeResultsRepository{}
	service := NewService(repository)

	result, err := service.GetProfessorResult(
		0,
		10,
		"2025-2026",
		"semester_1",
	)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidProfessorID)
}

func TestService_GetProfessorResult_InvalidCourseID(t *testing.T) {
	repository := &FakeResultsRepository{}
	service := NewService(repository)

	result, err := service.GetProfessorResult(
		1,
		0,
		"2025-2026",
		"semester_1",
	)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidCourseID)
}

func TestService_GetProfessorResult_AcademicYearRequired(t *testing.T) {
	repository := &FakeResultsRepository{}
	service := NewService(repository)

	result, err := service.GetProfessorResult(
		1,
		10,
		"",
		"semester_1",
	)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrAcademicYearRequired)
}

func TestService_GetProfessorResult_AcademicYearWhitespace(t *testing.T) {
	repository := &FakeResultsRepository{}
	service := NewService(repository)

	result, err := service.GetProfessorResult(
		1,
		10,
		"   ",
		"semester_1",
	)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrAcademicYearRequired)
}

func TestService_GetProfessorResult_PeriodRequired(t *testing.T) {
	repository := &FakeResultsRepository{}
	service := NewService(repository)

	result, err := service.GetProfessorResult(
		1,
		10,
		"2025-2026",
		"",
	)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrPeriodRequired)
}

func TestService_GetProfessorResult_PeriodWhitespace(t *testing.T) {
	repository := &FakeResultsRepository{}
	service := NewService(repository)

	result, err := service.GetProfessorResult(
		1,
		10,
		"2025-2026",
		"   ",
	)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrPeriodRequired)
}

func TestService_GetProfessorResult_RepositoryError(t *testing.T) {
	repositoryError := errors.New("database connection failed")

	repository := &FakeResultsRepository{
		err: repositoryError,
	}

	service := NewService(repository)

	result, err := service.GetProfessorResult(
		1,
		10,
		"2025-2026",
		"semester_1",
	)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repositoryError)
}

func TestService_GetProfessorResult_NoEvaluations(t *testing.T) {
	repository := &FakeResultsRepository{
		evaluations: []evaluations.Evaluation{},
	}

	service := NewService(repository)

	result, err := service.GetProfessorResult(
		1,
		10,
		"2025-2026",
		"semester_1",
	)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrNoEvaluations)
}

func TestService_GetProfessorResult_NoAnswers(t *testing.T) {
	repository := &FakeResultsRepository{
		evaluations: []evaluations.Evaluation{
			{
				ID:          1,
				ProfessorID: 1,
				CourseID:    10,
				AcademicYear: "2025-2026",
				Period:      "semester_1",
				Answers:     []evaluations.EvaluationAnswer{},
			},
		},
	}

	service := NewService(repository)

	result, err := service.GetProfessorResult(
		1,
		10,
		"2025-2026",
		"semester_1",
	)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrNoAnswers)
}

func TestService_GetProfessorResult_IgnoresInvalidCriterionID(t *testing.T) {
	repository := &FakeResultsRepository{
		evaluations: []evaluations.Evaluation{
			{
				ID: 1,
				Answers: []evaluations.EvaluationAnswer{
					{
						CriterionID: 0,
						Score:       5,
					},
					{
						CriterionID: 1,
						Score:       4,
					},
				},
			},
		},
	}

	service := NewService(repository)

	result, err := service.GetProfessorResult(
		1,
		10,
		"2025-2026",
		"semester_1",
	)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, result.Criteria[0].CriterionID)
	assert.Equal(t, 4.0, result.GlobalAverage)
	assert.Equal(t, 1, result.Criteria[0].Responses)
}
