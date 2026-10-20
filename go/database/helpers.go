package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"

	domainerrors "cargo_management_app/go/errors"
)

// Domain error sentinels live in the shared errors package. They are
// re-exported here so the database layer and the binding layer classify
// failures with exactly the same rules.
var (
	ErrNotFound   = domainerrors.ErrNotFound
	ErrValidation = domainerrors.ErrValidation
	ErrConflict   = domainerrors.ErrConflict
	ErrInternal   = domainerrors.ErrInternal
)

var (
	vinPattern        = regexp.MustCompile(`^[A-HJ-NPR-Z0-9]{17}$`)
	licensePattern    = regexp.MustCompile(`^[A-Za-z0-9-]{4,30}$`)
	ticketPattern     = regexp.MustCompile(`^[A-Za-z0-9._/-]{3,64}$`)
	identifierPattern = regexp.MustCompile(`^[A-Za-z0-9._/-]+$`)
	permitPattern     = regexp.MustCompile(`^[A-Za-z0-9-._/]{2,64}$`)
	phonePattern      = regexp.MustCompile(`^[0-9()+\s.-]{7,20}$`)
)

// ValidationError is an alias of the shared domain error so callers can use
// errors.As to recover the offending field.
type ValidationError = domainerrors.ValidationError

// invalid builds a validation failure that carries both the sentinel and the
// typed error, so errors.Is and errors.As both work on the result.
func invalid(field, message string) error {
	return fmt.Errorf("%w: %w", ErrValidation, domainerrors.ValidationError{Field: field, Message: message})
}

// notFound builds a not-found failure that carries the entity name.
func notFound(entity string) error {
	return fmt.Errorf("%w: %w", ErrNotFound, domainerrors.NotFoundError{Entity: entity})
}

// internalError builds an internal failure. These never describe user input
// problems, so the binding layer hides their detail from the UI.
func internalError(message string) error {
	return fmt.Errorf("%w: %w", ErrInternal, domainerrors.InternalError{Message: message})
}

func required(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return invalid(field, "is required")
	}
	return nil
}

func maxLen(value, field string, max int) error {
	if len([]rune(strings.TrimSpace(value))) > max {
		return invalid(field, fmt.Sprintf("must be at most %d characters", max))
	}
	return nil
}

func positive(value float64, field string) error {
	if value < 0 {
		return invalid(field, "must be zero or greater")
	}
	return nil
}

func positiveFloat(value float64, field string) error {
	if value <= 0 {
		return invalid(field, "must be greater than zero")
	}
	return nil
}

func positiveID(value int, field string) error {
	if value <= 0 {
		return invalid(field, "must be greater than zero")
	}
	return nil
}

func email(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return invalid(field, "must be a valid email address")
	}
	return nil
}

func phone(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	if !phonePattern.MatchString(strings.TrimSpace(value)) {
		return invalid(field, "must be a valid phone number")
	}
	return nil
}

func date(value, field string, requiredDate bool) error {
	value = strings.TrimSpace(value)
	if value == "" && !requiredDate {
		return nil
	}
	if value == "" {
		return invalid(field, "is required")
	}
	if _, err := time.Parse("2006-01-02", value); err != nil {
		return invalid(field, "must be a valid date in YYYY-MM-DD format")
	}
	return nil
}

func dateRange(start, end, startField, endField string) error {
	if err := date(start, startField, true); err != nil {
		return err
	}
	if err := date(end, endField, true); err != nil {
		return err
	}
	startAt, err := time.Parse("2006-01-02", strings.TrimSpace(start))
	if err != nil {
		return invalid(startField, "must be a valid date in YYYY-MM-DD format")
	}
	endAt, err := time.Parse("2006-01-02", strings.TrimSpace(end))
	if err != nil {
		return invalid(endField, "must be a valid date in YYYY-MM-DD format")
	}
	if endAt.Before(startAt) {
		return invalid(endField, "must be greater than or equal to the start date")
	}
	return nil
}

func oneOf(value, field string, allowed ...string) error {
	for _, v := range allowed {
		if value == v {
			return nil
		}
	}
	return invalid(field, "has an unsupported value")
}

func validateVIN(vin string) error {
	value := strings.TrimSpace(vin)
	if value == "" {
		return invalid("vin", "is required")
	}
	if !vinPattern.MatchString(value) {
		return invalid("vin", "must be a valid 17-character VIN with letters and numbers only")
	}
	return nil
}

func validateLicenseIdentifier(value, field string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return invalid(field, "is required")
	}
	if !licensePattern.MatchString(trimmed) {
		return invalid(field, "must contain only letters, numbers, and hyphens")
	}
	return nil
}

func validateTicketIdentifier(value, field string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return invalid(field, "is required")
	}
	if !ticketPattern.MatchString(trimmed) {
		return invalid(field, "must be a valid ticket or reference identifier")
	}
	return nil
}

