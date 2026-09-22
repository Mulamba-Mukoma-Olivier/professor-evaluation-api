package professors

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *Repository) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	repository := NewRepository(db)

	return db, mock, repository
}

func TestNewRepository(t *testing.T) {
	db, _, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	assert.NotNil(t, repository)
	assert.Equal(t, db, repository.db)
}

func TestRepository_GetAll_Success(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"matricule",
		"first_name",
		"last_name",
		"email",
		"department",
		"grade",
		"active",
		"status",
	}).
		AddRow(1, "PROF001", "Jean", "Dupont", "jean@test.com", "Informatique", "Professeur", true, "active").
		AddRow(2, "PROF002", "Paul", "Martin", "paul@test.com", "Mathématiques", "Assistant", true, "active")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "professors"`,
	)).
		WillReturnRows(rows)

	result, err := repository.GetAll()

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "PROF001", result[0].Matricule)
	assert.Equal(t, "Jean", result[0].FirstName)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAll_Error(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "professors"`,
	)).
		WillReturnError(assert.AnError)

	result, err := repository.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAllPaginated_Success(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	// page = 2, pageSize = 10
	// offset = (2 - 1) * 10 = 10

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "professors"`,
	)).
		WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(25),
		)

	rows := sqlmock.NewRows([]string{
		"id",
		"matricule",
		"first_name",
		"last_name",
		"email",
		"department",
		"grade",
		"active",
		"status",
	}).
		AddRow(11, "PROF011", "Jean", "Dupont", "jean@test.com", "Informatique", "Professeur", true, "active").
		AddRow(12, "PROF012", "Paul", "Martin", "paul@test.com", "Mathématiques", "Assistant", true, "active")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "professors" LIMIT $1 OFFSET $2`,
	)).
		WithArgs(10, 10).
		WillReturnRows(rows)

	result, total, err := repository.GetAllPaginated(2, 10)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 25, total)
	assert.Equal(t, "PROF011", result[0].Matricule)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAllPaginated_CountError(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "professors"`,
	)).
		WillReturnError(assert.AnError)

	result, total, err := repository.GetAllPaginated(1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, total)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAllPaginated_FindError(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "professors"`,
	)).
		WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(25),
		)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "professors" LIMIT $1 OFFSET $2`,
	)).
		WithArgs(10, 0).
		WillReturnError(assert.AnError)

	result, total, err := repository.GetAllPaginated(1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, total)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetActive_Success(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"matricule",
		"first_name",
		"last_name",
		"email",
		"department",
		"grade",
		"active",
		"status",
	}).
		AddRow(1, "PROF001", "Jean", "Dupont", "jean@test.com", "Informatique", "Professeur", true, "active")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "professors" WHERE active = $1`,
	)).
		WithArgs(true).
		WillReturnRows(rows)

	result, err := repository.GetActive()

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.True(t, result[0].Active)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetActive_Error(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "professors" WHERE active = $1`,
	)).
		WithArgs(true).
		WillReturnError(assert.AnError)

	result, err := repository.GetActive()

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByStatus_Success(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"matricule",
		"first_name",
		"last_name",
		"email",
		"department",
		"grade",
		"active",
		"status",
	}).
		AddRow(1, "PROF001", "Jean", "Dupont", "jean@test.com", "Informatique", "Professeur", false, "on_leave")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "professors" WHERE status = $1`,
	)).
		WithArgs("on_leave").
		WillReturnRows(rows)

	result, err := repository.GetByStatus("on_leave")

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "on_leave", result[0].Status)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByStatus_Error(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "professors" WHERE status = $1`,
	)).
		WithArgs("active").
		WillReturnError(assert.AnError)

	result, err := repository.GetByStatus("active")

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_Success(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"matricule",
		"first_name",
		"last_name",
		"email",
		"department",
		"grade",
		"active",
		"status",
	}).
		AddRow(1, "PROF001", "Jean", "Dupont", "jean@test.com", "Informatique", "Professeur", true, "active")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "professors" WHERE "professors"."id" = $1 ORDER BY "professors"."id" LIMIT $2`,
	)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	result, err := repository.GetByID(1)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "PROF001", result.Matricule)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "professors" WHERE "professors"."id" = $1 ORDER BY "professors"."id" LIMIT $2`,
	)).
		WithArgs(999, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"matricule",
			"first_name",
			"last_name",
			"email",
			"department",
			"grade",
			"active",
			"status",
		}))

	result, err := repository.GetByID(999)

	assert.Nil(t, result)
	assert.EqualError(t, err, "professor not found")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_Error(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "professors" WHERE "professors"."id" = $1 ORDER BY "professors"."id" LIMIT $2`,
	)).
		WithArgs(1, 1).
		WillReturnError(assert.AnError)

	result, err := repository.GetByID(1)

	assert.Nil(t, result)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Create_Success(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	professor := Professor{
		Matricule:  "PROF001",
		FirstName:  "Jean",
		LastName:   "Dupont",
		Email:      "jean@test.com",
		Department: "Informatique",
		Grade:      "Professeur",
		Active:     true,
		Status:     "active",
	}

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "professors"`,
	)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1),
		)

	mock.ExpectCommit()

	result, err := repository.Create(professor)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "PROF001", result.Matricule)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Create_Error(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	professor := Professor{
		Matricule:  "PROF001",
		FirstName:  "Jean",
		LastName:   "Dupont",
		Department: "Informatique",
	}

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "professors"`,
	)).
		WillReturnError(assert.AnError)

	mock.ExpectRollback()

	result, err := repository.Create(professor)

	assert.Nil(t, result)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update_Success(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "professors" WHERE "professors"."id" = $1 ORDER BY "professors"."id" LIMIT $2`,
	)).
		WithArgs(1, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"matricule",
				"first_name",
				"last_name",
				"email",
				"department",
				"grade",
				"active",
				"status",
			}).
				AddRow(1, "PROF001", "Jean", "Dupont", "old@test.com", "Informatique", "Assistant", true, "active"),
		)

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "professors"`,
	)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	professor := Professor{
		Matricule:  "PROF001",
		FirstName:  "Jean",
		LastName:   "Dupont",
		Email:      "new@test.com",
		Department: "Informatique",
		Grade:      "Professeur",
		Active:     true,
		Status:     "active",
	}

	result, err := repository.Update(1, professor)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new@test.com", result.Email)
	assert.Equal(t, "Professeur", result.Grade)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update_NotFound(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "professors" WHERE "professors"."id" = $1 ORDER BY "professors"."id" LIMIT $2`,
	)).
		WithArgs(999, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"matricule",
			"first_name",
			"last_name",
			"email",
			"department",
			"grade",
			"active",
			"status",
		}))

	professor := Professor{
		Matricule:  "PROF999",
		FirstName:  "Jean",
		LastName:   "Dupont",
		Department: "Informatique",
	}

	result, err := repository.Update(999, professor)

	assert.Nil(t, result)
	assert.EqualError(t, err, "professor not found")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete_Success(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "professors"`,
	)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = repository.Delete(1)

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete_NotFound(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "professors"`,
	)).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	err = repository.Delete(999)

	assert.EqualError(t, err, "professor not found")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete_Error(t *testing.T) {
	db, mock, repository := setupTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "professors"`,
	)).
		WithArgs(1).
		WillReturnError(assert.AnError)

	mock.ExpectRollback()

	err = repository.Delete(1)

	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
