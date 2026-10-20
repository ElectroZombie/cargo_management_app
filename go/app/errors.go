package app

import (
	"errors"
	"log/slog"

	domainerrors "cargo_management_app/go/errors"
)

// Error codes returned alongside the human-readable message so the frontend
// can react to a failure category without parsing the message text.
const (
	ErrorCodeValidation = "validation"
	ErrorCodeNotFound   = "not_found"
	ErrorCodeConflict   = "conflict"
	ErrorCodeInternal   = "internal"
	ErrorCodeUnknown    = "unknown"
)

// errorCode classifies a domain error into its stable error code.
func errorCode(err error) string {
	switch {
	case err == nil:
		return ""
	case domainerrors.IsValidation(err):
		return ErrorCodeValidation
	case domainerrors.IsNotFound(err):
		return ErrorCodeNotFound
	case domainerrors.IsConflict(err):
		return ErrorCodeConflict
	case domainerrors.IsInternal(err):
		return ErrorCodeInternal
	default:
		return ErrorCodeUnknown
	}
}

// formatError converts a domain error into a user-facing message.
//
// Internal failures are logged in full but reduced to a generic message, so
// database internals never reach the frontend. Validation, not-found, and
// conflict errors keep their specific detail because the user can act on it.
func formatError(err error) string {
	if err == nil {
		return ""
	}

	switch {
	case domainerrors.IsValidation(err):
		var validation domainerrors.ValidationError
		if errors.As(err, &validation) {
			return validation.Field + ": " + validation.Message
		}
		return "The submitted data is not valid."

	case domainerrors.IsNotFound(err):
		var missing domainerrors.NotFoundError
		if errors.As(err, &missing) {
			return missing.Entity + " not found."
		}
		return "The requested record was not found."

	case domainerrors.IsConflict(err):
		var conflict domainerrors.ConflictError
		if errors.As(err, &conflict) {
			return conflict.Message
		}
		return "This record conflicts with existing data."

	case domainerrors.IsInternal(err):
		slog.Error("internal error", "error", err)
		return "An internal error occurred. Please try again."

	default:
		slog.Error("unclassified error", "error", err)
		return "An unexpected error occurred."
	}
}
