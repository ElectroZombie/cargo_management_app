package database

import (
	"testing"

	domainerrors "cargo_management_app/go/errors"
)

// createVehicleDriver seeds the driver a vehicle must reference.
func createVehicleDriver(t *testing.T, d *Database) int {
	t.Helper()

	driver, err := d.CreateDriver(DriverDTO{Name: "Vehicle Driver", LicenseNumber: "DL-VEH", Phone: "5551111111", Email: "veh@test.com", Address: "1 Test St"})
	if err != nil {
		t.Fatalf("create driver should succeed: %v", err)
	}
	return driver.ID
}

func TestCreateVehicle_Success(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	driverID := createVehicleDriver(t, d)
	dto := VehicleDTO{
		VecID:             "VEC-9001",
		LicenseID:         "LIC-9001",
		TrailerNumber:     12,
		VIN:               "3T2BK1BA4KC123456",
		LicenseExpiration: "2028-03-31",
		DriverID:          driverID,
	}
	vehicle, err := d.CreateVehicle(dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vehicle.VecID != "VEC-9001" {
		t.Errorf("expected VEC-9001, got %s", vehicle.VecID)
	}
	if vehicle.DriverID != driverID {
		t.Errorf("expected driver %d, got %d", driverID, vehicle.DriverID)
	}
}

func TestCreateVehicle_InvalidVIN(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	driverID := createVehicleDriver(t, d)
	_, err := d.CreateVehicle(VehicleDTO{
		VecID:             "VEC-BADVIN",
		LicenseID:         "LIC-BADVIN",
		VIN:               "TOO-SHORT",
		LicenseExpiration: "2028-03-31",
		DriverID:          driverID,
	})
	if err == nil {
		t.Fatal("expected validation error for malformed VIN")
	}
	if !domainerrors.IsValidation(err) {
		t.Errorf("expected a validation error, got %v", err)
	}
}

func TestCreateVehicle_UnknownDriver(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	_, err := d.CreateVehicle(VehicleDTO{
		VecID:             "VEC-ORPHAN",
		LicenseID:         "LIC-ORPHAN",
		VIN:               "4T2BK1BA4KC123456",
		LicenseExpiration: "2028-03-31",
		DriverID:          999,
	})
	if err == nil {
		t.Fatal("expected reference error for unknown driver")
	}
	if !domainerrors.IsValidation(err) {
		t.Errorf("expected a validation error, got %v", err)
	}
}

func TestCreateVehicle_DuplicateVIN(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	driverID := createVehicleDriver(t, d)
	first := VehicleDTO{VecID: "VEC-DUP", LicenseID: "LIC-DUP", VIN: "5T2BK1BA4KC123456", LicenseExpiration: "2028-03-31", DriverID: driverID}
	if _, err := d.CreateVehicle(first); err != nil {
		t.Fatalf("first create should succeed: %v", err)
	}

	second := VehicleDTO{VecID: "VEC-DUP2", LicenseID: "LIC-DUP2", VIN: first.VIN, LicenseExpiration: "2028-03-31", DriverID: driverID}
	if _, err := d.CreateVehicle(second); err == nil {
		t.Error("expected conflict error for duplicate VIN")
	} else if !domainerrors.IsConflict(err) {
		t.Errorf("expected a conflict error, got %v", err)
	}
}

func TestGetVehicle_NotFound(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	_, err := d.GetVehicle(999)
	if err == nil {
		t.Fatal("expected not found error")
	}
	if !domainerrors.IsNotFound(err) {
		t.Errorf("expected a not found error, got %v", err)
	}
}

func TestUpdateVehicle_Success(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	driverID := createVehicleDriver(t, d)
	vehicle, err := d.CreateVehicle(VehicleDTO{VecID: "VEC-UPD", LicenseID: "LIC-UPD", VIN: "6T2BK1BA4KC123456", LicenseExpiration: "2028-03-31", DriverID: driverID})
	if err != nil {
		t.Fatalf("create should succeed: %v", err)
	}

	updateDTO := VehicleDTO{VecID: vehicle.VecID, LicenseID: "LIC-UPD2", TrailerNumber: vehicle.TrailerNumber, VIN: vehicle.VIN, LicenseExpiration: vehicle.LicenseExpiration, DriverID: driverID}
	if err := d.UpdateVehicle(vehicle.ID, updateDTO); err != nil {
		t.Fatalf("update should succeed: %v", err)
	}

	found, err := d.GetVehicle(vehicle.ID)
	if err != nil {
		t.Fatalf("get should succeed: %v", err)
	}
	if found.LicenseID != "LIC-UPD2" {
		t.Errorf("expected LIC-UPD2, got %s", found.LicenseID)
	}
}

func TestDeleteVehicle_Success(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	driverID := createVehicleDriver(t, d)
	vehicle, err := d.CreateVehicle(VehicleDTO{VecID: "VEC-DEL", LicenseID: "LIC-DEL", VIN: "7T2BK1BA4KC123456", LicenseExpiration: "2028-03-31", DriverID: driverID})
	if err != nil {
		t.Fatalf("create should succeed: %v", err)
	}
	if err := d.DeleteVehicle(vehicle.ID); err != nil {
		t.Fatalf("delete should succeed: %v", err)
	}
	if _, err := d.GetVehicle(vehicle.ID); err == nil {
		t.Error("expected not found after delete")
	}
}
