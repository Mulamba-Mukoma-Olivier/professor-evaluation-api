package courses

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return db, mock
}

func closeTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	sqlDB, err := db.DB()
	require.NoError(t, err)

	require.NoError(t, sqlDB.Close())
}

func TestRepository_GetAll_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer closeTestDB(t, db)

	repository := NewRepository(db)

	rows := sqlmock.NewRows([]string{
		"id",
		"code",
		"name",
		"description",
		"department",
		"academic_year",
	}).
		AddRow(
			1,
			"INF301",
			"Programmation parallèle",
			"Introduction à la programmation parallèle",
			"Informatique",
			"2025-2026",
		).
		AddRow(
			2,
			"INF302",
			"Bases de données",
			"Introduction aux bases de données",
			"Informatique",
			"2025-2026",
		)

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "courses"`,
		),
	).
		WillReturnRows(rows)

	courses, err := repository.GetAll()

	require.NoError(t, err)
	require.Len(t, courses, 2)

	assert.Equal(t, 1, courses[0].ID)
	assert.Equal(t, "INF301", courses[0].Code)
	assert.Equal(t, "Programmation parallèle", courses[0].Name)

	assert.Equal(t, 2, courses[1].ID)
	assert.Equal(t, "INF302", courses[1].Code)
	assert.Equal(t, "Bases de données", courses[1].Name)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAll_DatabaseError(t *testing.T) {
	db, mock := setupTestDB(t)
	defer closeTestDB(t, db)

	repository := NewRepository(db)

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "courses"`,
		),
	).
		WillReturnError(assert.AnError)

	courses, err := repository.GetAll()

	assert.Error(t, err)
	assert.Nil(t, courses)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer closeTestDB(t, db)

	repository := NewRepository(db)

	rows := sqlmock.NewRows([]string{
		"id",
		"code",
		"name",
		"description",
		"department",
		"academic_year",
	}).
		AddRow(
			1,
			"INF301",
			"Programmation parallèle",
			"Introduction à la programmation parallèle",
			"Informatique",
			"2025-2026",
		)

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "courses" WHERE "courses"."id" = $1 ORDER BY "courses"."id" LIMIT $2`,
		),
	).
		WithArgs(1, 1).
		WillReturnRows(rows)

	course, err := repository.GetByID(1)

	require.NoError(t, err)
	require.NotNil(t, course)

	assert.Equal(t, 1, course.ID)
	assert.Equal(t, "INF301", course.Code)
	assert.Equal(t, "Programmation parallèle", course.Name)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	defer closeTestDB(t, db)

	repository := NewRepository(db)

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "courses" WHERE "courses"."id" = $1 ORDER BY "courses"."id" LIMIT $2`,
		),
	).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	course, err := repository.GetByID(999)

	assert.Nil(t, course)
	assert.ErrorIs(t, err, ErrCourseNotFound)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_DatabaseError(t *testing.T) {
	db, mock := setupTestDB(t)
	defer closeTestDB(t, db)

	repository := NewRepository(db)

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "courses" WHERE "courses"."id" = $1 ORDER BY "courses"."id" LIMIT $2`,
		),
	).
		WithArgs(1, 1).
		WillReturnError(assert.AnError)

	course, err := repository.GetByID(1)

	assert.Nil(t, course)
	assert.Error(t, err)
	assert.NotErrorIs(t, err, ErrCourseNotFound)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Create_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer closeTestDB(t, db)

	repository := NewRepository(db)

	course := Course{
		Code:         "INF303",
		Name:         "Algorithmique avancée",
		Description:  "Cours avancé d'algorithmique",
		Department:   "Informatique",
		AcademicYear: "2025-2026",
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "courses" ("code","name","description","department","academic_year") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`,
		),
	).
		WithArgs(
			course.Code,
			course.Name,
			course.Description,
			course.Department,
			course.AcademicYear,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(3),
		)

	createdCourse, err := repository.Create(course)

	require.NoError(t, err)
	require.NotNil(t, createdCourse)

	assert.Equal(t, 3, createdCourse.ID)
	assert.Equal(t, "INF303", createdCourse.Code)
	assert.Equal(t, "Algorithmique avancée", createdCourse.Name)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Create_DatabaseError(t *testing.T) {
	db, mock := setupTestDB(t)
	defer closeTestDB(t, db)

	repository := NewRepository(db)

	course := Course{
		Code:         "INF303",
		Name:         "Algorithmique avancée",
		Description:  "Cours avancé d'algorithmique",
		Department:   "Informatique",
		AcademicYear: "2025-2026",
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "courses" ("code","name","description","department","academic_year") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`,
		),
	).
		WithArgs(
			course.Code,
			course.Name,
			course.Description,
			course.Department,
			course.AcademicYear,
		).
		WillReturnError(assert.AnError)

	createdCourse, err := repository.Create(course)

	assert.Nil(t, createdCourse)
	assert.Error(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer closeTestDB(t, db)

	repository := NewRepository(db)

	existingRows := sqlmock.NewRows([]string{
		"id",
		"code",
		"name",
		"description",
		"department",
		"academic_year",
	}).
		AddRow(
			1,
			"INF301",
			"Ancien cours",
			"Ancienne description",
			"Informatique",
			"2024-2025",
		)

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "courses" WHERE "courses"."id" = $1 ORDER BY "courses"."id" LIMIT $2`,
		),
	).
		WithArgs(1, 1).
		WillReturnRows(existingRows)

	mock.ExpectExec(
		regexp.QuoteMeta(
			`UPDATE "courses" SET "code"=$1,"name"=$2,"description"=$3,"department"=$4,"academic_year"=$5 WHERE "id" = $6`,
		),
	).
		WithArgs(
			"INF301",
			"Nouveau cours",
			"Nouvelle description",
			"Informatique",
			"2025-2026",
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	course := Course{
		Code:         "INF301",
		Name:         "Nouveau cours",
		Description:  "Nouvelle description",
		Department:   "Informatique",
		AcademicYear: "2025-2026",
	}

	updatedCourse, err := repository.Update(1, course)

	require.NoError(t, err)
	require.NotNil(t, updatedCourse)

	assert.Equal(t, 1, updatedCourse.ID)
	assert.Equal(t, "INF301", updatedCourse.Code)
	assert.Equal(t, "Nouveau cours", updatedCourse.Name)
	assert.Equal(t, "2025-2026", updatedCourse.AcademicYear)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	defer closeTestDB(t, db)

	repository := NewRepository(db)

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "courses" WHERE "courses"."id" = $1 ORDER BY "courses"."id" LIMIT $2`,
		),
	).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	course := Course{
		Code:         "INF301",
		Name:         "Nouveau cours",
		Description:  "Description",
		Department:   "Informatique",
		AcademicYear: "2025-2026",
	}

	updatedCourse, err := repository.Update(999, course)

	assert.Nil(t, updatedCourse)
	assert.ErrorIs(t, err, ErrCourseNotFound)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer closeTestDB(t, db)

	repository := NewRepository(db)

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(
			`DELETE FROM "courses" WHERE "courses"."id" = $1`,
		),
	).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := repository.Delete(1)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	defer closeTestDB(t, db)

	repository := NewRepository(db)

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(
			`DELETE FROM "courses" WHERE "courses"."id" = $1`,
		),
	).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectCommit()

	err := repository.Delete(999)

	assert.ErrorIs(t, err, ErrCourseNotFound)

	require.NoError(t, mock.ExpectationsWereMet())
}
