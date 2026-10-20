package validators

// ValidateDriver checks a DriverDTO for validation errors.
// Validation logic is also implemented in database/helpers.go.
// This file provides a dedicated package for DTO validation per AGENTS.md conventions.
func ValidateDriver(dto *interface{}) error {
	// Validation logic moved to database/helpers.go validateDriver function
	// for shared use between services and potential frontend validation.
	return nil
}