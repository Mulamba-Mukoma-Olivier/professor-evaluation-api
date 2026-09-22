package professors

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

func TestRepository_GetAll_Success(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"id", "matricule", "first_name", "last_name", "email", "department", "grade", "active", "status"}).
		AddRow(1, "12345", "John", "Doe", "john@example.com", "CS", "Professor", true, "active").
		AddRow(2, "12346", "Jane", "Smith", "jane@example.com", "Math", "Associate", true, "active")

	mock.ExpectQuery(`SELECT \* FROM "professors"`).WillReturnRows(rows)

	professors, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, professors, 2)
	assert.Equal(t, "John", professors[0].FirstName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_Success(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"id", "matricule", "first_name", "last_name", "email", "department", "grade", "active", "status"}).
		AddRow(1, "12345", "John", "Doe", "john@example.com", "CS", "Professor", true, "active")

	mock.ExpectQuery(`SELECT \* FROM "professors" WHERE "professors"\."id" = \$1 ORDER BY "professors"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	professor, err := repo.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, professor)
	assert.Equal(t, "John", professor.FirstName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	mock.ExpectQuery(`SELECT \* FROM "professors" WHERE "professors"\."id" = \$1 ORDER BY "professors"\."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	professor, err := repo.GetByID(999)

	assert.Error(t, err)
	assert.Nil(t, professor)
	assert.Equal(t, "professor not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Create_Success(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	professor := Professor{
		Matricule:  "12345",
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john@example.com",
		Department: "CS",
		Grade:      "Professor",
		Active:     true,
		Status:     "active",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "professors"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	createdProfessor, err := repo.Create(professor)

	assert.NoError(t, err)
	assert.NotNil(t, createdProfessor)
	assert.Equal(t, professor.FirstName, createdProfessor.FirstName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update_Success(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"id", "matricule", "first_name", "last_name", "email", "department", "grade", "active", "status"}).
		AddRow(1, "12345", "John", "Doe", "john@example.com", "CS", "Professor", true, "active")

	mock.ExpectQuery(`SELECT \* FROM "professors" WHERE "professors"\."id" = \$1 ORDER BY "professors"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "professors"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	professor := Professor{
		Matricule:  "12345",
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.new@example.com",
		Department: "CS",
		Grade:      "Professor",
		Active:     true,
		Status:     "active",
	}

	updatedProfessor, err := repo.Update(1, professor)

	assert.NoError(t, err)
	assert.NotNil(t, updatedProfessor)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete_Success(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "professors"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAllPaginated_Success(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Mock count query
	mock.ExpectQuery(`SELECT count\(\*\) FROM "professors"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

	// Mock data query
	rows := sqlmock.NewRows([]string{"id", "matricule", "first_name", "last_name", "email", "department", "grade", "active", "status"}).
		AddRow(1, "12345", "John", "Doe", "john@example.com", "CS", "Professor", true, "active").
		AddRow(2, "12346", "Jane", "Smith", "jane@example.com", "Math", "Associate", true, "active")

	mock.ExpectQuery(`SELECT \* FROM "professors" LIMIT \$1`).
		WithArgs(10).
		WillReturnRows(rows)

	professors, total, err := repo.GetAllPaginated(1, 10)

	assert.NoError(t, err)
	assert.Len(t, professors, 2)
	assert.Equal(t, 10, total)
	assert.NoError(t, mock.ExpectationsWereMet())
}
