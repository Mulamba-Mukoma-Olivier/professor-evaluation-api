package auth

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *Repository) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	repo := NewRepository(gormDB)

	return gormDB, mock, repo
}

func closeTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	sqlDB, err := db.DB()
	require.NoError(t, err)

	require.NoError(t, sqlDB.Close())
}

func TestNewRepository(t *testing.T) {
	db, _, repo := setupTestDB(t)
	defer closeTestDB(t, db)

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.db)
}

func TestRepository_FindByEmail_Success(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer closeTestDB(t, db)

	rows := sqlmock.NewRows([]string{
		"id",
		"matricule",
		"name",
		"email",
		"password_hash",
		"role",
	}).AddRow(
		1,
		"12345",
		"Test User",
		"test@example.com",
		"hashed_password",
		"STUDENT",
	)

	mock.ExpectQuery(
		`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"\."id" LIMIT \$2`,
	).
		WithArgs("test@example.com", 1).
		WillReturnRows(rows)

	user, err := repo.FindByEmail("test@example.com")

	require.NoError(t, err)
	require.NotNil(t, user)

	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "12345", user.Matricule)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashed_password", user.PasswordHash)
	assert.Equal(t, "STUDENT", user.Role)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_FindByEmail_NotFound(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer closeTestDB(t, db)

	mock.ExpectQuery(
		`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"\."id" LIMIT \$2`,
	).
		WithArgs("nonexistent@example.com", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.FindByEmail("nonexistent@example.com")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "user not found")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_FindByMatricule_Success(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer closeTestDB(t, db)

	rows := sqlmock.NewRows([]string{
		"id",
		"matricule",
		"name",
		"email",
		"password_hash",
		"role",
	}).AddRow(
		1,
		"12345",
		"Test User",
		"test@example.com",
		"hashed_password",
		"STUDENT",
	)

	mock.ExpectQuery(
		`SELECT \* FROM "users" WHERE matricule = \$1 ORDER BY "users"\."id" LIMIT \$2`,
	).
		WithArgs("12345", 1).
		WillReturnRows(rows)

	user, err := repo.FindByMatricule("12345")

	require.NoError(t, err)
	require.NotNil(t, user)

	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "12345", user.Matricule)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "STUDENT", user.Role)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Create_Success(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer closeTestDB(t, db)

	user := User{
		Matricule:    "12345",
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Role:         "STUDENT",
	}

	mock.ExpectBegin()

	mock.ExpectQuery(
		`INSERT INTO "users"`,
	).
		WithArgs(
			user.Matricule,
			user.Name,
			user.Email,
			user.PasswordHash,
			user.Role,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1),
		)

	mock.ExpectCommit()

	createdUser, err := repo.Create(user)

	require.NoError(t, err)
	require.NotNil(t, createdUser)

	assert.Equal(t, 1, createdUser.ID)
	assert.Equal(t, user.Matricule, createdUser.Matricule)
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.PasswordHash, createdUser.PasswordHash)
	assert.Equal(t, user.Role, createdUser.Role)

	require.NoError(t, mock.ExpectationsWereMet())
}
