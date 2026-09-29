package database

import (
    "database/sql"
    "errors"
    "fmt"
    "net/mail"
    "regexp"
    "strings"
    "time"
)

var (
    ErrNotFound = errors.New("record not found")
    ErrValidation = errors.New("validation failed")
    ErrConflict = errors.New("record conflicts with existing data")
)

type ValidationError struct { Field string; Message string }
func (e ValidationError) Error() string { return fmt.Sprintf("%s: %s", e.Field, e.Message) }

func invalid(field, message string) error { return fmt.Errorf("%w: %s", ErrValidation, ValidationError{field, message}) }
func required(value, field string) error { if strings.TrimSpace(value)=="" { return invalid(field,"is required") }; return nil }
func maxLen(value, field string, max int) error { if len([]rune(strings.TrimSpace(value)))>max { return invalid(field, fmt.Sprintf("must be at most %d characters", max)) }; return nil }
func positive(value float64, field string) error { if value < 0 { return invalid(field, "must be zero or greater") }; return nil }
func positiveID(value int, field string) error { if value <= 0 { return invalid(field, "must be greater than zero") }; return nil }
func email(value, field string) error { if strings.TrimSpace(value)=="" { return nil }; if _, err := mail.ParseAddress(value); err != nil { return invalid(field, "must be a valid email address") }; return nil }
func date(value, field string, requiredDate bool) error {
    value = strings.TrimSpace(value)
    if value == "" && !requiredDate { return nil }
    if value == "" { return invalid(field, "is required") }
    if _, err := time.Parse("2006-01-02", value); err != nil {
        return invalid(field, "must be a valid date in YYYY-MM-DD format")
    }
    return nil
}
func oneOf(value, field string, allowed ...string) error {
    for _, v := range allowed { if value == v { return nil } }
    return invalid(field, "has an unsupported value")
}

func validateDriver(v DriverDTO) error {
    if err := required(v.Name, "name"); err != nil { return err }
    if err := required(v.LicenseNumber, "license_number"); err != nil { return err }
    if err := maxLen(v.Name, "name", 200); err != nil { return err }
    if err := maxLen(v.Phone, "phone", 50); err != nil { return err }
    if err := email(v.Email, "email"); err != nil { return err }
    if err := maxLen(v.Address, "address", 500); err != nil { return err }
    return nil
}

func validateTruck(v TruckDTO) error {
    if err := required(v.Name, "name"); err != nil { return err }
    if err := required(v.VIN, "vin"); err != nil { return err }
    if err := required(v.LicenseExpiration, "license_expiration"); err != nil { return err }
    if err := date(v.LicenseExpiration, "license_expiration", true); err != nil { return err }
    if err := positiveID(v.DriverID, "driver_id"); err != nil { return err }
    return nil
}

func validateTrailer(v TrailerDTO) error {
    if err := required(v.Name, "name"); err != nil { return err }
    if v.Number < 0 { return invalid("number", "must be zero or greater") }
    return nil
}

func validateWell(v WellDTO) error {
    if err := required(v.Name, "name"); err != nil { return err }
    if err := required(v.Location, "location"); err != nil { return err }
    return nil
}

func validateTrip(v TripDTO) error {
    if err := required(v.Date, "date"); err != nil { return err }
    if err := date(v.Date, "date", true); err != nil { return err }
    if err := required(v.LoadNumber, "load_number"); err != nil { return err }
    if err := required(v.TicketNumber, "ticket_number"); err != nil { return err }
    if err := positiveID(v.WellID, "well_id"); err != nil { return err }
    if err := positiveID(v.DriverID, "driver_id"); err != nil { return err }
    if err := required(v.TruckID, "truck_id"); err != nil { return err }
    if err := positiveID(v.TrailerID, "id_trailer"); err != nil { return err }
    return nil
}

func validateCategory(v CategoryDTO) error {
    if err := required(v.Type, "type"); err != nil { return err }
    return nil
}

func validateExpense(v ExpenseDTO) error {
    if err := required(v.Date, "date"); err != nil { return err }
    if err := date(v.Date, "date", true); err != nil { return err }
    if err := required(v.Concept, "concept"); err != nil { return err }
    if err := positiveID(v.CategoryID, "category_id"); err != nil { return err }
    if err := positive(v.Amount, "amount"); err != nil { return err }
    if err := positive(v.Quantity, "quantity"); err != nil { return err }
    if err := positive(v.Total, "total"); err != nil { return err }
    return nil
}

