package results

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupResultsRepositoryTest(t *testing.T) (*Repository, sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	repository := NewRepository(db)

	cleanup := func() {
		sqlDB.Close()
	}

	return repository, mock, cleanup
}

func TestRepository_GetByProfessor_Success(t *testing.T) {
	repository, mock, cleanup := setupResultsRepositoryTest(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "evaluations" WHERE professor_id = $1 AND course_id = $2 AND academic_year = $3 AND period = $4`,
	)).
		WithArgs(1, 10, "2025-2026", "semester_1").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"student_id",
				"professor_id",
				"course_id",
				"academic_year",
				"period",
				"submitted_at",
			}).
				AddRow(
					1,
					100,
					1,
					10,
					"2025-2026",
					"semester_1",
					"2026-09-20 10:00:00",
				).
				AddRow(
					2,
					101,
					1,
					10,
					"2025-2026",
					"semester_1",
					"2026-09-20 11:00:00",
				),
		)

	// GORM fait une deuxième requête pour Preload("Answers").
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "evaluation_answers" WHERE "evaluation_answers"."evaluation_id" IN ($1,$2)`,
	)).
		WithArgs(1, 2).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"evaluation_id",
				"criterion_id",
				"score",
			}).
				AddRow(1, 1, 1, 4).
				AddRow(2, 1, 2, 5).
				AddRow(3, 2, 1, 3).
				AddRow(4, 2, 2, 4),
		)

	results, err := repository.GetByProfessor(
		1,
		10,
		"2025-2026",
		"semester_1",
	)

	require.NoError(t, err)
	require.Len(t, results, 2)

	assert.Equal(t, 1, results[0].ID)
	assert.Equal(t, 100, results[0].StudentID)
	assert.Equal(t, 1, results[0].ProfessorID)
	assert.Equal(t, 10, results[0].CourseID)

	require.Len(t, results[0].Answers, 2)
	assert.Equal(t, 1, results[0].Answers[0].CriterionID)
	assert.Equal(t, 4, results[0].Answers[0].Score)

	require.Len(t, results[1].Answers, 2)
	assert.Equal(t, 2, results[1].Answers[1].CriterionID)
	assert.Equal(t, 4, results[1].Answers[1].Score)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByProfessor_NoResults(t *testing.T) {
	repository, mock, cleanup := setupResultsRepositoryTest(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "evaluations" WHERE professor_id = $1 AND course_id = $2 AND academic_year = $3 AND period = $4`,
	)).
		WithArgs(1, 10, "2025-2026", "semester_1").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"student_id",
				"professor_id",
				"course_id",
				"academic_year",
				"period",
				"submitted_at",
			}),
		)

	results, err := repository.GetByProfessor(
		1,
		10,
		"2025-2026",
		"semester_1",
	)

	require.NoError(t, err)
	assert.Empty(t, results)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByProfessor_DatabaseError(t *testing.T) {
	repository, mock, cleanup := setupResultsRepositoryTest(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "evaluations" WHERE professor_id = $1 AND course_id = $2 AND academic_year = $3 AND period = $4`,
	)).
		WithArgs(1, 10, "2025-2026", "semester_1").
		WillReturnError(assert.AnError)

	results, err := repository.GetByProfessor(
		1,
		10,
		"2025-2026",
		"semester_1",
	)

	assert.Nil(t, results)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByProfessor_LoadsAnswers(t *testing.T) {
	repository, mock, cleanup := setupResultsRepositoryTest(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "evaluations" WHERE professor_id = $1 AND course_id = $2 AND academic_year = $3 AND period = $4`,
	)).
		WithArgs(5, 20, "2025-2026", "semester_2").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"student_id",
				"professor_id",
				"course_id",
				"academic_year",
				"period",
				"submitted_at",
			}).
				AddRow(
					7,
					200,
					5,
					20,
					"2025-2026",
					"semester_2",
					"2026-09-21 09:00:00",
				),
		)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "evaluation_answers" WHERE "evaluation_answers"."evaluation_id" = $1`,
	)).
		WithArgs(7).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"evaluation_id",
				"criterion_id",
				"score",
			}).
				AddRow(10, 7, 1, 5).
				AddRow(11, 7, 2, 4),
		)

	results, err := repository.GetByProfessor(
		5,
		20,
		"2025-2026",
		"semester_2",
	)

	require.NoError(t, err)
	require.Len(t, results, 1)

	assert.Equal(t, 7, results[0].ID)

	require.IsType(t, []evaluations.Evaluation{}, results)

	require.Len(t, results[0].Answers, 2)
	assert.Equal(t, 1, results[0].Answers[0].CriterionID)
	assert.Equal(t, 5, results[0].Answers[0].Score)
	assert.Equal(t, 2, results[0].Answers[1].CriterionID)
	assert.Equal(t, 4, results[0].Answers[1].Score)

	assert.NoError(t, mock.ExpectationsWereMet())
}
