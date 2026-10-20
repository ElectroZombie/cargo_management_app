// Package errors defines the domain error types shared by the database,
// service, and Wails binding layers.
//
// Every domain error unwraps to one of the sentinel values below, so callers
// can classify a failure with errors.Is even after an intermediate layer has
// wrapped it with fmt.Errorf.
package errors

import (
	stderrors "errors"
	"fmt"
)

// Sentinel errors used to classify domain failures.
var (
	ErrValidation = stderrors.New("validation failed")
	ErrNotFound   = stderrors.New("record not found")
	ErrConflict   = stderrors.New("record conflicts with existing data")
	ErrInternal   = stderrors.New("internal error")
)

// ValidationError reports a field-level validation failure.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Unwrap links the error to ErrValidation.
func (e ValidationError) Unwrap() error { return ErrValidation }

// NotFoundError reports a missing entity.
type NotFoundError struct {
	Entity string
	ID     interface{}
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Entity)
}

// Unwrap links the error to ErrNotFound.
func (e NotFoundError) Unwrap() error { return ErrNotFound }

// ConflictError reports a conflict with existing data.
type ConflictError struct {
	Field   string
	Message string
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("conflict: %s", e.Message)
}

// Unwrap links the error to ErrConflict.
func (e ConflictError) Unwrap() error { return ErrConflict }

// InternalError reports an unexpected failure that must not reach the UI.
type InternalError struct {
	Message string
}

func (e InternalError) Error() string {
	return fmt.Sprintf("internal error: %s", e.Message)
}

// Unwrap links the error to ErrInternal.
func (e InternalError) Unwrap() error { return ErrInternal }

// IsValidation reports whether err is a validation failure.
func IsValidation(err error) bool {
	var target ValidationError
	return stderrors.As(err, &target) || stderrors.Is(err, ErrValidation)
}

// IsNotFound reports whether err is a missing-record failure.
func IsNotFound(err error) bool {
	var target NotFoundError
	return stderrors.As(err, &target) || stderrors.Is(err, ErrNotFound)
}

// IsConflict reports whether err is a conflict with existing data.
func IsConflict(err error) bool {
	var target ConflictError
	return stderrors.As(err, &target) || stderrors.Is(err, ErrConflict)
}

// IsInternal reports whether err is an internal failure.
func IsInternal(err error) bool {
	var target InternalError
	return stderrors.As(err, &target) || stderrors.Is(err, ErrInternal)
}
