package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/auth"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/courses"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/criteria"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/eligibility"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/professors"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/results"

	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/pkg/database"
	appjwt "github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/pkg/jwt"

	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/routes"
)

func main() {

	// ========================================
	// ENVIRONMENT
	// ========================================

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// ========================================
	// GIN
	// ========================================

	gin.SetMode(gin.DebugMode)

	// ========================================
	// DATABASE
	// ========================================

	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Name:     getEnv("DB_NAME", "professor_evaluation"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	log.Println("PostgreSQL connected successfully")

	// ========================================
	// DATABASE MIGRATION
	// ========================================

	if err := database.Migrate(db); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	log.Println("Database migration completed")

	// ========================================
	// JWT
	// ========================================

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	jwtManager := appjwt.NewManager(
		jwtSecret,
		24*time.Hour,
	)

	// ========================================
	// AUTH
	// ========================================

	authRepository := auth.NewRepository()

	authService := auth.NewService(
		authRepository,
	)

	authHandler := auth.NewHandler(
		authService,
		jwtManager,
	)

	// ========================================
	// COURSES
	// ========================================

	courseRepository := courses.NewRepository(db)

	courseService := courses.NewService(
		courseRepository,
	)

	courseHandler := courses.NewHandler(
		courseService,
	)

	// ========================================
	// CRITERIA
	// ========================================

	criteriaRepository := criteria.NewRepository(db)

	criteriaService := criteria.NewService(
		criteriaRepository,
	)

	criteriaHandler := criteria.NewHandler(
		criteriaService,
	)

	// ========================================
	// ELIGIBILITY
	// ========================================

	eligibilityRepository := eligibility.NewRepository(db)

	eligibilityService := eligibility.NewService(
		eligibilityRepository,
	)

	eligibilityHandler := eligibility.NewHandler(
		eligibilityService,
	)

	// ========================================
	// PROFESSORS
	// ========================================

	professorRepository := professors.NewRepository(db)

	professorService := professors.NewService(
		professorRepository,
	)

	professorHandler := professors.NewHandler(
		professorService,
	)

	// ========================================
	// EVALUATIONS
	// ========================================

	evaluationRepository := evaluations.NewRepository(db)

	evaluationService := evaluations.NewService(
		evaluationRepository,
	)

	evaluationHandler := evaluations.NewHandler(
		evaluationService,
	)

	// ========================================
	// RESULTS
	// ========================================

	resultRepository := results.NewRepository(db)

	resultService := results.NewService(
		resultRepository,
	)

	resultHandler := results.NewHandler(
		resultService,
	)

	// ========================================
	// ROUTER
	// ========================================

	router := routes.SetupRouter(
		routes.RouterDependencies{
			AuthHandler:      authHandler,
			CourseHandler:    courseHandler,
			CriteriaHandler:  criteriaHandler,
			Eligibility:      eligibilityHandler,
			Evaluation:       evaluationHandler,
			ProfessorHandler: professorHandler,
			ResultHandler:    resultHandler,
			JWTManager:       jwtManager,
		},
	)

	// ========================================
	// SERVER
	// ========================================

	port := getEnv("PORT", "8080")

	log.Println("========================================")
	log.Println("Professor Evaluation API")
	log.Printf("Server running on http://localhost:%s", port)
	log.Println("========================================")

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

// ========================================
// ENVIRONMENT HELPER
// ========================================

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}