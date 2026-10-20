package database

import (
	"testing"

	domainerrors "cargo_management_app/go/errors"
)

func TestCreateWell_Success(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	dto := WellDTO{Name: "North Basin", Client: "Acme Oil", Location: "West Field"}
	well, err := d.CreateWell(dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if well.Name != "North Basin" {
		t.Errorf("expected North Basin, got %s", well.Name)
	}
	if well.Client != "Acme Oil" {
		t.Errorf("expected Acme Oil, got %s", well.Client)
	}
}

func TestCreateWell_MissingLocation(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	_, err := d.CreateWell(WellDTO{Name: "No Location"})
	if err == nil {
		t.Fatal("expected validation error for empty location")
	}
	if !domainerrors.IsValidation(err) {
		t.Errorf("expected a validation error, got %v", err)
	}
}

func TestGetWell_NotFound(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	_, err := d.GetWell(999)
	if err == nil {
		t.Fatal("expected not found error")
	}
	if !domainerrors.IsNotFound(err) {
		t.Errorf("expected a not found error, got %v", err)
	}
}

func TestUpdateWell_Success(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	well, err := d.CreateWell(WellDTO{Name: "Original", Location: "Field A"})
	if err != nil {
		t.Fatalf("create should succeed: %v", err)
	}

	updateDTO := WellDTO{Name: "Updated", Client: well.Client, Location: well.Location}
	if err := d.UpdateWell(well.ID, updateDTO); err != nil {
		t.Fatalf("update should succeed: %v", err)
	}

	found, err := d.GetWell(well.ID)
	if err != nil {
		t.Fatalf("get should succeed: %v", err)
	}
	if found.Name != "Updated" {
		t.Errorf("expected Updated, got %s", found.Name)
	}
}

func TestDeleteWell_Success(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	well, err := d.CreateWell(WellDTO{Name: "ToBeDeleted", Location: "Field B"})
	if err != nil {
		t.Fatalf("create should succeed: %v", err)
	}
	if err := d.DeleteWell(well.ID); err != nil {
		t.Fatalf("delete should succeed: %v", err)
	}
	if _, err := d.GetWell(well.ID); err == nil {
		t.Error("expected not found after delete")
	}
}

func TestDeleteWell_Referenced(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	loaderID, wellID, driverID, vecID := createLoadRefs(t, d)

	load := LoadDTO{
		Date:         "2026-03-01",
		LoadNumber:   "LDR-WELL",
		LoaderID:     loaderID,
		WellID:       wellID,
		TicketNumber: "TK-WELL",
		Miles:        60,
		DriverID:     driverID,
		VecID:        vecID,
		NetWeight:    60,
		Tons:         30,
		TonRate:      9,
		TotalValue:   270,
		Status:       true,
	}
	if _, err := d.CreateLoad(load); err != nil {
		t.Fatalf("create load should succeed: %v", err)
	}

	// Deleting a well that still has assigned loads must fail on the FK constraint.
	err := d.DeleteWell(wellID)
	if err == nil {
		t.Fatal("expected FK constraint error when deleting well with assigned loads")
	}
	if !domainerrors.IsConflict(err) {
		t.Errorf("expected a conflict error, got %v", err)
	}
}
