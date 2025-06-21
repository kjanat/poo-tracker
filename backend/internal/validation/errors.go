package validation

import "fmt"

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// ValidationErrors holds multiple validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	if len(e) == 1 {
		return e[0].Error()
	}
	msg := fmt.Sprintf("%d validation errors: ", len(e))
	for i, err := range e {
		if i > 0 {
			msg += "; "
		}
		msg += err.Error()
	}
	return msg
}

// HasErrors returns true if there are any validation errors
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// Add adds a new validation error
func (e *ValidationErrors) Add(field, message string) {
	*e = append(*e, ValidationError{Field: field, Message: message})
}

// FormatValidationError formats validation errors from gin binding
func FormatValidationError(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// safeAppendValidationError safely appends an error to ValidationErrors if it's a ValidationError
func safeAppendValidationError(errors *ValidationErrors, err error) {
	if validationErr, ok := err.(ValidationError); ok {
		*errors = append(*errors, validationErr)
	}
}
