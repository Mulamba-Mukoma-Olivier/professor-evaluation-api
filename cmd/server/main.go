package main

import (
	"log"
	"os"
	"time"

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
	// =========================================================
	// 1. Charger les variables d'environnement
	// =========================================================

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// =========================================================
	// 2. Configuration PostgreSQL
	// =========================================================

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
		log.Fatalf("Database connection failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	defer sqlDB.Close()

	log.Println("Database connected successfully")

	// =========================================================
	// 3. Database migrations
	// =========================================================

	if err := database.Migrate(db, dbConfig); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	log.Println("Database migration completed successfully")

	// =========================================================
	// 4. JWT
	// =========================================================

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	jwtManager := appjwt.NewManager(
		jwtSecret,
		24*time.Hour,
	)

	// =========================================================
	// 5. Repositories
	// =========================================================

	userRepository := auth.NewRepository(db)
	professorRepository := professors.NewRepository(db)
	courseRepository := courses.NewRepository(db)
	criteriaRepository := criteria.NewRepository(db)
	evaluationRepository := evaluations.NewRepository(db)
	eligibilityRepository := eligibility.NewRepository(db)
	resultRepository := results.NewRepository(db)

	// =========================================================
	// 6. Services
	// =========================================================

	authService := auth.NewService(
		userRepository,
	)

	professorService := professors.NewService(
		professorRepository,
	)

	courseService := courses.NewService(
		courseRepository,
	)

	criteriaService := criteria.NewService(
		criteriaRepository,
	)

	evaluationService := evaluations.NewService(
		evaluationRepository,
	)

	eligibilityService := eligibility.NewService(
		eligibilityRepository,
	)

	resultService := results.NewService(
		resultRepository,
	)

	// =========================================================
	// 7. Handlers
	// =========================================================

	authHandler := auth.NewHandler(
		authService,
		jwtManager,
	)

	professorHandler := professors.NewHandler(
		professorService,
	)

	courseHandler := courses.NewHandler(
		courseService,
	)

	criteriaHandler := criteria.NewHandler(
		criteriaService,
	)

	evaluationHandler := evaluations.NewHandler(
		evaluationService,
	)

	eligibilityHandler := eligibility.NewHandler(
		eligibilityService,
	)

	resultHandler := results.NewHandler(
		resultService,
	)

	// =========================================================
	// 8. Router
	// =========================================================

	router := routes.SetupRouter(
		routes.RouterDependencies{
			AuthHandler:      authHandler,
			CourseHandler:    courseHandler,
			CriteriaHandler:  criteriaHandler,
			Eligibility:      eligibilityHandler,
			Evaluation:      evaluationHandler,
			ProfessorHandler: professorHandler,
			ResultHandler:    resultHandler,
			JWTManager:       jwtManager,
		},
	)

	// =========================================================
	// 9. Serveur HTTP
	// =========================================================

	port := getEnv("PORT", "8080")

	log.Printf(
		"Professor Evaluation API running on http://localhost:%s",
		port,
	)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// getEnv récupère une variable d'environnement.
// Si elle n'existe pas, la valeur par défaut est utilisée.
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}