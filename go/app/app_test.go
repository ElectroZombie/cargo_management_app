package app

import (
	"path/filepath"
	"strings"
	"testing"

	"cargo_management_app/go/database"
	domainerrors "cargo_management_app/go/errors"
)

// setupTestApp creates an isolated database and a bound application for a test.
func setupTestApp(t *testing.T) *App {
	t.Helper()

	db, err := database.NewDatabase(filepath.Join(t.TempDir(), "cargo_app_test.db"))
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	return NewApp(db)
}

func TestCreateDriver_Success(t *testing.T) {
	t.Helper()
	a := setupTestApp(t)

	response := a.CreateDriver(database.DriverDTO{Name: "John Doe", LicenseNumber: "DL-2001", Phone: "5551234567", Email: "john@example.com", Address: "123 Main St"})
	if !response.Success {
		t.Fatalf("expected success, got error: %s", response.Error)
	}
	if response.Error != "" {
		t.Errorf("expected no error message, got %s", response.Error)
	}
	if response.Code != "" {
		t.Errorf("expected no error code, got %s", response.Code)
	}

	driver, ok := response.Data.(*database.Driver)
	if !ok {
		t.Fatalf("expected a driver payload, got %T", response.Data)
	}
	if driver.Name != "John Doe" {
		t.Errorf("expected John Doe, got %s", driver.Name)
	}
	if driver.ID == 0 {
		t.Error("expected the created driver to have an ID")
	}
}

func TestCreateDriver_ValidationError(t *testing.T) {
	t.Helper()
	a := setupTestApp(t)

	response := a.CreateDriver(database.DriverDTO{})
	if response.Success {
		t.Fatal("expected failure for an empty driver")
	}
	if response.Code != ErrorCodeValidation {
		t.Errorf("expected code %q, got %q", ErrorCodeValidation, response.Code)
	}
	if response.Error == "" {
		t.Error("expected a non-empty error message")
	}
}

func TestCreateDriver_ConflictError(t *testing.T) {
	t.Helper()
	a := setupTestApp(t)

	first := a.CreateDriver(database.DriverDTO{Name: "John Doe", LicenseNumber: "DL-DUP"})
	if !first.Success {
		t.Fatalf("first create should succeed: %s", first.Error)
	}

	second := a.CreateDriver(database.DriverDTO{Name: "Jane Doe", LicenseNumber: "DL-DUP"})
	if second.Success {
		t.Fatal("expected failure for a duplicate license number")
	}
	if second.Code != ErrorCodeConflict {
		t.Errorf("expected code %q, got %q", ErrorCodeConflict, second.Code)
	}
}

func TestGetDriver_NotFound(t *testing.T) {
	t.Helper()
	a := setupTestApp(t)

	response := a.GetDriver(999)
	if response.Success {
		t.Fatal("expected failure for a missing driver")
	}
	if response.Code != ErrorCodeNotFound {
		t.Errorf("expected code %q, got %q", ErrorCodeNotFound, response.Code)
	}
	if !strings.Contains(response.Error, "not found") {
		t.Errorf("expected a not-found message, got %q", response.Error)
	}
}

func TestUpdateDriver_Success(t *testing.T) {
	t.Helper()
	a := setupTestApp(t)

	created := a.CreateDriver(database.DriverDTO{Name: "Original", LicenseNumber: "DL-UPD"})
	if !created.Success {
		t.Fatalf("create should succeed: %s", created.Error)
	}
	driver := created.Data.(*database.Driver)

	response := a.UpdateDriver(driver.ID, database.DriverDTO{Name: "Updated", LicenseNumber: driver.LicenseNumber})
	if !response.Success {
		t.Fatalf("expected success, got error: %s", response.Error)
	}
	updated := response.Data.(*database.Driver)
	if updated.Name != "Updated" {
		t.Errorf("expected Updated, got %s", updated.Name)
	}
}

