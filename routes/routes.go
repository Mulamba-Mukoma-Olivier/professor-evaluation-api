package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/auth"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/courses"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/criteria"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/eligibility"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/professors"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/results"

	appjwt "github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/pkg/jwt"
	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/pkg/middleware"
)

type RouterDependencies struct {
	AuthHandler      *auth.Handler
	CourseHandler    *courses.Handler
	CriteriaHandler  *criteria.Handler
	Eligibility      *eligibility.Handler
	Evaluation       *evaluations.Handler
	ProfessorHandler *professors.Handler
	ResultHandler    *results.Handler
	JWTManager       *appjwt.Manager
}

func SetupRouter(deps RouterDependencies) *gin.Engine {

	router := gin.Default()

	// =========================================================
	// CORS
	// =========================================================

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:3002",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// =========================================================
	// TRUSTED PROXIES
	// =========================================================

	if err := router.SetTrustedProxies([]string{
		"127.0.0.1",
		"::1",
	}); err != nil {
		panic(err)
	}

	// =========================================================
	// HEALTH CHECK
	// =========================================================

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Professor Evaluation API is running",
		})
	})

	// =========================================================
	// PUBLIC ROUTES
	// =========================================================

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", deps.AuthHandler.Login)
		authRoutes.POST("/register", deps.AuthHandler.Register)
	}

	// =========================================================
	// PROTECTED ROUTES
	// =========================================================

	protected := router.Group("")
	protected.Use(
		middleware.AuthRequired(deps.JWTManager),
	)

	// ---------------------------------------------------------
	// Eligibility
	// ---------------------------------------------------------

	protected.GET(
		"/eligibility/:student_id",
		deps.Eligibility.Check,
	)

	// ---------------------------------------------------------
	// Professors
	// ---------------------------------------------------------

	protected.GET(
		"/professors",
		deps.ProfessorHandler.GetAll,
	)

	protected.GET(
		"/professors/active",
		deps.ProfessorHandler.GetActive,
	)

	protected.GET(
		"/professors/:id",
		deps.ProfessorHandler.GetByID,
	)

	protected.POST(
		"/professors",
		middleware.RequireRole("ADMIN", "SUPER_ADMIN"),
		deps.ProfessorHandler.Create,
	)

	protected.PUT(
		"/professors/:id",
		middleware.RequireRole("ADMIN", "SUPER_ADMIN"),
		deps.ProfessorHandler.Update,
	)

	protected.DELETE(
		"/professors/:id",
		middleware.RequireRole("ADMIN", "SUPER_ADMIN"),
		deps.ProfessorHandler.Delete,
	)

	// ---------------------------------------------------------
	// Courses
	// ---------------------------------------------------------

	protected.GET(
		"/courses",
		deps.CourseHandler.GetAll,
	)

	protected.GET(
		"/courses/:id",
		deps.CourseHandler.GetByID,
	)

	protected.POST(
		"/courses",
		middleware.RequireRole("ADMIN", "SUPER_ADMIN"),
		deps.CourseHandler.Create,
	)

	// ---------------------------------------------------------
	// Criteria
	// ---------------------------------------------------------

	protected.GET(
		"/criteria",
		deps.CriteriaHandler.GetAll,
	)

	protected.GET(
		"/criteria/active",
		deps.CriteriaHandler.GetActive,
	)

	protected.GET(
		"/criteria/:id",
		deps.CriteriaHandler.GetByID,
	)

	protected.POST(
		"/criteria",
		middleware.RequireRole("ADMIN", "SUPER_ADMIN"),
		deps.CriteriaHandler.Create,
	)

	protected.PUT(
		"/criteria/:id",
		middleware.RequireRole("ADMIN", "SUPER_ADMIN"),
		deps.CriteriaHandler.Update,
	)

	protected.DELETE(
		"/criteria/:id",
		middleware.RequireRole("ADMIN", "SUPER_ADMIN"),
		deps.CriteriaHandler.Delete,
	)

	// ---------------------------------------------------------
	// Evaluations
	// ---------------------------------------------------------

	protected.GET(
		"/evaluations",
		middleware.RequireRole("ADMIN", "SUPER_ADMIN"),
		deps.Evaluation.GetAll,
	)

	protected.GET(
		"/evaluations/:id",
		middleware.RequireRole("ADMIN", "SUPER_ADMIN"),
		deps.Evaluation.GetByID,
	)

	protected.POST(
		"/evaluations/:student_id",
		middleware.RequireRole("STUDENT"),
		deps.Evaluation.Create,
	)

	// ---------------------------------------------------------
	// Results
	// ---------------------------------------------------------

	protected.GET(
		"/results/professors/:professor_id",
		middleware.RequireRole("ADMIN", "SUPER_ADMIN"),
		deps.ResultHandler.GetProfessorResult,
	)

	return router
}