package database

import (
	"context"
	"testing"
)

func setupTestDB2(t *testing.T) *Database {
	t.Helper()
	path := "/tmp/cargo_test_load.db"
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

func TestCreateLoad_Success(t *testing.T) {
	t.Helper()
	db := setupTestDB2(t)
	defer db.Close()
	d := NewDatabase(db)

	// Ensure driver ID 1 exists by creating it
	dto := DriverDTO{Name: "TestDriver", LicenseNumber: "DL001", Phone: "5551111111", Email: "test@test.com", Address: "1 Test St"}
	_, err := d.CreateDriver(context.Background(), dto)
	if err != nil {
		t.Fatalf("create driver should succeed: %v", err)
	}

	dto2 := LoadDTO{Date: "2024-01-15", LoadNumber: "LDR-010", LoaderID: 1, WellID: 1, TicketNumber: "TK-100", Miles: 150, DriverID: 1, VecID: "VIN999TT0123", TrailerNumber: 2, NetWeight: 200, Tons: 100, TonRate: 25, TotalValue: 2500, Status: true, TicketID: "TICKET-001"}
	load, err := d.CreateLoad(context.Background(), dto2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if load.LoadNumber != "LDR-010" {
		t.Errorf("expected LDR-010, got %s", load.LoadNumber)
	}
}

func TestCreateLoad_InvalidLoaderID(t *testing.T) {
	t.Helper()
	db := setupTestDB2(t)
	defer db.Close()
	d := NewDatabase(db)

	// Create a driver first
	dto := DriverDTO{Name: "TestDriver", LicenseNumber: "DL001", Phone: "5551111111", Email: "test@test.com", Address: "1 Test St"}
	_, err := d.CreateDriver(context.Background(), dto)
	if err != nil {
		t.Fatalf("create driver should succeed: %v", err)
	}

	// LoaderID=0 should fail validation
	dto := LoadDTO{Date: "2024-01-15", LoadNumber: "LDR-011", LoaderID: 0, WellID: 1, TicketNumber: "TK-101", Miles: 50, DriverID: 1, VecID: "VIN999TT0123", TrailerNumber: 0, NetWeight: 50, Tons: 25, TonRate: 10, TotalValue: 250, Status: true, TicketID: "TK-001"}
	_, err = d.CreateLoad(context.Background(), dto)
	if err == nil {
		t.Error("expected validation error for invalid loader_id=0")
	}
}

func TestCreateLoad_FKConstraint_Driver(t *testing.T) {
	t.Helper()
	db := setupTestDB2(t)
	defer db.Close()
	d := NewDatabase(db)

	dto := LoadDTO{Date: "2024-01-15", LoadNumber: "LDR-012", LoaderID: 1, WellID: 1, TicketNumber: "TK-102", Miles: 50, DriverID: 999, VecID: "VIN999TT0123", TrailerNumber: 0, NetWeight: 50, Tons: 25, TonRate: 10, TotalValue: 250, Status: true, TicketID: "TK-001"}
	_, err := d.CreateLoad(context.Background(), dto)
	if err == nil {
		t.Error("expected FK constraint error for invalid driver_id=999")
	}
}