func validateIdentifier(value, field string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return invalid(field, "is required")
	}
	if !identifierPattern.MatchString(trimmed) {
		return invalid(field, "contains unsupported characters")
	}
	return nil
}

func validateOptionalIdentifier(value, field string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	if !permitPattern.MatchString(trimmed) {
		return invalid(field, "contains unsupported characters")
	}
	return nil
}

func validateDriver(v DriverDTO) error {
	if err := required(v.Name, "name"); err != nil {
		return err
	}
	if err := required(v.LicenseNumber, "license_number"); err != nil {
		return err
	}
	if err := maxLen(v.Name, "name", 200); err != nil {
		return err
	}
	if err := validateLicenseIdentifier(v.LicenseNumber, "license_number"); err != nil {
		return err
	}
	if err := phone(v.Phone, "phone"); err != nil {
		return err
	}
	if err := email(v.Email, "email"); err != nil {
		return err
	}
	if err := maxLen(v.Address, "address", 500); err != nil {
		return err
	}
	return nil
}

func validateLoader(v LoaderDTO) error {
	if err := required(v.Name, "name"); err != nil {
		return err
	}
	if err := maxLen(v.Name, "name", 200); err != nil {
		return err
	}
	if err := phone(v.Phone, "phone"); err != nil {
		return err
	}
	if err := email(v.Email, "email"); err != nil {
		return err
	}
	if err := maxLen(v.Address, "address", 500); err != nil {
		return err
	}
	if err := maxLen(v.Company, "company", 200); err != nil {
		return err
	}
	return nil
}

func validateVehicle(v VehicleDTO) error {
	if err := validateIdentifier(v.VecID, "vec_id"); err != nil {
		return err
	}
	if err := validateLicenseIdentifier(v.LicenseID, "license_id"); err != nil {
		return err
	}
	if err := validateVIN(v.VIN); err != nil {
		return err
	}
	if err := date(v.LicenseExpiration, "license_expiration", true); err != nil {
		return err
	}
	if err := positiveID(v.DriverID, "driver_id"); err != nil {
		return err
	}
	if err := validateOptionalIdentifier(v.OverweightPermitID, "overweight_permit_id"); err != nil {
		return err
	}
	if v.TrailerNumber < 0 {
		return invalid("trailer_number", "must be zero or greater")
	}
	return nil
}

func validateWell(v WellDTO) error {
	if err := required(v.Name, "name"); err != nil {
		return err
	}
	if err := required(v.Location, "location"); err != nil {
		return err
	}
	if err := maxLen(v.Name, "name", 200); err != nil {
		return err
	}
	if err := maxLen(v.Client, "client", 200); err != nil {
		return err
	}
	if err := maxLen(v.Location, "location", 500); err != nil {
		return err
	}
	return nil
}

func validateLoad(v LoadDTO) error {
	if err := date(v.Date, "date", true); err != nil {
		return err
	}
	if err := required(v.LoadNumber, "load_number"); err != nil {
		return err
	}
	if err := validateTicketIdentifier(v.LoadNumber, "load_number"); err != nil {
		return err
	}
	if err := positiveID(v.LoaderID, "loader_id"); err != nil {
		return err
	}
	if err := positiveID(v.WellID, "well_id"); err != nil {
		return err
	}
	if err := required(v.TicketNumber, "ticket_number"); err != nil {
		return err
	}
	if err := validateTicketIdentifier(v.TicketNumber, "ticket_number"); err != nil {
		return err
	}
	if err := positiveFloat(v.Miles, "miles"); err != nil {
		return err
	}
	if err := positiveID(v.DriverID, "driver_id"); err != nil {
		return err
	}
	if err := validateIdentifier(v.VecID, "vec_id"); err != nil {
		return err
	}
	if v.TrailerNumber < 0 {
		return invalid("trailer_number", "must be zero or greater")
	}
	if err := positiveFloat(v.NetWeight, "net_weight"); err != nil {
		return err
	}
	if err := positiveFloat(v.Tons, "tons"); err != nil {
		return err
	}
	if err := positiveFloat(v.TonRate, "ton_rate"); err != nil {
		return err
	}
	if err := positiveFloat(v.TotalValue, "total_value"); err != nil {
		return err
	}
	if err := validateOptionalIdentifier(v.TicketID, "ticket_id"); err != nil {
		return err
	}
	return nil
}

func validateCategory(v CategoryDTO) error {
	if err := required(v.Type, "type"); err != nil {
		return err
	}
	return nil
}

