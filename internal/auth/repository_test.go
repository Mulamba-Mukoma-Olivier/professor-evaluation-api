package auth

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *Repository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := NewRepository(gormDB)
	return gormDB, mock, repo
}

func TestNewRepository(t *testing.T) {
	db, _, _ := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	repo := NewRepository(db)
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.db)
}

func TestRepository_FindByEmail_Success(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"id", "matricule", "name", "email", "password_hash", "role"}).
		AddRow(1, "12345", "Test User", "test@example.com", "hashed_password", "STUDENT")

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs("test@example.com", 1).
		WillReturnRows(rows)

	user, err := repo.FindByEmail("test@example.com")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_FindByEmail_NotFound(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs("nonexistent@example.com", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.FindByEmail("nonexistent@example.com")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_FindByMatricule_Success(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"id", "matricule", "name", "email", "password_hash", "role"}).
		AddRow(1, "12345", "Test User", "test@example.com", "hashed_password", "STUDENT")

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE matricule = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs("12345", 1).
		WillReturnRows(rows)

	user, err := repo.FindByMatricule("12345")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "12345", user.Matricule)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Create_Success(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	user := User{
		Matricule:    "12345",
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Role:         "STUDENT",
	}

	// GORM uses Query with RETURNING instead of Exec
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs(user.Matricule, user.Name, user.Email, user.PasswordHash, user.Role).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	createdUser, err := repo.Create(user)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}
