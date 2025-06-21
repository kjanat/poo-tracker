package shared

import (
	"time"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error     ErrorDetail `json:"error"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorDetail contains detailed error information
type ErrorDetail struct {
	Code        string       `json:"code"`
	Message     string       `json:"message"`
	Details     string       `json:"details,omitempty"`
	Fields      []FieldError `json:"fields,omitempty"`
	Suggestions []string     `json:"suggestions,omitempty"`
}

// FieldError represents a field-specific validation error
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
}

// SuccessResponse represents a standardized success response
type SuccessResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
	Data      any       `json:"data,omitempty"`
	RequestID string    `json:"request_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
	TotalCount int64 `json:"total_count"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       any            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
	RequestID  string         `json:"request_id,omitempty"`
	Timestamp  time.Time      `json:"timestamp"`
}

// BaseListRequest represents common list request parameters
type BaseListRequest struct {
	Page     int    `json:"page,omitempty" form:"page" binding:"omitempty,min=1"`
	PageSize int    `json:"page_size,omitempty" form:"page_size" binding:"omitempty,min=1,max=100"`
	SortBy   string `json:"sort_by,omitempty" form:"sort_by"`
	SortDir  string `json:"sort_dir,omitempty" form:"sort_dir" binding:"omitempty,oneof=asc desc"`
	Search   string `json:"search,omitempty" form:"search"`
}

// DateRangeRequest represents common date range parameters
type DateRangeRequest struct {
	StartDate *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty" form:"end_date"`
}

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status      string            `json:"status"`
	Version     string            `json:"version"`
	Environment string            `json:"environment"`
	Timestamp   time.Time         `json:"timestamp"`
	Services    map[string]string `json:"services"`
	Uptime      string            `json:"uptime"`
}

// NewErrorResponse creates a new standardized error response
func NewErrorResponse(code, message, details string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
	}
}

// NewValidationErrorResponse creates a new validation error response
func NewValidationErrorResponse(fieldErrors []FieldError) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
			Fields:  fieldErrors,
		},
		Timestamp: time.Now(),
	}
}

// NewSuccessResponse creates a new standardized success response
func NewSuccessResponse(message string, data any) SuccessResponse {
	return SuccessResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data any, pagination PaginationMeta) PaginatedResponse {
	return PaginatedResponse{
		Data:       data,
		Pagination: pagination,
		Timestamp:  time.Now(),
	}
}

// NewPaginationMeta creates pagination metadata
func NewPaginationMeta(page, pageSize int, totalCount int64) PaginationMeta {
	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))
	if totalPages < 1 {
		totalPages = 1
	}

	return PaginationMeta{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalCount: totalCount,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// GetDefaultPagination returns default pagination values
func (r *BaseListRequest) GetDefaultPagination() (page, pageSize int) {
	page = 1
	if r.Page > 0 {
		page = r.Page
	}

	pageSize = 20
	if r.PageSize > 0 {
		pageSize = r.PageSize
	}

	return page, pageSize
}

// GetDefaultSorting returns default sorting values
func (r *BaseListRequest) GetDefaultSorting() (sortBy, sortDir string) {
	sortBy = "created_at"
	if r.SortBy != "" {
		sortBy = r.SortBy
	}

	sortDir = "desc"
	if r.SortDir != "" {
		sortDir = r.SortDir
	}

	return sortBy, sortDir
}

// Validate validates the BaseListRequest
func (r *BaseListRequest) Validate() error {
	if r.Page < 1 {
		return NewValidationError("page", "page must be greater than 0")
	}
	if r.PageSize < 1 || r.PageSize > 100 {
		return NewValidationError("page_size", "page_size must be between 1 and 100")
	}
	return nil
}

// Validate validates the DateRangeRequest
func (r *DateRangeRequest) Validate() error {
	if r.StartDate != nil && r.EndDate != nil && r.EndDate.Before(*r.StartDate) {
		return NewValidationError("date_range", "end_date must be after start_date")
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}