func validateExpense(v ExpenseDTO) error {
	if err := required(v.Date, "date"); err != nil {
		return err
	}
	if err := date(v.Date, "date", true); err != nil {
		return err
	}
	if err := required(v.Concept, "concept"); err != nil {
		return err
	}
	if err := positiveID(v.CategoryID, "category_id"); err != nil {
		return err
	}
	if err := positive(v.Amount, "amount"); err != nil {
		return err
	}
	if err := positive(v.Quantity, "quantity"); err != nil {
		return err
	}
	if err := positive(v.Total, "total"); err != nil {
		return err
	}
	return nil
}

func validateIncome(v IncomeDTO) error {
	if err := required(v.Date, "date"); err != nil {
		return err
	}
	if err := date(v.Date, "date", true); err != nil {
		return err
	}
	if err := required(v.Concept, "concept"); err != nil {
		return err
	}
	if err := positiveID(v.CategoryID, "category_id"); err != nil {
		return err
	}
	if err := positive(v.Amount, "amount"); err != nil {
		return err
	}
	if err := positive(v.Quantity, "quantity"); err != nil {
		return err
	}
	if err := positive(v.Total, "total"); err != nil {
		return err
	}
	return nil
}

func validateFuelExpense(v FuelExpenseDTO) error {
	if err := positiveID(v.ExpenseID, "id_expense"); err != nil {
		return err
	}
	if err := positiveID(v.DriverID, "id_driver"); err != nil {
		return err
	}
	return nil
}

func validateStats(v StatsDTO) error {
	if v.NumberWeek <= 0 {
		return invalid("number_week", "must be greater than zero")
	}
	if err := required(v.DateRangeInit, "date_range_init"); err != nil {
		return err
	}
	if err := required(v.DateRangeEnd, "date_range_end"); err != nil {
		return err
	}
	if err := dateRange(v.DateRangeInit, v.DateRangeEnd, "date_range_init", "date_range_end"); err != nil {
		return err
	}
	if err := positive(v.ExpenseTotal, "expense_total"); err != nil {
		return err
	}
	if err := positive(v.IncomeTotal, "income_total"); err != nil {
		return err
	}
	if err := positive(v.NetTotal, "net_total"); err != nil {
		return err
	}
	return nil
}

func affected(result sql.Result, entity string) error {
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s update: %w", entity, err)
	}
	if n == 0 {
		return notFound(entity)
	}
	return nil
}

// databaseError translates a raw SQLite failure into a domain error so the
// frontend never receives driver or constraint internals.
func databaseError(operation string, err error) error {
	if err == nil {
		return nil
	}
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "unique constraint") || (strings.Contains(message, "constraint failed") && strings.Contains(message, "unique")) {
		return fmt.Errorf("%w during %s: %w", ErrConflict, operation, domainerrors.ConflictError{Message: "a record with the same unique value already exists"})
	}
	if strings.Contains(message, "foreign key constraint") {
		return fmt.Errorf("%w during %s: %w", ErrValidation, operation, domainerrors.ValidationError{Field: "reference", Message: "referenced record does not exist"})
	}
	return fmt.Errorf("%s: %w", operation, domainerrors.InternalError{Message: "database operation failed"})
}

// deleteError translates a failed delete into a domain error. A delete blocked
// by a foreign key is a conflict: the record is still referenced, so the user
// must remove the dependent records first.
func deleteError(operation, entity string, err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(strings.ToLower(err.Error()), "foreign key constraint") {
		return fmt.Errorf("%w during %s: %w", ErrConflict, operation, domainerrors.ConflictError{Field: entity, Message: "cannot delete " + entity + " because it is still referenced by other records"})
	}
	return databaseError(operation, err)
}

func scanDriver(row interface{ Scan(...any) error }) (*Driver, error) {
	var v Driver
	var meta sql.NullString
	err := row.Scan(&v.ID, &v.Name, &v.LicenseNumber, &v.Phone, &v.Email, &v.Address, &meta, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, err
	}
	v.Meta = rawMeta(meta)
	return &v, nil
}

// rawMeta converts a nullable text column into a JSON raw message.
// The meta column is optional, so a NULL value becomes an empty payload.
func rawMeta(value sql.NullString) json.RawMessage {
	if !value.Valid || strings.TrimSpace(value.String) == "" {
		return nil
	}
	return json.RawMessage(value.String)
}

func scanLoader(row interface{ Scan(...any) error }) (*Loader, error) {
	var v Loader
	err := row.Scan(&v.ID, &v.Name, &v.Phone, &v.Email, &v.Address, &v.Company, &v.CreatedAt, &v.UpdatedAt)
	return &v, err
}

func scanVehicle(row interface{ Scan(...any) error }) (*Vehicle, error) {
	var v Vehicle
	err := row.Scan(&v.ID, &v.VecID, &v.LicenseID, &v.TrailerNumber, &v.VIN, &v.LicenseExpiration, &v.DriverID, &v.OverweightPermitID, &v.CreatedAt, &v.UpdatedAt)
	return &v, err
}

