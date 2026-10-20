package database

import "testing"

// createLoadRefs seeds the loader, well, driver, and vehicle a load depends on
// and returns the identifiers needed to build a valid LoadDTO.
func createLoadRefs(t *testing.T, d *Database) (loaderID, wellID, driverID int, vecID string) {
	t.Helper()

	driver, err := d.CreateDriver(DriverDTO{Name: "Load Driver", LicenseNumber: "DL-LDR", Phone: "5551111111", Email: "load@test.com", Address: "1 Test St"})
	if err != nil {
		t.Fatalf("create driver should succeed: %v", err)
	}
	loader, err := d.CreateLoader(LoaderDTO{Name: "Load Loader", Company: "Co"})
	if err != nil {
		t.Fatalf("create loader should succeed: %v", err)
	}
	well, err := d.CreateWell(WellDTO{Name: "Load Well", Location: "Field"})
	if err != nil {
		t.Fatalf("create well should succeed: %v", err)
	}
	vehicle, err := d.CreateVehicle(VehicleDTO{
		VecID:             "VEC-LDR",
		LicenseID:         "LIC-LDR",
		VIN:               "2T2BK1BA4KC123456",
		LicenseExpiration: "2027-06-30",
		DriverID:          driver.ID,
	})
	if err != nil {
		t.Fatalf("create vehicle should succeed: %v", err)
	}

	return loader.ID, well.ID, driver.ID, vehicle.VecID
}

func TestCreateLoad_Success(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	loaderID, wellID, driverID, vecID := createLoadRefs(t, d)

	dto := LoadDTO{Date: "2026-01-15", LoadNumber: "LDR-010", LoaderID: loaderID, WellID: wellID, TicketNumber: "TK-100", Miles: 150, DriverID: driverID, VecID: vecID, TrailerNumber: 2, NetWeight: 200, Tons: 100, TonRate: 25, TotalValue: 2500, Status: true, TicketID: "TICKET-001"}
	load, err := d.CreateLoad(dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if load.LoadNumber != "LDR-010" {
		t.Errorf("expected LDR-010, got %s", load.LoadNumber)
	}
}

func TestCreateLoad_InvalidLoaderID(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	_, wellID, driverID, vecID := createLoadRefs(t, d)

	dto := LoadDTO{Date: "2026-01-15", LoadNumber: "LDR-011", LoaderID: 0, WellID: wellID, TicketNumber: "TK-101", Miles: 50, DriverID: driverID, VecID: vecID, TrailerNumber: 0, NetWeight: 50, Tons: 25, TonRate: 10, TotalValue: 250, Status: true, TicketID: "TK-101"}
	if _, err := d.CreateLoad(dto); err == nil {
		t.Error("expected validation error for invalid loader_id=0")
	}
}

func TestCreateLoad_FKConstraint_Driver(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	loaderID, wellID, _, vecID := createLoadRefs(t, d)

	dto := LoadDTO{Date: "2026-01-15", LoadNumber: "LDR-012", LoaderID: loaderID, WellID: wellID, TicketNumber: "TK-102", Miles: 50, DriverID: 999, VecID: vecID, TrailerNumber: 0, NetWeight: 50, Tons: 25, TonRate: 10, TotalValue: 250, Status: true, TicketID: "TK-102"}
	if _, err := d.CreateLoad(dto); err == nil {
		t.Error("expected FK constraint error for invalid driver_id=999")
	}
}
