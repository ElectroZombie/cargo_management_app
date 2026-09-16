package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// ============ LOAD CRUD OPERATIONS ============

// CreateLoad adds a new load to the database
func (d *Database) CreateLoad(date time.Time, loadNumber string, loaderID int, wellID int, ticketNumber string, miles float64, driverID int, vecID string, trailerNumber int, netWeight, tons, tonRate, totalValue float64, status bool, ticketID string) (*Load, error) {
	query := `
	INSERT INTO load (date, load_number, loader_id, well_id, ticket_number, miles, driver_id, vec_id, trailer_number, net_weight, tons, ton_rate, total_value, status, ticket_id, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := d.db.Exec(query, date, loadNumber, loaderID, wellID, ticketNumber, miles, driverID, vecID, trailerNumber, netWeight, tons, tonRate, totalValue, status, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to create load: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	load := &Load{
		ID:            int(id),
		Date:          date,
		LoadNumber:    loadNumber,
		LoaderID:      loaderID,
		WellID:        wellID,
		TicketNumber:  ticketNumber,
		Miles:         miles,
		DriverID:      driverID,
		VecID:         vecID,
		TrailerNumber: trailerNumber,
		NetWeight:     netWeight,
		Tons:          tons,
		TonRate:       tonRate,
		TotalValue:    totalValue,
		Status:        status,
		TicketID:      ticketID,
	}

	log.Printf("Load created successfully: %s (ID: %d)", loadNumber, id)
	return load, nil
}

// GetAllLoads retrieves all loads from the database
func (d *Database) GetAllLoads() ([]Load, error) {
	query := `SELECT id, date, load_number, loader_id, well_id, ticket_number, miles, driver_id, vec_id, trailer_number, net_weight, tons, ton_rate, total_value, status, ticket_id, created_at, updated_at FROM load`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query loads: %w", err)
	}
	defer rows.Close()

	var loads []Load
	for rows.Next() {
		var load Load
		if err := rows.Scan(&load.ID, &load.Date, &load.LoadNumber, &load.LoaderID, &load.WellID, &load.TicketNumber, &load.Miles, &load.DriverID, &load.VecID, &load.TrailerNumber, &load.NetWeight, &load.Tons, &load.TonRate, &load.TotalValue, &load.Status, &load.TicketID, &load.CreatedAt, &load.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan load: %w", err)
		}
		loads = append(loads, load)
	}

	return loads, nil
}

// GetLoadByID retrieves a specific load by ID
func (d *Database) GetLoadByID(id int) (*Load, error) {
	query := `SELECT id, date, load_number, loader_id, well_id, ticket_number, miles, driver_id, vec_id, trailer_number, net_weight, tons, ton_rate, total_value, status, ticket_id, created_at, updated_at FROM load WHERE id = ?`

	var load Load
	err := d.db.QueryRow(query, id).Scan(&load.ID, &load.Date, &load.LoadNumber, &load.LoaderID, &load.WellID, &load.TicketNumber, &load.Miles, &load.DriverID, &load.VecID, &load.TrailerNumber, &load.NetWeight, &load.Tons, &load.TonRate, &load.TotalValue, &load.Status, &load.TicketID, &load.CreatedAt, &load.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("load not found")
		}
		return nil, fmt.Errorf("failed to query load: %w", err)
	}

	return &load, nil
}

// UpdateLoad updates an existing load
func (d *Database) UpdateLoad(id int, date time.Time, loadNumber string, loaderID int, wellID int, ticketNumber string, miles float64, driverID int, vecID string, trailerNumber int, netWeight, tons, tonRate, totalValue float64, status bool, ticketID string) error {
	query := `
	UPDATE load 
	SET date = ?, load_number = ?, loader_id = ?, well_id = ?, ticket_number = ?, miles = ?, driver_id = ?, vec_id = ?, trailer_number = ?, net_weight = ?, tons = ?, ton_rate = ?, total_value = ?, status = ?, ticket_id = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	result, err := d.db.Exec(query, date, loadNumber, loaderID, wellID, ticketNumber, miles, driverID, vecID, trailerNumber, netWeight, tons, tonRate, totalValue, status, ticketID, id)
	if err != nil {
		return fmt.Errorf("failed to update load: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("load not found")
	}

	log.Printf("Load updated successfully: ID %d", id)
	return nil
}

// DeleteLoad removes a load from the database
func (d *Database) DeleteLoad(id int) error {
	query := `DELETE FROM load WHERE id = ?`

	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete load: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("load not found")
	}

	log.Printf("Load deleted successfully: ID %d", id)
	return nil
}

