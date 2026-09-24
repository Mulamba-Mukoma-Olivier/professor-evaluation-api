package evaluations

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupEvaluationRepositoryTest prépare une base GORM
// simulée avec sqlmock.
func setupEvaluationRepositoryTest(
	t *testing.T,
) (*Repository, sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn: sqlDB,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	repository := NewRepository(db)

	cleanup := func() {
		_ = sqlDB.Close()
	}

	return repository, mock, cleanup
}

// =========================================================
// CREATE
// =========================================================

// ---------------------------------------------------------
// Create - succès
// ---------------------------------------------------------

func TestRepository_Create_Success(t *testing.T) {
	repository, mock, cleanup := setupEvaluationRepositoryTest(t)
	defer cleanup()

	submittedAt := time.Now()

	mock.ExpectBegin()

	// Création de l'évaluation.
	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "evaluations" ("student_id","professor_id","course_id","academic_year","period","submitted_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
		),
	).
		WithArgs(
			5,
			10,
			20,
			"2025-2026",
			"S1",
			sqlmock.AnyArg(),
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1),
		)

	// Création de la réponse.
	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "evaluation_answers" ("evaluation_id","criterion_id","score") VALUES ($1,$2,$3) RETURNING "id"`,
		),
	).
		WithArgs(1, 1, 4).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1),
		)

	mock.ExpectCommit()

	// Rechargement de l'évaluation.
	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "evaluations" WHERE "evaluations"."id" = $1 ORDER BY "evaluations"."id" LIMIT $2`,
		),
	).
		WithArgs(1, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"student_id",
				"professor_id",
				"course_id",
				"academic_year",
				"period",
				"submitted_at",
			}).AddRow(
				1,
				5,
				10,
				20,
				"2025-2026",
				"S1",
				submittedAt,
			),
		)

	// Preload des réponses.
	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "evaluation_answers" WHERE "evaluation_answers"."evaluation_id" = $1`,
		),
	).
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"evaluation_id",
				"criterion_id",
				"score",
			}).AddRow(
				1,
				1,
				1,
				4,
			),
		)

	evaluation := Evaluation{
		StudentID:    5,
		ProfessorID:  10,
		CourseID:     20,
		AcademicYear: "2025-2026",
		Period:       "S1",
		SubmittedAt:  submittedAt,
		Answers: []EvaluationAnswer{
			{
				CriterionID: 1,
				Score:       4,
			},
		},
	}

	result, err := repository.Create(evaluation)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, result.ID)
	assert.Equal(t, 5, result.StudentID)
	assert.Equal(t, 10, result.ProfessorID)
	assert.Equal(t, 20, result.CourseID)
	assert.Equal(t, "2025-2026", result.AcademicYear)
	assert.Equal(t, "S1", result.Period)

	require.Len(t, result.Answers, 1)

	assert.Equal(t, 1, result.Answers[0].CriterionID)
	assert.Equal(t, 4, result.Answers[0].Score)

	require.NoError(t, mock.ExpectationsWereMet())
}

// ---------------------------------------------------------
// Create - erreur lors de la création de l'évaluation
// ---------------------------------------------------------

func TestRepository_Create_EvaluationError(t *testing.T) {
	repository, mock, cleanup := setupEvaluationRepositoryTest(t)
	defer cleanup()

	mock.ExpectBegin()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "evaluations" ("student_id","professor_id","course_id","academic_year","period","submitted_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
		),
	).
		WithArgs(
			5,
			10,
			20,
			"2025-2026",
			"S1",
			sqlmock.AnyArg(),
		).
		WillReturnError(assert.AnError)

	mock.ExpectRollback()

	evaluation := Evaluation{
		StudentID:    5,
		ProfessorID:  10,
		CourseID:     20,
		AcademicYear: "2025-2026",
		Period:       "S1",
		Answers: []EvaluationAnswer{
			{
				CriterionID: 1,
				Score:       4,
			},
		},
	}

	result, err := repository.Create(evaluation)

	assert.Nil(t, result)
	assert.Error(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

// ---------------------------------------------------------
// Create - erreur lors de la création d'une réponse
// ---------------------------------------------------------

func TestRepository_Create_AnswerError(t *testing.T) {
	repository, mock, cleanup := setupEvaluationRepositoryTest(t)
	defer cleanup()

	mock.ExpectBegin()

	// L'évaluation est créée avec succès.
	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "evaluations" ("student_id","professor_id","course_id","academic_year","period","submitted_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
		),
	).
		WithArgs(
			5,
			10,
			20,
			"2025-2026",
			"S1",
			sqlmock.AnyArg(),
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1),
		)

	// La réponse échoue.
	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "evaluation_answers" ("evaluation_id","criterion_id","score") VALUES ($1,$2,$3) RETURNING "id"`,
		),
	).
		WithArgs(1, 1, 4).
		WillReturnError(assert.AnError)

	// La transaction doit être annulée.
	mock.ExpectRollback()

	evaluation := Evaluation{
		StudentID:    5,
		ProfessorID:  10,
		CourseID:     20,
		AcademicYear: "2025-2026",
		Period:       "S1",
		Answers: []EvaluationAnswer{
			{
				CriterionID: 1,
				Score:       4,
			},
		},
	}

	result, err := repository.Create(evaluation)

	assert.Nil(t, result)
	assert.Error(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

// =========================================================
// GET ALL
// =========================================================

// ---------------------------------------------------------
// GetAll - succès
// ---------------------------------------------------------

func TestRepository_GetAll_Success(t *testing.T) {
	repository, mock, cleanup := setupEvaluationRepositoryTest(t)
	defer cleanup()

	submittedAt := time.Now()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "evaluations"`,
		),
	).
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
					5,
					10,
					20,
					"2025-2026",
					"S1",
					submittedAt,
				).
				AddRow(
					2,
					6,
					11,
					21,
					"2025-2026",
					"S2",
					submittedAt,
				),
		)

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "evaluation_answers" WHERE "evaluation_answers"."evaluation_id" IN ($1,$2)`,
		),
	).
		WithArgs(1, 2).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"evaluation_id",
				"criterion_id",
				"score",
			}).
				AddRow(1, 1, 1, 4).
				AddRow(2, 2, 2, 5),
		)

	result, err := repository.GetAll()

	require.NoError(t, err)
	require.Len(t, result, 2)

	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, 2, result[1].ID)

	assert.Len(t, result[0].Answers, 1)
	assert.Len(t, result[1].Answers, 1)

	require.NoError(t, mock.ExpectationsWereMet())
}