func validateIncome(v IncomeDTO) error {
    if err := required(v.Date, "date"); err != nil { return err }
    if err := date(v.Date, "date", true); err != nil { return err }
    if err := required(v.Concept, "concept"); err != nil { return err }
    if err := positiveID(v.CategoryID, "category_id"); err != nil { return err }
    if err := positive(v.Amount, "amount"); err != nil { return err }
    if err := positive(v.Quantity, "quantity"); err != nil { return err }
    if err := positive(v.Total, "total"); err != nil { return err }
    return nil
}

func validateFuelExpense(v FuelExpenseDTO) error {
    if err := positiveID(v.ExpenseID, "id_expense"); err != nil { return err }
    if err := positiveID(v.DriverID, "id_driver"); err != nil { return err }
    return nil
}

func validateStats(v StatsDTO) error {
    if v.NumberWeek <= 0 { return invalid("number_week", "must be greater than zero") }
    if err := required(v.DateRangeInit, "date_range_init"); err != nil { return err }
    if err := required(v.DateRangeEnd, "date_range_end"); err != nil { return err }
    if err := date(v.DateRangeInit, "date_range_init", true); err != nil { return err }
    if err := date(v.DateRangeEnd, "date_range_end", true); err != nil { return err }
    if err := positive(v.ExpenseTotal, "expense_total"); err != nil { return err }
    if err := positive(v.IncomeTotal, "income_total"); err != nil { return err }
    if err := positive(v.NetTotal, "net_total"); err != nil { return err }
    return nil
}

func affected(result sql.Result, entity string) error {
    n, err := result.RowsAffected()
    if err != nil { return fmt.Errorf("%s update: %w", entity, err) }
    if n == 0 { return fmt.Errorf("%w: %s", ErrNotFound, entity) }
    return nil
}

func databaseError(operation string, err error) error {
    if err == nil { return nil }
    message := strings.ToLower(err.Error())
    if strings.Contains(message, "unique constraint") || (strings.Contains(message, "constraint failed") && strings.Contains(message, "unique")) {
        return fmt.Errorf("%w during %s: %v", ErrConflict, operation, err)
    }
    if strings.Contains(message, "foreign key constraint") {
        return fmt.Errorf("%w during %s: referenced record does not exist", ErrValidation, operation)
    }
    return fmt.Errorf("%s: %w", operation, err)
}

func scanDriver(row interface{ Scan(...any) error }) (*Driver, error) {
    var v Driver
    err := row.Scan(&v.ID, &v.Name, &v.LicenseNumber, &v.Phone, &v.Email, &v.Address, &v.Meta, &v.CreatedAt, &v.UpdatedAt)
    return &v, err
}

func scanTruck(row interface{ Scan(...any) error }) (*Truck, error) {
    var v Truck
    err := row.Scan(&v.ID, &v.Name, &v.VIN, &v.LicenseExpiration, &v.DriverID, &v.OverweightPermitID, &v.CreatedAt, &v.UpdatedAt)
    return &v, err
}

func scanTrailer(row interface{ Scan(...any) error }) (*Trailer, error) {
    var v Trailer
    err := row.Scan(&v.ID, &v.Name, &v.Number, &v.CreatedAt, &v.UpdatedAt)
    return &v, err
}

func scanWell(row interface{ Scan(...any) error }) (*Well, error) {
    var v Well
    err := row.Scan(&v.ID, &v.Name, &v.Location, &v.CreatedAt, &v.UpdatedAt)
    return &v, err
}

func scanTrip(row interface{ Scan(...any) error }) (*Trip, error) {
    var v Trip
    err := row.Scan(&v.ID, &v.Date, &v.LoadNumber, &v.WellID, &v.TicketNumber, &v.DriverID, &v.TruckID, &v.TrailerID, &v.Meta, &v.Status, &v.TicketID, &v.CreatedAt, &v.UpdatedAt)
    return &v, err
}

func scanCategory(row interface{ Scan(...any) error }) (*Category, error) {
    var v Category
    err := row.Scan(&v.ID, &v.Type, &v.CreatedAt, &v.UpdatedAt)
    return &v, err
}

func scanExpense(row interface{ Scan(...any) error }) (*Expense, error) {
    var v Expense
    err := row.Scan(&v.ID, &v.Date, &v.Concept, &v.CategoryID, &v.Amount, &v.Quantity, &v.Total, &v.ReceiptID, &v.Meta, &v.CreatedAt, &v.UpdatedAt)
    return &v, err
}

func scanIncome(row interface{ Scan(...any) error }) (*Income, error) {
    var v Income
    err := row.Scan(&v.ID, &v.Date, &v.Concept, &v.CategoryID, &v.Amount, &v.Quantity, &v.Total, &v.ReceiptID, &v.Meta, &v.CreatedAt, &v.UpdatedAt)
    return &v, err
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

var identifierPattern = regexp.MustCompile(`^[A-Za-z0-9._/-]+$`)
