package server

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BowelMovement represents a bowel movement entry
// This is a simplified version of the previous Node.js model
// ID is generated on creation
// UserID would normally come from authentication middleware
// For demo purposes we keep it simple

type BowelMovement struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	BristolType int       `json:"bristolType"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// In-memory store and mutex for thread safety
var (
	store = make(map[string]BowelMovement)
	mu    sync.RWMutex
)

func registerRoutes(r *gin.Engine) {
	api := r.Group("/api")
	bm := api.Group("/bowel-movements")
	bm.GET("", listBowelMovements)
	bm.POST("", createBowelMovement)
	bm.GET(":id", getBowelMovement)
	bm.PUT(":id", updateBowelMovement)
	bm.DELETE(":id", deleteBowelMovement)

	api.GET("/analytics", getAnalytics)
}

func listBowelMovements(c *gin.Context) {
	mu.RLock()
	defer mu.RUnlock()
	movements := make([]BowelMovement, 0, len(store))
	for _, m := range store {
		movements = append(movements, m)
	}
	c.JSON(http.StatusOK, gin.H{"data": movements})
}

func createBowelMovement(c *gin.Context) {
	var req struct {
		UserID      string `json:"userId"`
		BristolType int    `json:"bristolType"`
		Notes       string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.NewString()
	now := time.Now().UTC()
	bm := BowelMovement{
		ID:          id,
		UserID:      req.UserID,
		BristolType: req.BristolType,
		Notes:       req.Notes,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mu.Lock()
	store[id] = bm
	mu.Unlock()

	c.JSON(http.StatusCreated, bm)
}

func getBowelMovement(c *gin.Context) {
	id := c.Param("id")
	mu.RLock()
	bm, ok := store[id]
	mu.RUnlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, bm)
}

func updateBowelMovement(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		BristolType *int    `json:"bristolType"`
		Notes       *string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	bm, ok := store[id]
	if !ok {
		mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if req.BristolType != nil {
		bm.BristolType = *req.BristolType
	}
	if req.Notes != nil {
		bm.Notes = *req.Notes
	}
	bm.UpdatedAt = time.Now().UTC()
	store[id] = bm
	mu.Unlock()

	c.JSON(http.StatusOK, bm)
}

func deleteBowelMovement(c *gin.Context) {
	id := c.Param("id")
	mu.Lock()
	_, ok := store[id]
	if ok {
		delete(store, id)
	}
	mu.Unlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func getAnalytics(c *gin.Context) {
	mu.RLock()
	defer mu.RUnlock()
	total := len(store)
	if total == 0 {
		c.JSON(http.StatusOK, gin.H{"total": 0})
		return
	}
	sum := 0
	for _, bm := range store {
		sum += bm.BristolType
	}
	avg := float64(sum) / float64(total)
	c.JSON(http.StatusOK, gin.H{"total": total, "avgBristol": avg})
}
