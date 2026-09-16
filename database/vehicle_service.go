package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// ============ VEHICLE CRUD OPERATIONS ============

// CreateVehicle adds a new vehicle to the database
func (d *Database) CreateVehicle(vecID, licenseID string, trailerNumber int, vin string, licenseExpiration time.Time, driverID int, overweightPermitID string) (*Vehicle, error) {
	query := `
	INSERT INTO vehicle (vec_id, license_id, trailer_number, vin, license_expiration, driver_id, overweight_permit_id, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := d.db.Exec(query, vecID, licenseID, trailerNumber, vin, licenseExpiration, driverID, overweightPermitID)
	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	vehicle := &Vehicle{
		ID:                 int(id),
		VecID:              vecID,
		LicenseID:          licenseID,
		TrailerNumber:      trailerNumber,
		VIN:                vin,
		LicenseExpiration:  licenseExpiration,
		DriverID:           driverID,
		OverweightPermitID: overweightPermitID,
	}

	log.Printf("Vehicle created successfully: %s (ID: %d)", vecID, id)
	return vehicle, nil
}

// GetAllVehicles retrieves all vehicles from the database
func (d *Database) GetAllVehicles() ([]Vehicle, error) {
	query := `SELECT id, vec_id, license_id, trailer_number, vin, license_expiration, driver_id, overweight_permit_id, created_at, updated_at FROM vehicle`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var vehicle Vehicle
		if err := rows.Scan(&vehicle.ID, &vehicle.VecID, &vehicle.LicenseID, &vehicle.TrailerNumber, &vehicle.VIN, &vehicle.LicenseExpiration, &vehicle.DriverID, &vehicle.OverweightPermitID, &vehicle.CreatedAt, &vehicle.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}
		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

// GetVehicleByID retrieves a specific vehicle by ID
func (d *Database) GetVehicleByID(id int) (*Vehicle, error) {
	query := `SELECT id, vec_id, license_id, trailer_number, vin, license_expiration, driver_id, overweight_permit_id, created_at, updated_at FROM vehicle WHERE id = ?`

	var vehicle Vehicle
	err := d.db.QueryRow(query, id).Scan(&vehicle.ID, &vehicle.VecID, &vehicle.LicenseID, &vehicle.TrailerNumber, &vehicle.VIN, &vehicle.LicenseExpiration, &vehicle.DriverID, &vehicle.OverweightPermitID, &vehicle.CreatedAt, &vehicle.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle not found")
		}
		return nil, fmt.Errorf("failed to query vehicle: %w", err)
	}

	return &vehicle, nil
}

// GetVehicleByVecID retrieves a specific vehicle by vec_id
func (d *Database) GetVehicleByVecID(vecID string) (*Vehicle, error) {
	query := `SELECT id, vec_id, license_id, trailer_number, vin, license_expiration, driver_id, overweight_permit_id, created_at, updated_at FROM vehicle WHERE vec_id = ?`

	var vehicle Vehicle
	err := d.db.QueryRow(query, vecID).Scan(&vehicle.ID, &vehicle.VecID, &vehicle.LicenseID, &vehicle.TrailerNumber, &vehicle.VIN, &vehicle.LicenseExpiration, &vehicle.DriverID, &vehicle.OverweightPermitID, &vehicle.CreatedAt, &vehicle.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle not found")
		}
		return nil, fmt.Errorf("failed to query vehicle: %w", err)
	}

	return &vehicle, nil
}

// UpdateVehicle updates an existing vehicle
func (d *Database) UpdateVehicle(id int, vecID, licenseID string, trailerNumber int, vin string, licenseExpiration time.Time, driverID int, overweightPermitID string) error {
	query := `
	UPDATE vehicle 
	SET vec_id = ?, license_id = ?, trailer_number = ?, vin = ?, license_expiration = ?, driver_id = ?, overweight_permit_id = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	result, err := d.db.Exec(query, vecID, licenseID, trailerNumber, vin, licenseExpiration, driverID, overweightPermitID, id)
	if err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle not found")
	}

	log.Printf("Vehicle updated successfully: ID %d", id)
	return nil
}

// DeleteVehicle removes a vehicle from the database
func (d *Database) DeleteVehicle(id int) error {
	query := `DELETE FROM vehicle WHERE id = ?`

	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle not found")
	}

	log.Printf("Vehicle deleted successfully: ID %d", id)
	return nil
}
