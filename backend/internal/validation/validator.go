// Package validation provides comprehensive validation functions for the poo-tracker backend.
// This file serves as a facade to re-export all validation functions from the modular structure.
package validation

// This file re-exports all validation functions from the split modules for backward compatibility.
// The validation package is now organized into focused modules:
//
// - errors.go: ValidationError and ValidationErrors types
// - basic.go: Basic validation functions (Bristol type, scales, enums)
// - user.go: User-related validations (email, password, name, etc.)
// - content.go: Content validation (notes, URLs, meal names, etc.)
// - security.go: Security-focused validation (strong passwords, policies)
// - models.go: Complex model validation (BowelMovement, Meal)
//
// All functions are available at the package level for existing code compatibility.
