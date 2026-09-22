package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/auth"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/courses"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/criteria"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/eligibility"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/professors"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/results"

	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/pkg/database"
	appjwt "github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/pkg/jwt"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/pkg/logger"

	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/routes"
)

func main() {

	// ========================================
	// LOGGER INITIALIZATION
	// ========================================

	if err := logger.InitLogger(); err != nil {
		log.Fatal("failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Info("Starting Professor Evaluation API")

	// ========================================
	// ENVIRONMENT
	// ========================================

	if err := godotenv.Load(); err != nil {
		logger.Warn(".env file not found")
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
		logger.Fatal("failed to connect to database", 
			zap.String("error", err.Error()))
	}

	logger.Info("PostgreSQL connected successfully")

	// ========================================
	// DATABASE MIGRATION
	// ========================================

	if err := database.RunMigrations(dbConfig); err != nil {
		logger.Fatal("failed to run database migrations",
			zap.String("error", err.Error()))
	}

	logger.Info("Database migrations completed")

	// ========================================
	// JWT
	// ========================================

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		logger.Fatal("JWT_SECRET is required")
	}

	jwtManager := appjwt.NewManager(
		jwtSecret,
		24*time.Hour,
	)

	// ========================================
	// AUTH
	// ========================================

	authRepository := auth.NewRepository(db)

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

	logger.Info("Server starting",
		zap.String("port", port),
		zap.String("host", "localhost"))

	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Server failed to start",
			zap.String("error", err.Error()))
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