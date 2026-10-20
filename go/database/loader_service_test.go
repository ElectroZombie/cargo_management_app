package database

import (
	"testing"

	domainerrors "cargo_management_app/go/errors"
)

func TestCreateLoader_Success(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	dto := LoaderDTO{Name: "Carlos Ruiz", Phone: "5552002002", Email: "carlos@example.com", Address: "78 Pine Rd", Company: "North Line"}
	loader, err := d.CreateLoader(dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loader.Name != "Carlos Ruiz" {
		t.Errorf("expected Carlos Ruiz, got %s", loader.Name)
	}
	if loader.Company != "North Line" {
		t.Errorf("expected North Line, got %s", loader.Company)
	}
}

func TestCreateLoader_MissingName(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	_, err := d.CreateLoader(LoaderDTO{Company: "No Name Co"})
	if err == nil {
		t.Fatal("expected validation error for empty name")
	}
	if !domainerrors.IsValidation(err) {
		t.Errorf("expected a validation error, got %v", err)
	}
}

func TestCreateLoader_InvalidEmail(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	_, err := d.CreateLoader(LoaderDTO{Name: "Bad Email", Email: "not-an-email"})
	if err == nil {
		t.Fatal("expected validation error for malformed email")
	}
	if !domainerrors.IsValidation(err) {
		t.Errorf("expected a validation error, got %v", err)
	}
}

func TestGetLoader_NotFound(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	_, err := d.GetLoader(999)
	if err == nil {
		t.Fatal("expected not found error")
	}
	if !domainerrors.IsNotFound(err) {
		t.Errorf("expected a not found error, got %v", err)
	}
}

func TestUpdateLoader_Success(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	loader, err := d.CreateLoader(LoaderDTO{Name: "Original", Company: "Co"})
	if err != nil {
		t.Fatalf("create should succeed: %v", err)
	}

	updateDTO := LoaderDTO{Name: "Updated", Phone: loader.Phone, Email: loader.Email, Address: loader.Address, Company: loader.Company}
	if err := d.UpdateLoader(loader.ID, updateDTO); err != nil {
		t.Fatalf("update should succeed: %v", err)
	}

	found, err := d.GetLoader(loader.ID)
	if err != nil {
		t.Fatalf("get should succeed: %v", err)
	}
	if found.Name != "Updated" {
		t.Errorf("expected Updated, got %s", found.Name)
	}
}

func TestDeleteLoader_Success(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	loader, err := d.CreateLoader(LoaderDTO{Name: "ToBeDeleted"})
	if err != nil {
		t.Fatalf("create should succeed: %v", err)
	}
	if err := d.DeleteLoader(loader.ID); err != nil {
		t.Fatalf("delete should succeed: %v", err)
	}
	if _, err := d.GetLoader(loader.ID); err == nil {
		t.Error("expected not found after delete")
	}
}

func TestDeleteLoader_Referenced(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	loaderID, wellID, driverID, vecID := createLoadRefs(t, d)

	load := LoadDTO{
		Date:         "2026-02-01",
		LoadNumber:   "LDR-LOADER",
		LoaderID:     loaderID,
		WellID:       wellID,
		TicketNumber: "TK-LOADER",
		Miles:        80,
		DriverID:     driverID,
		VecID:        vecID,
		NetWeight:    80,
		Tons:         40,
		TonRate:      12,
		TotalValue:   480,
		Status:       true,
	}
	if _, err := d.CreateLoad(load); err != nil {
		t.Fatalf("create load should succeed: %v", err)
	}

	// Deleting a loader that still has assigned loads must fail on the FK constraint.
	err := d.DeleteLoader(loaderID)
	if err == nil {
		t.Fatal("expected FK constraint error when deleting loader with assigned loads")
	}
	if !domainerrors.IsConflict(err) {
		t.Errorf("expected a conflict error, got %v", err)
	}
}

func TestSearchLoaders_Fuzzy(t *testing.T) {
	t.Helper()
	d := setupTestDB(t)

	if _, err := d.CreateLoader(LoaderDTO{Name: "Maria Gomez", Company: "South Point"}); err != nil {
		t.Fatalf("create should succeed: %v", err)
	}

	// A near-miss typo still clears the similarity threshold.
	items, err := d.SearchLoaders("Maria Gome")
	if err != nil {
		t.Fatalf("search should succeed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 match, got %d", len(items))
	}
	if items[0].Name != "Maria Gomez" {
		t.Errorf("expected Maria Gomez, got %s", items[0].Name)
	}
}