func scanWell(row interface{ Scan(...any) error }) (*Well, error) {
	var v Well
	err := row.Scan(&v.ID, &v.Name, &v.Client, &v.Location, &v.CreatedAt, &v.UpdatedAt)
	return &v, err
}

func scanLoad(row interface{ Scan(...any) error }) (*Load, error) {
	var v Load
	err := row.Scan(&v.ID, &v.Date, &v.LoadNumber, &v.LoaderID, &v.WellID, &v.TicketNumber, &v.Miles, &v.DriverID, &v.VecID, &v.TrailerNumber, &v.NetWeight, &v.Tons, &v.TonRate, &v.TotalValue, &v.Status, &v.TicketID, &v.CreatedAt, &v.UpdatedAt)
	return &v, err
}

func scanCategory(row interface{ Scan(...any) error }) (*Category, error) {
	var v Category
	err := row.Scan(&v.ID, &v.Type, &v.CreatedAt, &v.UpdatedAt)
	return &v, err
}

func scanExpense(row interface{ Scan(...any) error }) (*Expense, error) {
	var v Expense
	var meta sql.NullString
	err := row.Scan(&v.ID, &v.Date, &v.Concept, &v.CategoryID, &v.Amount, &v.Quantity, &v.Total, &v.ReceiptID, &meta, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, err
	}
	v.Meta = rawMeta(meta)
	return &v, nil
}

func scanIncome(row interface{ Scan(...any) error }) (*Income, error) {
	var v Income
	var meta sql.NullString
	err := row.Scan(&v.ID, &v.Date, &v.Concept, &v.CategoryID, &v.Amount, &v.Quantity, &v.Total, &v.ReceiptID, &meta, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, err
	}
	v.Meta = rawMeta(meta)
	return &v, nil
}

func scanFuelExpense(row interface{ Scan(...any) error }) (*FuelExpense, error) {
	var v FuelExpense
	err := row.Scan(&v.ID, &v.ExpenseID, &v.DriverID, &v.CreatedAt, &v.UpdatedAt)
	return &v, err
}

func scanStats(row interface{ Scan(...any) error }) (*Stats, error) {
	var v Stats
	err := row.Scan(&v.ID, &v.NumberWeek, &v.DateRangeInit, &v.DateRangeEnd, &v.ExpenseTotal, &v.IncomeTotal, &v.NetTotal, &v.CreatedAt, &v.UpdatedAt)
	return &v, err
}

func (d *Database) validateVehicleReferences(v VehicleDTO) error {
	if err := positiveID(v.DriverID, "driver_id"); err != nil {
		return err
	}
	var exists int
	if err := d.runner.QueryRow(`SELECT 1 FROM driver WHERE id = ? LIMIT 1`, v.DriverID).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%w: %w", ErrValidation, domainerrors.ValidationError{Field: "driver_id", Message: fmt.Sprintf("%d does not reference an existing driver", v.DriverID)})
		}
		return fmt.Errorf("check driver reference: %w", err)
	}
	return nil
}

func (d *Database) validateLoadReferences(v LoadDTO) error {
	if err := positiveID(v.LoaderID, "loader_id"); err != nil {
		return err
	}
	if err := positiveID(v.WellID, "well_id"); err != nil {
		return err
	}
	if err := positiveID(v.DriverID, "driver_id"); err != nil {
		return err
	}

	if err := d.validateIDExists("loader", v.LoaderID, "loader_id"); err != nil {
		return err
	}
	if err := d.validateIDExists("well", v.WellID, "well_id"); err != nil {
		return err
	}
	if err := d.validateIDExists("driver", v.DriverID, "driver_id"); err != nil {
		return err
	}
	if strings.TrimSpace(v.VecID) != "" {
		var exists int
		if err := d.runner.QueryRow(`SELECT 1 FROM vehicle WHERE vec_id = ? LIMIT 1`, strings.TrimSpace(v.VecID)).Scan(&exists); err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("%w: %w", ErrValidation, domainerrors.ValidationError{Field: "vec_id", Message: fmt.Sprintf("%q does not reference an existing vehicle", v.VecID)})
			}
			return fmt.Errorf("check vehicle reference: %w", err)
		}
	}
	return nil
}

func (d *Database) validateIDExists(table string, id int, field string) error {
	var exists int
	if err := d.runner.QueryRow(fmt.Sprintf(`SELECT 1 FROM %s WHERE id = ? LIMIT 1`, table), id).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%w: %w", ErrValidation, domainerrors.ValidationError{Field: field, Message: fmt.Sprintf("%d does not reference an existing %s", id, table)})
		}
		return fmt.Errorf("check %s reference: %w", field, err)
	}
	return nil
}
