package database

import (
	"database/sql"
	"testing"
)

func setupTestDB(t *testing.T) *Database {
	t.Helper()
	path := "/tmp/cargo_test_" + string(rune(t.TimetNano()%1000)) + ".db"
	db, err := NewDatabase(path)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	err = db.createSchema()
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}
	return db
}

func TestCreateDriver_Success(t *testing.T) {
	t.Helper()
	db := setupTestDB(t)
	defer db.Close()
	d := NewDatabase(db)

	dto := DriverDTO{Name: "John Doe", LicenseNumber: "DL123", Phone: "5551234567", Email: "john@example.com", Address: "123 Main St"}
	driver, err := d.CreateDriver(context.Background(), dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if driver.Name != "John Doe" {
		t.Errorf("expected John Doe, got %s", driver.Name)
	}
}

func TestCreateDriver_DuplicateLicense(t *testing.T) {
	t.Helper()
	db := setupTestDB(t)
	defer db.Close()
	d := NewDatabase(db)

	dto1 := DriverDTO{Name: "John Doe", LicenseNumber: "DL123", Phone: "5551234567", Email: "john@example.com", Address: "123 Main St"}
	_, err := d.CreateDriver(context.Background(), dto1)
	if err != nil {
		t.Fatalf("first create should succeed: %v", err)
	}

	dto2 := DriverDTO{Name: "Jane Doe", LicenseNumber: "DL123", Phone: "5559876543", Email: "jane@example.com", Address: "456 Oak St"}
	_, err = d.CreateDriver(context.Background(), dto2)
	if err == nil {
		t.Error("expected error for duplicate license number")
	}
}

func TestGetDriver_NotFound(t *testing.T) {
	t.Helper()
	db := setupTestDB(t)
	defer db.Close()
	d := NewDatabase(db)

	_, err := d.GetDriver(999)
	if err == nil {
		t.Error("expected not found error")
	}
}

func TestGetDriver_Success(t *testing.T) {
	t.Helper()
	db := setupTestDB(t)
	defer db.Close()
	d := NewDatabase(db)

	dto := DriverDTO{Name: "Jane Doe", LicenseNumber: "DL456", Phone: "5554567890", Email: "jane@example.com", Address: "789 Elm St"}
	driver, err := d.CreateDriver(context.Background(), dto)
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
	db := setupTestDB(t)
	defer db.Close()
	d := NewDatabase(db)

	dto := DriverDTO{Name: "Original", LicenseNumber: "DL789", Phone: "5551111111", Email: "original@example.com", Address: "111 Original St"}
	driver, err := d.CreateDriver(context.Background(), dto)
	if err != nil {
		t.Fatalf("create should succeed: %v", err)
	}

	updateDTO := DriverDTO{Name: "Updated", LicenseNumber: driver.LicenseNumber, Phone: driver.Phone, Email: driver.Email, Address: driver.Address}
	err = d.UpdateDriver(driver.ID, updateDTO)
	if err != nil {
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
	db := setupTestDB(t)
	defer db.Close()
	d := NewDatabase(db)

	dto := DriverDTO{Name: "ToBeDeleted", LicenseNumber: "DL999", Phone: "5551231234", Email: "delete@example.com", Address: "999 Delete St"}
	driver, err := d.CreateDriver(context.Background(), dto)
	if err != nil {
		t.Fatalf("create should succeed: %v", err)
	}

	err = d.DeleteDriver(driver.ID)
	if err != nil {
		t.Fatalf("delete should succeed: %v", err)
	}

	_, err = d.GetDriver(driver.ID)
	if err == nil {
		t.Error("expected not found after delete")
	}
}

func TestDeleteDriver_Referenced(t *testing.T) {
	t.Helper()
	db := setupTestDB(t)
	defer db.Close()
	d := NewDatabase(db)

	dto := DriverDTO{Name: "DriverRef", LicenseNumber: "DL888", Phone: "5551111111", Email: "ref@example.com", Address: "888 Ref St"}
	driver, err := d.CreateDriver(context.Background(), dto)
	if err != nil {
		t.Fatalf("create should succeed: %v", err)
	}

	// Create a load referencing this driver
	loadDTO := LoadDTO{Date: "2024-01-01", LoadNumber: "LDR-001", LoaderID: 1, WellID: 1, TicketNumber: "TK-001", Miles: 100, DriverID: driver.ID, VecID: "VIN123", TrailerNumber: 0, NetWeight: 100, Tons: 50, TonRate: 10, TotalValue: 500, Status: true, TicketID: "TK-001"}
	_, err = d.CreateLoad(context.Background(), loadDTO)
	if err != nil {
		t.Fatalf("create load should succeed: %v", err)
	}

	// Try to delete driver - should fail due to FK constraint
	err = d.DeleteDriver(driver.ID)
	if err == nil {
		t.Error("expected FK constraint error when deleting driver with assigned loads")
	}
}