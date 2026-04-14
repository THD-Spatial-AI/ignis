package models

import "fmt"

// Custom error types for better error handling

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// CalculationError represents a calculation error
type CalculationError struct {
	Level   string
	Message string
}

func (e *CalculationError) Error() string {
	return fmt.Sprintf("calculation error at level '%s': %s", e.Level, e.Message)
}

// NewCalculationError creates a new calculation error
func NewCalculationError(level, message string) *CalculationError {
	return &CalculationError{
		Level:   level,
		Message: message,
	}
}

// DatabaseError represents a database error
type DatabaseError struct {
	Operation string
	Message   string
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during '%s': %s", e.Operation, e.Message)
}

// NewDatabaseError creates a new database error
func NewDatabaseError(operation, message string) *DatabaseError {
	return &DatabaseError{
		Operation: operation,
		Message:   message,
	}
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID '%s' not found", e.Resource, e.ID)
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}
