package database

import (
	"path/filepath"
	"testing"

	domainerrors "cargo_management_app/go/errors"
)

// setupTestDB creates an isolated SQLite database in a temporary directory for a test.
func setupTestDB(t *testing.T) *Database {
	t.Helper()
	db, err := NewDatabase(filepath.Join(t.TempDir(), "cargo_test.db"))
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestCreateDriver_Success(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	dto := DriverDTO{Name: "John Doe", LicenseNumber: "DL123", Phone: "5551234567", Email: "john@example.com", Address: "123 Main St"}
	driver, err := d.CreateDriver(dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if driver.Name != "John Doe" {
		t.Errorf("expected John Doe, got %s", driver.Name)
	}
}

func TestCreateDriver_DuplicateLicense(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	dto1 := DriverDTO{Name: "John Doe", LicenseNumber: "DL123", Phone: "5551234567", Email: "john@example.com", Address: "123 Main St"}
	if _, err := d.CreateDriver(dto1); err != nil {
		t.Fatalf("first create should succeed: %v", err)
	}

	dto2 := DriverDTO{Name: "Jane Doe", LicenseNumber: "DL123", Phone: "5559876543", Email: "jane@example.com", Address: "456 Oak St"}
	if _, err := d.CreateDriver(dto2); err == nil {
		t.Error("expected error for duplicate license number")
	}
}

func TestGetDriver_NotFound(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	if _, err := d.GetDriver(999); err == nil {
		t.Error("expected not found error")
	}
}

func TestGetDriver_Success(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	dto := DriverDTO{Name: "Jane Doe", LicenseNumber: "DL456", Phone: "5554567890", Email: "jane@example.com", Address: "789 Elm St"}
	driver, err := d.CreateDriver(dto)
	if err != nil {
		t.Fatalf("create should succeed: %v", err)
	}

	found, err := d.GetDriver(driver.ID)
	if err != nil {
		t.Fatalf("get should succeed: %v", err)
	}
	if found.Name != "Jane Doe" {
		t.Errorf("expected Jane Doe, got %s", found.Name)
	}
}

func TestUpdateDriver_Success(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	dto := DriverDTO{Name: "Original", LicenseNumber: "DL789", Phone: "5551111111", Email: "original@example.com", Address: "111 Original St"}
	driver, err := d.CreateDriver(dto)
	if err != nil {
		t.Fatalf("create should succeed: %v", err)
	}

	updateDTO := DriverDTO{Name: "Updated", LicenseNumber: driver.LicenseNumber, Phone: driver.Phone, Email: driver.Email, Address: driver.Address}
	if err := d.UpdateDriver(driver.ID, updateDTO); err != nil {
		t.Fatalf("update should succeed: %v", err)
	}

	found, err := d.GetDriver(driver.ID)
	if err != nil {
		t.Fatalf("get should succeed: %v", err)
	}
	if found.Name != "Updated" {
		t.Errorf("expected Updated, got %s", found.Name)
	}
}

func TestDeleteDriver_Success(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	dto := DriverDTO{Name: "ToBeDeleted", LicenseNumber: "DL999", Phone: "5551231234", Email: "delete@example.com", Address: "999 Delete St"}
	driver, err := d.CreateDriver(dto)
	if err != nil {
		t.Fatalf("create should succeed: %v", err)
	}

	if err := d.DeleteDriver(driver.ID); err != nil {
		t.Fatalf("delete should succeed: %v", err)
	}

	if _, err := d.GetDriver(driver.ID); err == nil {
		t.Error("expected not found after delete")
	}
}

func TestDeleteDriver_Referenced(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	driver, err := d.CreateDriver(DriverDTO{Name: "DriverRef", LicenseNumber: "DL888", Phone: "5551111111", Email: "ref@example.com", Address: "888 Ref St"})
	if err != nil {
		t.Fatalf("create driver should succeed: %v", err)
	}
	loader, err := d.CreateLoader(LoaderDTO{Name: "LoaderRef", Company: "Ref Co"})
	if err != nil {
		t.Fatalf("create loader should succeed: %v", err)
	}
	well, err := d.CreateWell(WellDTO{Name: "WellRef", Location: "Ref Field"})
	if err != nil {
		t.Fatalf("create well should succeed: %v", err)
	}
	vehicle, err := d.CreateVehicle(VehicleDTO{
		VecID:             "VEC-REF",
		LicenseID:         "LIC-REF",
		VIN:               "1HGBH41JXMN109186",
		LicenseExpiration: "2027-01-31",
		DriverID:          driver.ID,
	})
	if err != nil {
		t.Fatalf("create vehicle should succeed: %v", err)
	}

	load := LoadDTO{
		Date:          "2026-01-01",
		LoadNumber:    "LDR-REF",
		LoaderID:      loader.ID,
		WellID:        well.ID,
		TicketNumber:  "TK-REF",
		Miles:         100,
		DriverID:      driver.ID,
		VecID:         vehicle.VecID,
		TrailerNumber: 0,
		NetWeight:     100,
		Tons:          50,
		TonRate:       10,
		TotalValue:    500,
		Status:        true,
		TicketID:      "TK-REF",
	}
	if _, err := d.CreateLoad(load); err != nil {
		t.Fatalf("create load should succeed: %v", err)
	}

	// Deleting a driver that still has assigned loads must fail on the FK constraint.
	err = d.DeleteDriver(driver.ID)
	if err == nil {
		t.Fatal("expected FK constraint error when deleting driver with assigned loads")
	}
	if !domainerrors.IsConflict(err) {
		t.Errorf("expected a conflict error, got %v", err)
	}
}
