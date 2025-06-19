package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/service"
)

type App struct {
	Engine    *gin.Engine
	repo      repository.BowelMovementRepository
	analytics *service.Service
}

func New(repo repository.BowelMovementRepository, strategy service.AnalyticsStrategy) *App {
	engine := gin.Default()
	app := &App{Engine: engine, repo: repo, analytics: service.New(repo, strategy)}

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
	bm.GET(":id", a.getBowelMovement)
	bm.PUT(":id", a.updateBowelMovement)
	bm.DELETE(":id", a.deleteBowelMovement)

	api.GET("/analytics", a.getAnalytics)
}