// ---------------------------------------------------------
// GetAll - erreur DB
// ---------------------------------------------------------

func TestRepository_GetAll_Error(t *testing.T) {
	repository, mock, cleanup := setupEvaluationRepositoryTest(t)
	defer cleanup()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "evaluations"`,
		),
	).
		WillReturnError(assert.AnError)

	result, err := repository.GetAll()

	assert.Nil(t, result)
	assert.Error(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

// =========================================================
// GET BY ID
// =========================================================

// ---------------------------------------------------------
// GetByID - succès
// ---------------------------------------------------------

func TestRepository_GetByID_Success(t *testing.T) {
	repository, mock, cleanup := setupEvaluationRepositoryTest(t)
	defer cleanup()

	submittedAt := time.Now()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "evaluations" WHERE "evaluations"."id" = $1 ORDER BY "evaluations"."id" LIMIT $2`,
		),
	).
		WithArgs(1, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"student_id",
				"professor_id",
				"course_id",
				"academic_year",
				"period",
				"submitted_at",
			}).AddRow(
				1,
				5,
				10,
				20,
				"2025-2026",
				"S1",
				submittedAt,
			),
		)

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "evaluation_answers" WHERE "evaluation_answers"."evaluation_id" = $1`,
		),
	).
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"evaluation_id",
				"criterion_id",
				"score",
			}).AddRow(
				1,
				1,
				1,
				4,
			),
		)

	result, err := repository.GetByID(1)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, result.ID)
	assert.Equal(t, 5, result.StudentID)
	assert.Equal(t, 10, result.ProfessorID)
	assert.Equal(t, 20, result.CourseID)

	require.Len(t, result.Answers, 1)

	require.NoError(t, mock.ExpectationsWereMet())
}

// ---------------------------------------------------------
// GetByID - introuvable
// ---------------------------------------------------------

func TestRepository_GetByID_NotFound(t *testing.T) {
	repository, mock, cleanup := setupEvaluationRepositoryTest(t)
	defer cleanup()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "evaluations" WHERE "evaluations"."id" = $1 ORDER BY "evaluations"."id" LIMIT $2`,
		),
	).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repository.GetByID(999)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrEvaluationNotFound)

	require.NoError(t, mock.ExpectationsWereMet())
}

// ---------------------------------------------------------
// GetByID - erreur DB
// ---------------------------------------------------------

