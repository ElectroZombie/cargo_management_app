package main

import (
	"database/sql"
	"fmt"
	"log"
)

// ============ DRIVER CRUD OPERATIONS ============

// CreateDriver adds a new driver to the database
func (d *Database) CreateDriver(name, license, phone, email, address string) (*Driver, error) {
	query := `
	INSERT INTO driver (name, license_number, phone, email, address, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := d.db.Exec(query, name, license, phone, email, address)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	driver := &Driver{
		ID:      int(id),
		Name:    name,
		License: license,
		Phone:   phone,
		Email:   email,
		Address: address,
	}

	log.Printf("Driver created successfully: %s (ID: %d)", name, id)
	return driver, nil
}

// GetAllDrivers retrieves all drivers from the database
func (d *Database) GetAllDrivers() ([]Driver, error) {
	query := `SELECT id, name, license_number, phone, email, address, created_at, updated_at FROM driver`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query drivers: %w", err)
	}
	defer rows.Close()

	var drivers []Driver
	for rows.Next() {
		var driver Driver
		if err := rows.Scan(&driver.ID, &driver.Name, &driver.License, &driver.Phone, &driver.Email, &driver.Address, &driver.CreatedAt, &driver.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan driver: %w", err)
		}
		drivers = append(drivers, driver)
	}

	return drivers, nil
}

// GetDriverByID retrieves a specific driver by ID
func (d *Database) GetDriverByID(id int) (*Driver, error) {
	query := `SELECT id, name, license_number, phone, email, address, created_at, updated_at FROM driver WHERE id = ?`

	var driver Driver
	err := d.db.QueryRow(query, id).Scan(&driver.ID, &driver.Name, &driver.License, &driver.Phone, &driver.Email, &driver.Address, &driver.CreatedAt, &driver.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("driver not found")
		}
		return nil, fmt.Errorf("failed to query driver: %w", err)
	}

	return &driver, nil
}

// UpdateDriver updates an existing driver
func (d *Database) UpdateDriver(id int, name, license, phone, email, address string) error {
	query := `
	UPDATE driver 
	SET name = ?, license_number = ?, phone = ?, email = ?, address = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	result, err := d.db.Exec(query, name, license, phone, email, address, id)
	if err != nil {
		return fmt.Errorf("failed to update driver: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("driver not found")
	}

	log.Printf("Driver updated successfully: ID %d", id)
	return nil
}

// DeleteDriver removes a driver from the database
func (d *Database) DeleteDriver(id int) error {
	query := `DELETE FROM driver WHERE id = ?`

	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete driver: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("driver not found")
	}

	log.Printf("Driver deleted successfully: ID %d", id)
	return nil
}
