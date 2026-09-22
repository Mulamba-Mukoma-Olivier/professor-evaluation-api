package criteria

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// createTestDB crée une base GORM utilisant sqlmock.
// Aucun vrai serveur PostgreSQL n'est utilisé pendant les tests.
func createTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

// ---------------------------------------------------------
// GetAll
// ---------------------------------------------------------

func TestRepository_GetAll_Success(t *testing.T) {
	db, mock := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"description",
		"max_score",
		"active",
	}).
		AddRow(1, "Clarté", "Clarté des explications", 5, true).
		AddRow(2, "Ponctualité", "Respect des horaires", 5, true)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "criteria"`,
	)).
		WillReturnRows(rows)

	repository := NewRepository(db)

	result, err := repository.GetAll()

	require.NoError(t, err)
	require.Len(t, result, 2)

	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, "Clarté", result[0].Name)
	assert.Equal(t, "Clarté des explications", result[0].Description)
	assert.Equal(t, 5, result[0].MaxScore)
	assert.True(t, result[0].Active)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAll_Error(t *testing.T) {
	db, mock := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "criteria"`,
	)).
		WillReturnError(assert.AnError)

	repository := NewRepository(db)

	result, err := repository.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---------------------------------------------------------
// GetActive
// ---------------------------------------------------------

func TestRepository_GetActive_Success(t *testing.T) {
	db, mock := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"description",
		"max_score",
		"active",
	}).
		AddRow(1, "Clarté", "Clarté des explications", 5, true)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "criteria" WHERE active = $1`,
	)).
		WithArgs(true).
		WillReturnRows(rows)

	repository := NewRepository(db)

	result, err := repository.GetActive()

	require.NoError(t, err)
	require.Len(t, result, 1)

	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, "Clarté", result[0].Name)
	assert.True(t, result[0].Active)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetActive_Error(t *testing.T) {
	db, mock := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "criteria" WHERE active = $1`,
	)).
		WithArgs(true).
		WillReturnError(assert.AnError)

	repository := NewRepository(db)

	result, err := repository.GetActive()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---------------------------------------------------------
// GetByID
// ---------------------------------------------------------

func TestRepository_GetByID_Success(t *testing.T) {
	db, mock := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"description",
		"max_score",
		"active",
	}).
		AddRow(1, "Clarté", "Clarté des explications", 5, true)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "criteria" WHERE "criteria"."id" = $1 ORDER BY "criteria"."id" LIMIT $2`,
	)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	repository := NewRepository(db)

	result, err := repository.GetByID(1)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Clarté", result.Name)
	assert.Equal(t, 5, result.MaxScore)
	assert.True(t, result.Active)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	db, mock := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "criteria" WHERE "criteria"."id" = $1 ORDER BY "criteria"."id" LIMIT $2`,
	)).
		WithArgs(99, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	repository := NewRepository(db)

	result, err := repository.GetByID(99)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrCriterionNotFound)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_Error(t *testing.T) {
	db, mock := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "criteria" WHERE "criteria"."id" = $1 ORDER BY "criteria"."id" LIMIT $2`,
	)).
		WithArgs(1, 1).
		WillReturnError(assert.AnError)

	repository := NewRepository(db)

	result, err := repository.GetByID(1)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---------------------------------------------------------
// Create
// ---------------------------------------------------------

func TestRepository_Create_Success(t *testing.T) {
	db, mock := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	criterion := Criterion{
		Name:        "Clarté",
		Description: "Clarté des explications",
		MaxScore:    5,
		Active:      true,
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "criteria" ("name","description","max_score","active") VALUES ($1,$2,$3,$4) RETURNING "id"`,
	)).
		WithArgs(
			criterion.Name,
			criterion.Description,
			criterion.MaxScore,
			criterion.Active,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1),
		)

	repository := NewRepository(db)

	result, err := repository.Create(criterion)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Clarté", result.Name)
	assert.Equal(t, "Clarté des explications", result.Description)
	assert.Equal(t, 5, result.MaxScore)
	assert.True(t, result.Active)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Create_Error(t *testing.T) {
	db, mock := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	criterion := Criterion{
		Name:        "Clarté",
		Description: "Clarté des explications",
		MaxScore:    5,
		Active:      true,
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "criteria" ("name","description","max_score","active") VALUES ($1,$2,$3,$4) RETURNING "id"`,
	)).
		WithArgs(
			criterion.Name,
			criterion.Description,
			criterion.MaxScore,
			criterion.Active,
		).
		WillReturnError(assert.AnError)

	repository := NewRepository(db)

	result, err := repository.Create(criterion)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---------------------------------------------------------
// Update
// ---------------------------------------------------------

func TestRepository_Update_Success(t *testing.T) {
	db, mock := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 1. GORM recherche d'abord le critère existant.
	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"description",
		"max_score",
		"active",
	}).
		AddRow(1, "Ancien nom", "Ancienne description", 5, true)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "criteria" WHERE "criteria"."id" = $1 ORDER BY "criteria"."id" LIMIT $2`,
	)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	// 2. GORM sauvegarde les nouvelles valeurs.
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "criteria" SET "name"=$1,"description"=$2,"max_score"=$3,"active"=$4 WHERE "id" = $5`,
	)).
		WithArgs(
			"Nouveau nom",
			"Nouvelle description",
			10,
			false,
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repository := NewRepository(db)

	criterion := Criterion{
		Name:        "Nouveau nom",
		Description: "Nouvelle description",
		MaxScore:    10,
		Active:      false,
	}

	result, err := repository.Update(1, criterion)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Nouveau nom", result.Name)
	assert.Equal(t, "Nouvelle description", result.Description)
	assert.Equal(t, 10, result.MaxScore)
	assert.False(t, result.Active)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update_NotFound(t *testing.T) {
	db, mock := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "criteria" WHERE "criteria"."id" = $1 ORDER BY "criteria"."id" LIMIT $2`,
	)).
		WithArgs(99, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	repository := NewRepository(db)

	result, err := repository.Update(99, Criterion{
		Name:     "Clarté",
		MaxScore: 5,
		Active:   true,
	})

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrCriterionNotFound)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---------------------------------------------------------
// Delete
// ---------------------------------------------------------

func TestRepository_Delete_Success(t *testing.T) {
	db, mock := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "criteria" WHERE "criteria"."id" = $1`,
	)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repository := NewRepository(db)

	err := repository.Delete(1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete_NotFound(t *testing.T) {
	db, mock := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "criteria" WHERE "criteria"."id" = $1`,
	)).
		WithArgs(99).
		WillReturnResult(sqlmock.NewResult(99, 0))

	repository := NewRepository(db)

	err := repository.Delete(99)

	assert.ErrorIs(t, err, ErrCriterionNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete_Error(t *testing.T) {
	db, mock := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "criteria" WHERE "criteria"."id" = $1`,
	)).
		WithArgs(1).
		WillReturnError(assert.AnError)

	repository := NewRepository(db)

	err := repository.Delete(1)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