func TestRepository_GetByID_Error(t *testing.T) {
	repository, mock, cleanup := setupEvaluationRepositoryTest(t)
	defer cleanup()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "evaluations" WHERE "evaluations"."id" = $1 ORDER BY "evaluations"."id" LIMIT $2`,
		),
	).
		WithArgs(1, 1).
		WillReturnError(assert.AnError)

	result, err := repository.GetByID(1)

	assert.Nil(t, result)
	assert.Error(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

// =========================================================
// EXISTS
// =========================================================

// ---------------------------------------------------------
// Exists - vrai
// ---------------------------------------------------------

func TestRepository_Exists_True(t *testing.T) {
	repository, mock, cleanup := setupEvaluationRepositoryTest(t)
	defer cleanup()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT count(*) FROM "evaluations" WHERE student_id = $1 AND professor_id = $2 AND course_id = $3 AND academic_year = $4 AND period = $5`,
		),
	).
		WithArgs(
			5,
			10,
			20,
			"2025-2026",
			"S1",
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(1),
		)

	result, err := repository.Exists(
		5,
		10,
		20,
		"2025-2026",
		"S1",
	)

	require.NoError(t, err)
	assert.True(t, result)

	require.NoError(t, mock.ExpectationsWereMet())
}

// ---------------------------------------------------------
// Exists - faux
// ---------------------------------------------------------

func TestRepository_Exists_False(t *testing.T) {
	repository, mock, cleanup := setupEvaluationRepositoryTest(t)
	defer cleanup()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT count(*) FROM "evaluations" WHERE student_id = $1 AND professor_id = $2 AND course_id = $3 AND academic_year = $4 AND period = $5`,
		),
	).
		WithArgs(
			5,
			10,
			20,
			"2025-2026",
			"S1",
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(0),
		)

	result, err := repository.Exists(
		5,
		10,
		20,
		"2025-2026",
		"S1",
	)

	require.NoError(t, err)
	assert.False(t, result)

	require.NoError(t, mock.ExpectationsWereMet())
}

// ---------------------------------------------------------
// Exists - erreur DB
// ---------------------------------------------------------

func TestRepository_Exists_Error(t *testing.T) {
	repository, mock, cleanup := setupEvaluationRepositoryTest(t)
	defer cleanup()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT count(*) FROM "evaluations" WHERE student_id = $1 AND professor_id = $2 AND course_id = $3 AND academic_year = $4 AND period = $5`,
		),
	).
		WithArgs(
			5,
			10,
			20,
			"2025-2026",
			"S1",
		).
		WillReturnError(assert.AnError)

	result, err := repository.Exists(
		5,
		10,
		20,
		"2025-2026",
		"S1",
	)

	assert.False(t, result)
	assert.Error(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

// =========================================================
// DELETE
// =========================================================

// ---------------------------------------------------------
// Delete - succès
// ---------------------------------------------------------

func TestRepository_Delete_Success(t *testing.T) {
	repository, mock, cleanup := setupEvaluationRepositoryTest(t)
	defer cleanup()

	mock.ExpectExec(
		regexp.QuoteMeta(
			`DELETE FROM "evaluations" WHERE "evaluations"."id" = $1`,
		),
	).
		WithArgs(1).
		WillReturnResult(
			sqlmock.NewResult(0, 1),
		)

	err := repository.Delete(1)

	assert.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

// ---------------------------------------------------------
// Delete - introuvable
// ---------------------------------------------------------

func TestRepository_Delete_NotFound(t *testing.T) {
	repository, mock, cleanup := setupEvaluationRepositoryTest(t)
	defer cleanup()

	mock.ExpectExec(
		regexp.QuoteMeta(
			`DELETE FROM "evaluations" WHERE "evaluations"."id" = $1`,
		),
	).
		WithArgs(999).
		WillReturnResult(
			sqlmock.NewResult(0, 0),
		)

	err := repository.Delete(999)

	assert.ErrorIs(t, err, ErrEvaluationNotFound)

	require.NoError(t, mock.ExpectationsWereMet())
}

// ---------------------------------------------------------
// Delete - erreur DB
// ---------------------------------------------------------

func TestRepository_Delete_Error(t *testing.T) {
	repository, mock, cleanup := setupEvaluationRepositoryTest(t)
	defer cleanup()

	mock.ExpectExec(
		regexp.QuoteMeta(
			`DELETE FROM "evaluations" WHERE "evaluations"."id" = $1`,
		),
	).
		WithArgs(1).
		WillReturnError(assert.AnError)

	err := repository.Delete(1)

	assert.Error(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}
