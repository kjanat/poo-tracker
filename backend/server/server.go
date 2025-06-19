package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/poo-tracker/backend/internal/middleware"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/service"
)

type App struct {
	Engine       *gin.Engine
	repo         repository.BowelMovementRepository
	details      repository.BowelMovementDetailsRepository
	meals        repository.MealRepository
	analytics    *service.Service
	userHandlers *UserAPIHandlers
}

func New(repo repository.BowelMovementRepository, details repository.BowelMovementDetailsRepository, meals repository.MealRepository, strategy service.AnalyticsStrategy, authService service.AuthService) *App {
	engine := gin.Default()
	userHandlers := NewUserAPIHandlers(authService)
	app := &App{
		Engine:       engine,
		repo:         repo,
		details:      details,
		meals:        meals,
		analytics:    service.New(repo, strategy),
		userHandlers: userHandlers,
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

	api.GET("/analytics", a.getAnalytics)

	// User management routes
	api.POST("/register", func(c *gin.Context) { a.userHandlers.RegisterHandler(c.Writer, c.Request) })
	api.POST("/login", func(c *gin.Context) { a.userHandlers.LoginHandler(c.Writer, c.Request) })
	api.GET("/profile", gin.WrapH(middleware.AuthMiddleware(a.userHandlers.AuthService)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { a.userHandlers.ProfileHandler(w, r) }))))
}
