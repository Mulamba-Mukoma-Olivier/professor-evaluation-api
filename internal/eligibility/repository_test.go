package eligibility

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupEligibilityRepositoryTest(
	t *testing.T,
) (*Repository, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	t.Cleanup(func() {
		sqlDB.Close()
	})

	db, err := gorm.Open(
		postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		}),
		&gorm.Config{},
	)

	if err != nil {
		t.Fatalf("failed to open gorm database: %v", err)
	}

	return NewRepository(db), mock
}

func TestRepository_GetByStudentID_Success(t *testing.T) {
	repository, mock := setupEligibilityRepositoryTest(t)

	rows := sqlmock.NewRows([]string{
		"id",
		"student_id",
		"enrollment",
		"academic_fees",
		"laboratory_fees",
		"access_fees",
		"eligible",
	}).
		AddRow(
			1,
			10,
			true,
			true,
			true,
			true,
			true,
		)

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "eligibilities" WHERE student_id = $1 ORDER BY "eligibilities"."id" LIMIT $2`,
		),
	).
		WithArgs(10, 1).
		WillReturnRows(rows)

	result, err := repository.GetByStudentID(10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, 10, result.StudentID)
	assert.True(t, result.Enrollment)
	assert.True(t, result.AcademicFees)
	assert.True(t, result.LaboratoryFees)
	assert.True(t, result.AccessFees)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByStudentID_NotFound(t *testing.T) {
	repository, mock := setupEligibilityRepositoryTest(t)

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "eligibilities" WHERE student_id = $1 ORDER BY "eligibilities"."id" LIMIT $2`,
		),
	).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repository.GetByStudentID(999)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrEligibilityNotFound)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByStudentID_DatabaseError(t *testing.T) {
	repository, mock := setupEligibilityRepositoryTest(t)

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "eligibilities" WHERE student_id = $1 ORDER BY "eligibilities"."id" LIMIT $2`,
		),
	).
		WithArgs(10, 1).
		WillReturnError(
			assert.AnError,
		)

	result, err := repository.GetByStudentID(10)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, assert.AnError)

	assert.NoError(t, mock.ExpectationsWereMet())
}
