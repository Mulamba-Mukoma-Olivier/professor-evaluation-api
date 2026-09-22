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

func setupServiceTest(t *testing.T) (*Service, sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(
		postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		}),
		&gorm.Config{},
	)
	require.NoError(t, err)

	repository := NewRepository(db)
	service := NewService(repository)

	cleanup := func() {
		sqlDB.Close()
	}

	return service, mock, cleanup
}

func TestService_GetAll_Success(t *testing.T) {
	service, mock, cleanup := setupServiceTest(t)
	defer cleanup()

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
		regexp.QuoteMeta(`SELECT * FROM "courses"`),
	).WillReturnRows(rows)

	courses, err := service.GetAll()

	require.NoError(t, err)
	require.Len(t, courses, 2)

	assert.Equal(t, "INF301", courses[0].Code)
	assert.Equal(t, "INF302", courses[1].Code)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestService_GetByID_InvalidID(t *testing.T) {
	service, _, cleanup := setupServiceTest(t)
	defer cleanup()

	course, err := service.GetByID(0)

	assert.Nil(t, course)
	assert.EqualError(t, err, "invalid course ID")
}

func TestService_GetByID_Success(t *testing.T) {
	service, mock, cleanup := setupServiceTest(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id",
		"code",
		"name",
		"description",
		"department",
		"academic_year",
	}).AddRow(
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

	course, err := service.GetByID(1)

	require.NoError(t, err)
	require.NotNil(t, course)

	assert.Equal(t, 1, course.ID)
	assert.Equal(t, "INF301", course.Code)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestService_Create_Success(t *testing.T) {
	service, mock, cleanup := setupServiceTest(t)
	defer cleanup()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "courses" ("code","name","description","department","academic_year") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`,
		),
	).
		WithArgs(
			"INF303",
			"Algorithmique avancée",
			"Algorithmes avancés",
			"Informatique",
			"2025-2026",
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(3),
		)

	request := CreateCourseRequest{
		Code:         " INF303 ",
		Name:         " Algorithmique avancée ",
		Description:  " Algorithmes avancés ",
		Department:   " Informatique ",
		AcademicYear: " 2025-2026 ",
	}

	course, err := service.Create(request)

	require.NoError(t, err)
	require.NotNil(t, course)

	assert.Equal(t, 3, course.ID)
	assert.Equal(t, "INF303", course.Code)
	assert.Equal(t, "Algorithmique avancée", course.Name)
	assert.Equal(t, "Informatique", course.Department)
	assert.Equal(t, "2025-2026", course.AcademicYear)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestService_Create_MissingCode(t *testing.T) {
	service, _, cleanup := setupServiceTest(t)
	defer cleanup()

	request := CreateCourseRequest{
		Name:         "Algorithmique",
		Department:   "Informatique",
		AcademicYear: "2025-2026",
	}

	course, err := service.Create(request)

	assert.Nil(t, course)
	assert.EqualError(t, err, "course code is required")
}

func TestService_Create_MissingName(t *testing.T) {
	service, _, cleanup := setupServiceTest(t)
	defer cleanup()

	request := CreateCourseRequest{
		Code:         "INF303",
		Department:   "Informatique",
		AcademicYear: "2025-2026",
	}

	course, err := service.Create(request)

	assert.Nil(t, course)
	assert.EqualError(t, err, "course name is required")
}

func TestService_Create_MissingDepartment(t *testing.T) {
	service, _, cleanup := setupServiceTest(t)
	defer cleanup()

	request := CreateCourseRequest{
		Code:         "INF303",
		Name:         "Algorithmique",
		AcademicYear: "2025-2026",
	}

	course, err := service.Create(request)

	assert.Nil(t, course)
	assert.EqualError(t, err, "department is required")
}

func TestService_Create_MissingAcademicYear(t *testing.T) {
	service, _, cleanup := setupServiceTest(t)
	defer cleanup()

	request := CreateCourseRequest{
		Code:       "INF303",
		Name:       "Algorithmique",
		Department: "Informatique",
	}

	course, err := service.Create(request)

	assert.Nil(t, course)
	assert.EqualError(t, err, "academic year is required")
}

func TestService_Update_InvalidID(t *testing.T) {
	service, _, cleanup := setupServiceTest(t)
	defer cleanup()

	request := UpdateCourseRequest{
		Code:         "INF303",
		Name:         "Algorithmique",
		Department:   "Informatique",
		AcademicYear: "2025-2026",
	}

	course, err := service.Update(0, request)

	assert.Nil(t, course)
	assert.EqualError(t, err, "invalid course ID")
}

func TestService_Update_Success(t *testing.T) {
	service, mock, cleanup := setupServiceTest(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id",
		"code",
		"name",
		"description",
		"department",
		"academic_year",
	}).AddRow(
		1,
		"INF301",
		"Ancien nom",
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
		WillReturnRows(rows)

	mock.ExpectExec(
		regexp.QuoteMeta(
			`UPDATE "courses" SET "code"=$1,"name"=$2,"description"=$3,"department"=$4,"academic_year"=$5 WHERE "id" = $6`,
		),
	).
		WithArgs(
			"INF301",
			"Nouveau nom",
			"Nouvelle description",
			"Informatique",
			"2025-2026",
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	request := UpdateCourseRequest{
		Code:         " INF301 ",
		Name:         " Nouveau nom ",
		Description:  " Nouvelle description ",
		Department:   " Informatique ",
		AcademicYear: " 2025-2026 ",
	}

	course, err := service.Update(1, request)

	require.NoError(t, err)
	require.NotNil(t, course)

	assert.Equal(t, "INF301", course.Code)
	assert.Equal(t, "Nouveau nom", course.Name)
	assert.Equal(t, "Nouvelle description", course.Description)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestService_Update_MissingCode(t *testing.T) {
	service, _, cleanup := setupServiceTest(t)
	defer cleanup()

	request := UpdateCourseRequest{
		Name:         "Algorithmique",
		Department:   "Informatique",
		AcademicYear: "2025-2026",
	}

	course, err := service.Update(1, request)

	assert.Nil(t, course)
	assert.EqualError(t, err, "course code is required")
}

func TestService_Update_MissingName(t *testing.T) {
	service, _, cleanup := setupServiceTest(t)
	defer cleanup()

	request := UpdateCourseRequest{
		Code:         "INF303",
		Department:   "Informatique",
		AcademicYear: "2025-2026",
	}

	course, err := service.Update(1, request)

	assert.Nil(t, course)
	assert.EqualError(t, err, "course name is required")
}

func TestService_Update_MissingDepartment(t *testing.T) {
	service, _, cleanup := setupServiceTest(t)
	defer cleanup()

	request := UpdateCourseRequest{
		Code:         "INF303",
		Name:         "Algorithmique",
		AcademicYear: "2025-2026",
	}

	course, err := service.Update(1, request)

	assert.Nil(t, course)
	assert.EqualError(t, err, "department is required")
}

func TestService_Update_MissingAcademicYear(t *testing.T) {
	service, _, cleanup := setupServiceTest(t)
	defer cleanup()

	request := UpdateCourseRequest{
		Code:       "INF303",
		Name:       "Algorithmique",
		Department: "Informatique",
	}

	course, err := service.Update(1, request)

	assert.Nil(t, course)
	assert.EqualError(t, err, "academic year is required")
}

func TestService_Delete_InvalidID(t *testing.T) {
	service, _, cleanup := setupServiceTest(t)
	defer cleanup()

	err := service.Delete(0)

	assert.EqualError(t, err, "invalid course ID")
}

func TestService_Delete_Success(t *testing.T) {
	service, mock, cleanup := setupServiceTest(t)
	defer cleanup()

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(
			`DELETE FROM "courses" WHERE "courses"."id" = $1`,
		),
	).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := service.Delete(1)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
