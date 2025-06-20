package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/poo-tracker/backend/internal/model"
	"github.com/kjanat/poo-tracker/backend/internal/service"
)

// responseWriter wraps gin.ResponseWriter to capture response data
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func newResponseWriter(w gin.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
	}
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

// AuditMiddleware creates a middleware that logs API actions for audit purposes
func AuditMiddleware(auditService *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip audit for certain paths
		if shouldSkipAudit(c.Request.URL.Path) {
			c.Next()
			return
		}

		userID := c.GetString("userID")
		if userID == "" {
			c.Next()
			return
		}

		// Read request body if it exists
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap response writer to capture response
		writer := newResponseWriter(c.Writer)
		c.Writer = writer

		// Process the request
		c.Next()

		// Determine action and entity type from request
		action, entityType, entityID := determineAuditInfo(c)
		if action == "" || entityType == "" {
			return
		}

		// Parse request and response data
		var oldData, newData interface{}

		if len(requestBody) > 0 {
			// Ignore JSON parsing errors for audit logging
			_ = json.Unmarshal(requestBody, &newData)
		}

		// For successful responses, capture the response data
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			var responseData interface{}
			if writer.body.Len() > 0 {
				// Ignore JSON parsing errors for audit logging
				if err := json.Unmarshal(writer.body.Bytes(), &responseData); err == nil {
					switch action {
					case model.AuditCreate, model.AuditUpdate:
						newData = responseData
					}
				}
			}
		}

		// Log the action
		auditService.LogAction(c.Request.Context(), userID, entityType, entityID, action, oldData, newData)
	}
}

// shouldSkipAudit determines if audit logging should be skipped for a path
func shouldSkipAudit(path string) bool {
	skipPaths := []string{
		"/health",
		"/api/login",
		"/api/register",
		"/api/profile",
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	return false
}

// determineAuditInfo extracts audit information from the request
func determineAuditInfo(c *gin.Context) (model.AuditAction, string, string) {
	method := c.Request.Method
	path := c.Request.URL.Path
	id := c.Param("id")

	var action model.AuditAction
	var entityType string

	// Determine action from HTTP method
	switch method {
	case "POST":
		action = model.AuditCreate
	case "PUT", "PATCH":
		action = model.AuditUpdate
	case "DELETE":
		action = model.AuditDelete
	default:
		return "", "", "" // Don't audit GET requests
	}

	// Determine entity type from path
	switch {
	case strings.Contains(path, "/bowel-movements"):
		entityType = "bowel_movement"
	case strings.Contains(path, "/meals"):
		entityType = "meal"
	case strings.Contains(path, "/symptoms"):
		entityType = "symptom"
	case strings.Contains(path, "/medications"):
		entityType = "medication"
	default:
		return "", "", ""
	}

	return action, entityType, id
}