func TestDeleteDriver_Success(t *testing.T) {
	t.Helper()
	a := setupTestApp(t)

	created := a.CreateDriver(database.DriverDTO{Name: "ToBeDeleted", LicenseNumber: "DL-DEL"})
	if !created.Success {
		t.Fatalf("create should succeed: %s", created.Error)
	}
	driver := created.Data.(*database.Driver)

	if response := a.DeleteDriver(driver.ID); !response.Success {
		t.Fatalf("expected success, got error: %s", response.Error)
	}
	if response := a.GetDriver(driver.ID); response.Success {
		t.Error("expected failure after delete")
	}
}

func TestListDrivers_ReturnsTotal(t *testing.T) {
	t.Helper()
	a := setupTestApp(t)

	for _, name := range []string{"One", "Two", "Three"} {
		if response := a.CreateDriver(database.DriverDTO{Name: name, LicenseNumber: "DL-" + name}); !response.Success {
			t.Fatalf("create should succeed: %s", response.Error)
		}
	}

	response := a.ListDrivers(database.DriverFilters{}, database.QueryOptions{Page: 1, PageSize: 2})
	if !response.Success {
		t.Fatalf("expected success, got error: %s", response.Error)
	}
	if response.Total != 3 {
		t.Errorf("expected total 3, got %d", response.Total)
	}

	drivers, ok := response.Data.([]database.Driver)
	if !ok {
		t.Fatalf("expected a driver slice, got %T", response.Data)
	}
	if len(drivers) != 2 {
		t.Errorf("expected 2 drivers on the first page, got %d", len(drivers))
	}
}

func TestListDrivers_Filter(t *testing.T) {
	t.Helper()
	a := setupTestApp(t)

	a.CreateDriver(database.DriverDTO{Name: "Alpha One", LicenseNumber: "DL-A"})
	a.CreateDriver(database.DriverDTO{Name: "Beta Two", LicenseNumber: "DL-B"})

	response := a.ListDrivers(database.DriverFilters{Name: "Alpha"}, database.QueryOptions{Page: 1, PageSize: 25})
	if !response.Success {
		t.Fatalf("expected success, got error: %s", response.Error)
	}
	if response.Total != 1 {
		t.Errorf("expected total 1, got %d", response.Total)
	}
}

func TestGetDashboardStats_Success(t *testing.T) {
	t.Helper()
	a := setupTestApp(t)

	a.CreateDriver(database.DriverDTO{Name: "Stats Driver", LicenseNumber: "DL-STATS"})

	response := a.GetDashboardStats()
	if !response.Success {
		t.Fatalf("expected success, got error: %s", response.Error)
	}
	stats, ok := response.Data.(*database.DashboardStats)
	if !ok {
		t.Fatalf("expected dashboard stats, got %T", response.Data)
	}
	if stats.DriverCount != 1 {
		t.Errorf("expected 1 driver, got %d", stats.DriverCount)
	}
}

func TestFormatError_InternalErrorHidesDetail(t *testing.T) {
	t.Helper()

	raw := domainerrors.InternalError{Message: "sqlite3: no such table: secret_table"}
	response := failure(raw)

	if response.Success {
		t.Fatal("expected failure")
	}
	if response.Code != ErrorCodeInternal {
		t.Errorf("expected code %q, got %q", ErrorCodeInternal, response.Code)
	}
	if strings.Contains(response.Error, "secret_table") {
		t.Errorf("expected internal detail to be hidden, got %q", response.Error)
	}
	if response.Error != "An internal error occurred. Please try again." {
		t.Errorf("unexpected message: %q", response.Error)
	}
}

func TestFormatError_NilError(t *testing.T) {
	t.Helper()

	if got := formatError(nil); got != "" {
		t.Errorf("expected an empty message, got %q", got)
	}
	if got := errorCode(nil); got != "" {
		t.Errorf("expected an empty code, got %q", got)
	}
}
