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
			"http://localhost:5500",
			"http://127.0.0.1:5500",

			"http://localhost:3000",
			"http://127.0.0.1:3000",

			"http://localhost:3001",
			"http://127.0.0.1:3001",

			"http://localhost:3002",
			"http://127.0.0.1:3002",
		},

		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
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

		MaxAge: 12 * time.Hour,
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
	// RATE LIMITING
	// =========================================================

	rateLimiter := middleware.NewRateLimiter(
		100,
		time.Minute,
	)

	router.Use(rateLimiter.Middleware())

	// =========================================================
	// HEALTH
	// =========================================================

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Professor Evaluation API is running",
		})
	})

	// =========================================================
	// AUTHENTIFICATION
	// =========================================================

	authRoutes := router.Group("/auth")

	{
		authRoutes.POST(
			"/login",
			deps.AuthHandler.Login,
		)

		authRoutes.POST(
			"/register",
			deps.AuthHandler.Register,
		)
	}

	// =========================================================
	// ROUTES PROTÉGÉES
	// =========================================================

	protected := router.Group("")

	protected.Use(
		middleware.AuthRequired(
			deps.JWTManager,
		),
	)

	// =========================================================
	// ELIGIBILITY
	// =========================================================

	// ---------------------------------------------------------
	// Consultation
	// STUDENT, ADMIN, SUPER_ADMIN, PROFESSOR
	// ---------------------------------------------------------

	protected.GET(
		"/eligibility/:student_id",
		middleware.RequireRole(
			"STUDENT",
			"ADMIN",
			"SUPER_ADMIN",
			"PROFESSOR",
		),
		deps.Eligibility.Check,
	)

	// ---------------------------------------------------------
	// Création
	// ADMIN, SUPER_ADMIN
	// ---------------------------------------------------------

	protected.POST(
		"/eligibility",
		middleware.RequireRole(
			"ADMIN",
			"SUPER_ADMIN",
		),
		deps.Eligibility.Create,
	)

	// ---------------------------------------------------------
	// Modification
	// ADMIN, SUPER_ADMIN
	// ---------------------------------------------------------

	protected.PUT(
		"/eligibility/:student_id",
		middleware.RequireRole(
			"ADMIN",
			"SUPER_ADMIN",
		),
		deps.Eligibility.Update,
	)

	// ---------------------------------------------------------
	// Suppression
	// ADMIN, SUPER_ADMIN
	// ---------------------------------------------------------

	protected.DELETE(
		"/eligibility/:student_id",
		middleware.RequireRole(
			"ADMIN",
			"SUPER_ADMIN",
		),
		deps.Eligibility.Delete,
	)

	// =========================================================
	// PROFESSORS
	// =========================================================

	// ---------------------------------------------------------
	// Lecture
	// STUDENT, ADMIN, SUPER_ADMIN, PROFESSOR
	// ---------------------------------------------------------

	protected.GET(
		"/professors",
		middleware.RequireRole(
			"STUDENT",
			"ADMIN",
			"SUPER_ADMIN",
			"PROFESSOR",
		),
		deps.ProfessorHandler.GetAll,
	)

	protected.GET(
		"/professors/active",
		middleware.RequireRole(
			"STUDENT",
			"ADMIN",
			"SUPER_ADMIN",
			"PROFESSOR",
		),
		deps.ProfessorHandler.GetActive,
	)

	protected.GET(
		"/professors/by-status",
		middleware.RequireRole(
			"STUDENT",
			"ADMIN",
			"SUPER_ADMIN",
			"PROFESSOR",
		),
		deps.ProfessorHandler.GetByStatus,
	)

	protected.GET(
		"/professors/:id",
		middleware.RequireRole(
			"STUDENT",
			"ADMIN",
			"SUPER_ADMIN",
			"PROFESSOR",
		),
		deps.ProfessorHandler.GetByID,
	)

	// ---------------------------------------------------------
	// Gestion
	// ADMIN, SUPER_ADMIN
	// ---------------------------------------------------------

	protected.POST(
		"/professors",
		middleware.RequireRole(
			"ADMIN",
			"SUPER_ADMIN",
		),
		deps.ProfessorHandler.Create,
	)

	protected.PUT(
		"/professors/:id",
		middleware.RequireRole(
			"ADMIN",
			"SUPER_ADMIN",
		),
		deps.ProfessorHandler.Update,
	)

	protected.DELETE(
		"/professors/:id",
		middleware.RequireRole(
			"ADMIN",
			"SUPER_ADMIN",
		),
		deps.ProfessorHandler.Delete,
	)

	// =========================================================
	// COURSES
	// =========================================================

	// ---------------------------------------------------------
	// Lecture
	// ---------------------------------------------------------

	protected.GET(
		"/courses",
		middleware.RequireRole(
			"STUDENT",
			"ADMIN",
			"SUPER_ADMIN",
			"PROFESSOR",
		),
		deps.CourseHandler.GetAll,
	)

	protected.GET(
		"/courses/:id",
		middleware.RequireRole(
			"STUDENT",
			"ADMIN",
			"SUPER_ADMIN",
			"PROFESSOR",
		),
		deps.CourseHandler.GetByID,
	)

	// ---------------------------------------------------------
	// Gestion
	// ADMIN, SUPER_ADMIN
	// ---------------------------------------------------------

	protected.POST(
		"/courses",
		middleware.RequireRole(
			"ADMIN",
			"SUPER_ADMIN",
		),
		deps.CourseHandler.Create,
	)

	protected.PUT(
		"/courses/:id",
		middleware.RequireRole(
			"ADMIN",
			"SUPER_ADMIN",
		),
		deps.CourseHandler.Update,
	)

	protected.DELETE(
		"/courses/:id",
		middleware.RequireRole(
			"ADMIN",
			"SUPER_ADMIN",
		),
		deps.CourseHandler.Delete,
	)

	// =========================================================
	// CRITERIA
	// =========================================================

	// ---------------------------------------------------------
	// Lecture
	// ---------------------------------------------------------

	protected.GET(
		"/criteria",
		middleware.RequireRole(
			"STUDENT",
			"ADMIN",
			"SUPER_ADMIN",
			"PROFESSOR",
		),
		deps.CriteriaHandler.GetAll,
	)

	protected.GET(
		"/criteria/active",
		middleware.RequireRole(
			"STUDENT",
			"ADMIN",
			"SUPER_ADMIN",
			"PROFESSOR",
		),
		deps.CriteriaHandler.GetActive,
	)

	protected.GET(
		"/criteria/:id",
		middleware.RequireRole(
			"STUDENT",
			"ADMIN",
			"SUPER_ADMIN",
			"PROFESSOR",
		),
		deps.CriteriaHandler.GetByID,
	)

	// ---------------------------------------------------------
	// Gestion
	// ADMIN, SUPER_ADMIN
	// ---------------------------------------------------------

	protected.POST(
		"/criteria",
		middleware.RequireRole(
			"ADMIN",
			"SUPER_ADMIN",
		),
		deps.CriteriaHandler.Create,
	)

	protected.PUT(
		"/criteria/:id",
		middleware.RequireRole(
			"ADMIN",
			"SUPER_ADMIN",
		),
		deps.CriteriaHandler.Update,
	)

	protected.DELETE(
		"/criteria/:id",
		middleware.RequireRole(
			"ADMIN",
			"SUPER_ADMIN",
		),
		deps.CriteriaHandler.Delete,
	)

	// =========================================================
	// EVALUATIONS
	// =========================================================

	// ---------------------------------------------------------
	// Consultation administrative
	// ADMIN, SUPER_ADMIN
	// ---------------------------------------------------------

	protected.GET(
		"/evaluations",
		middleware.RequireRole(
			"ADMIN",
			"SUPER_ADMIN",
		),
		deps.Evaluation.GetAll,
	)

	protected.GET(
		"/evaluations/:id",
		middleware.RequireRole(
			"ADMIN",
			"SUPER_ADMIN",
		),
		deps.Evaluation.GetByID,
	)

	// ---------------------------------------------------------
	// Création
	// STUDENT
	// ---------------------------------------------------------

	protected.POST(
		"/evaluations",
		middleware.RequireRole(
			"STUDENT",
		),
		deps.Evaluation.Create,
	)

	// =========================================================
	// RESULTS
	// =========================================================

	// ---------------------------------------------------------
	// Consultation des résultats
	// STUDENT, ADMIN, SUPER_ADMIN, PROFESSOR
	// ---------------------------------------------------------

	protected.GET(
		"/results/professors/:professor_id",
		middleware.RequireRole(
			"STUDENT",
			"ADMIN",
			"SUPER_ADMIN",
			"PROFESSOR",
		),
		deps.ResultHandler.GetProfessorResult,
	)

	// =========================================================
	// RETURN ROUTER
	// =========================================================

	return router
}