// GetLoadsByDriverID retrieves all loads for a specific driver
func (d *Database) GetLoadsByDriverID(driverID int) ([]Load, error) {
	query := `SELECT id, date, load_number, loader_id, well_id, ticket_number, miles, driver_id, vec_id, trailer_number, net_weight, tons, ton_rate, total_value, status, ticket_id, created_at, updated_at FROM load WHERE driver_id = ?`

	rows, err := d.db.Query(query, driverID)
	if err != nil {
		return nil, fmt.Errorf("failed to query loads: %w", err)
	}
	defer rows.Close()

	var loads []Load
	for rows.Next() {
		var load Load
		if err := rows.Scan(&load.ID, &load.Date, &load.LoadNumber, &load.LoaderID, &load.WellID, &load.TicketNumber, &load.Miles, &load.DriverID, &load.VecID, &load.TrailerNumber, &load.NetWeight, &load.Tons, &load.TonRate, &load.TotalValue, &load.Status, &load.TicketID, &load.CreatedAt, &load.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan load: %w", err)
		}
		loads = append(loads, load)
	}

	return loads, nil
}

// GetLoadsByWellID retrieves all loads for a specific well
func (d *Database) GetLoadsByWellID(wellID int) ([]Load, error) {
	query := `SELECT id, date, load_number, loader_id, well_id, ticket_number, miles, driver_id, vec_id, trailer_number, net_weight, tons, ton_rate, total_value, status, ticket_id, created_at, updated_at FROM load WHERE well_id = ?`

	rows, err := d.db.Query(query, wellID)
	if err != nil {
		return nil, fmt.Errorf("failed to query loads: %w", err)
	}
	defer rows.Close()

	var loads []Load
	for rows.Next() {
		var load Load
		if err := rows.Scan(&load.ID, &load.Date, &load.LoadNumber, &load.LoaderID, &load.WellID, &load.TicketNumber, &load.Miles, &load.DriverID, &load.VecID, &load.TrailerNumber, &load.NetWeight, &load.Tons, &load.TonRate, &load.TotalValue, &load.Status, &load.TicketID, &load.CreatedAt, &load.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan load: %w", err)
		}
		loads = append(loads, load)
	}

	return loads, nil
}

// GetLoadsByStatus retrieves all loads with a specific status
func (d *Database) GetLoadsByStatus(status bool) ([]Load, error) {
	query := `SELECT id, date, load_number, loader_id, well_id, ticket_number, miles, driver_id, vec_id, trailer_number, net_weight, tons, ton_rate, total_value, status, ticket_id, created_at, updated_at FROM load WHERE status = ?`

	rows, err := d.db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query loads: %w", err)
	}
	defer rows.Close()

	var loads []Load
	for rows.Next() {
		var load Load
		if err := rows.Scan(&load.ID, &load.Date, &load.LoadNumber, &load.LoaderID, &load.WellID, &load.TicketNumber, &load.Miles, &load.DriverID, &load.VecID, &load.TrailerNumber, &load.NetWeight, &load.Tons, &load.TonRate, &load.TotalValue, &load.Status, &load.TicketID, &load.CreatedAt, &load.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan load: %w", err)
		}
		loads = append(loads, load)
	}

	return loads, nil
}
