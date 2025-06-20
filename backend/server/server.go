package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/poo-tracker/backend/internal/middleware"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/service"
)

type App struct {
	Engine               *gin.Engine
	repo                 repository.BowelMovementRepository
	details              repository.BowelMovementDetailsRepository
	meals                repository.MealRepository
	symptoms             repository.SymptomRepository
	medications          repository.MedicationRepository
	mealBowelRelations   repository.MealBowelMovementRelationRepository
	mealSymptomRelations repository.MealSymptomRelationRepository
	analytics            *service.Service
	authService          service.AuthService
	auditService         *service.AuditService
	userHandlers         *UserAPIHandlers
	symptomHandler       *SymptomHandler
	medicationHandler    *MedicationHandler
	// relationsHandler     *RelationsHandler // TODO: implement
}

func New(repo repository.BowelMovementRepository, details repository.BowelMovementDetailsRepository, meals repository.MealRepository, symptoms repository.SymptomRepository, medications repository.MedicationRepository, mealBowelRelations repository.MealBowelMovementRelationRepository, mealSymptomRelations repository.MealSymptomRelationRepository, strategy service.AnalyticsStrategy, authService service.AuthService) *App {
	engine := gin.Default()

	// Add global middleware
	engine.Use(middleware.SecurityHeadersMiddleware())
	rateLimiter := middleware.NewRateLimiter(100, time.Minute) // 100 requests per minute
	engine.Use(middleware.RateLimitMiddleware(rateLimiter))

	auditService := service.NewAuditService()
	engine.Use(middleware.AuditMiddleware(auditService))

	userHandlers := NewUserAPIHandlers(authService)
	symptomHandler := NewSymptomHandler(symptoms)
	medicationHandler := NewMedicationHandler(medications)
	// relationsHandler := NewRelationsHandler(mealBowelRelations, mealSymptomRelations) // TODO: implement
	app := &App{
		Engine:               engine,
		repo:                 repo,
		details:              details,
		meals:                meals,
		symptoms:             symptoms,
		medications:          medications,
		mealBowelRelations:   mealBowelRelations,
		mealSymptomRelations: mealSymptomRelations,
		analytics:            service.New(repo, strategy),
		authService:          authService,
		auditService:         auditService,
		userHandlers:         userHandlers,
		symptomHandler:       symptomHandler,
		medicationHandler:    medicationHandler,
		// relationsHandler:     relationsHandler, // TODO: implement
	}

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	app.registerRoutes()
	return app
}

func (a *App) registerRoutes() {
	api := a.Engine.Group("/api")
	bm := api.Group("/bowel-movements")
	bm.GET("", a.listBowelMovements)
	bm.POST("", a.createBowelMovement)
	bm.GET("/:id", a.getBowelMovement)
	bm.PUT("/:id", a.updateBowelMovement)
	bm.DELETE("/:id", a.deleteBowelMovement)

	// BowelMovementDetails routes
	bm.POST("/:id/details", a.createBowelMovementDetails)
	bm.GET("/:id/details", a.getBowelMovementDetails)
	bm.PUT("/:id/details", a.updateBowelMovementDetails)
	bm.DELETE("/:id/details", a.deleteBowelMovementDetails)

	meals := api.Group("/meals")
	meals.GET("", a.listMeals)
	meals.POST("", a.createMeal)
	meals.GET("/:id", a.getMeal)
	meals.PUT("/:id", a.updateMeal)
	meals.DELETE("/:id", a.deleteMeal)

	// Symptom routes
	symptoms := api.Group("/symptoms")
	symptoms.Use(middleware.JWTAuthMiddleware(a.authService))
	symptoms.GET("", a.symptomHandler.GetSymptoms)
	symptoms.POST("", a.symptomHandler.CreateSymptom)
	symptoms.GET("/:id", a.symptomHandler.GetSymptom)
	symptoms.PUT("/:id", a.symptomHandler.UpdateSymptom)
	symptoms.DELETE("/:id", a.symptomHandler.DeleteSymptom)
	symptoms.GET("/date-range", a.symptomHandler.GetSymptomsByDateRange)
	symptoms.GET("/category/:category", a.symptomHandler.GetSymptomsByCategory)
	symptoms.GET("/type/:type", a.symptomHandler.GetSymptomsByType)

	// Medication routes
	medications := api.Group("/medications")
	medications.Use(middleware.JWTAuthMiddleware(a.authService))
	medications.GET("", a.medicationHandler.GetMedications)
	medications.POST("", a.medicationHandler.CreateMedication)
	medications.GET("/:id", a.medicationHandler.GetMedication)
	medications.PUT("/:id", a.medicationHandler.UpdateMedication)
	medications.DELETE("/:id", a.medicationHandler.DeleteMedication)
	medications.GET("/active", a.medicationHandler.GetActiveMedications)
	medications.POST("/:id/taken", a.medicationHandler.MarkMedicationAsTaken)
	medications.GET("/category/:category", a.medicationHandler.GetMedicationsByCategory)

	api.GET("/analytics", a.getAnalytics)

	// User management routes
	api.POST("/register", func(c *gin.Context) { a.userHandlers.RegisterHandler(c.Writer, c.Request) })
	api.POST("/login", func(c *gin.Context) { a.userHandlers.LoginHandler(c.Writer, c.Request) })
	api.GET("/profile", gin.WrapH(middleware.AuthMiddleware(a.userHandlers.AuthService)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { a.userHandlers.ProfileHandler(w, r) }))))
}
