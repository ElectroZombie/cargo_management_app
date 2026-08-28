package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{db: db}

	// Create tables
	if err := database.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database initialized successfully")
	return database, nil
}

// createTables creates the necessary database tables
func (d *Database) createTables() error {
	createCargoTable := `
	CREATE TABLE IF NOT EXISTS cargo (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'Pending',
		weight REAL,
		destination TEXT NOT NULL,
		date TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT DEFAULT 'user',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	createShipmentsTable := `
	CREATE TABLE IF NOT EXISTS shipments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cargo_id INTEGER NOT NULL,
		origin TEXT NOT NULL,
		destination TEXT NOT NULL,
		departure_date TEXT,
		arrival_date TEXT,
		status TEXT DEFAULT 'In Transit',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(cargo_id) REFERENCES cargo(id) ON DELETE CASCADE
	);
	`

	queries := []string{createCargoTable, createUsersTable, createShipmentsTable}

	for _, query := range queries {
		if _, err := d.db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// ============ CARGO OPERATIONS ============

// Cargo represents a cargo item
type Cargo struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Status      string  `json:"status"`
	Weight      float64 `json:"weight"`
	Destination string  `json:"destination"`
	Date        string  `json:"date"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// CreateCargo adds a new cargo to the database
func (d *Database) CreateCargo(name, status, destination, date string, weight float64) (*Cargo, error) {
	query := `
	INSERT INTO cargo (name, status, destination, date, weight)
	VALUES (?, ?, ?, ?, ?)
	`

	result, err := d.db.Exec(query, name, status, destination, date, weight)
	if err != nil {
		return nil, fmt.Errorf("failed to create cargo: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	cargo := &Cargo{
		ID:          int(id),
		Name:        name,
		Status:      status,
		Weight:      weight,
		Destination: destination,
		Date:        date,
	}

	return cargo, nil
}

// GetAllCargo retrieves all cargo items
func (d *Database) GetAllCargo() ([]Cargo, error) {
	query := `SELECT id, name, status, weight, destination, date, created_at, updated_at FROM cargo`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query cargo: %w", err)
	}
	defer rows.Close()

	var cargoList []Cargo
	for rows.Next() {
		var cargo Cargo
		if err := rows.Scan(&cargo.ID, &cargo.Name, &cargo.Status, &cargo.Weight, &cargo.Destination, &cargo.Date, &cargo.CreatedAt, &cargo.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan cargo: %w", err)
		}
		cargoList = append(cargoList, cargo)
	}

	return cargoList, nil
}

// GetCargoByID retrieves a specific cargo by ID
func (d *Database) GetCargoByID(id int) (*Cargo, error) {
	query := `SELECT id, name, status, weight, destination, date, created_at, updated_at FROM cargo WHERE id = ?`

	var cargo Cargo
	err := d.db.QueryRow(query, id).Scan(&cargo.ID, &cargo.Name, &cargo.Status, &cargo.Weight, &cargo.Destination, &cargo.Date, &cargo.CreatedAt, &cargo.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cargo not found")
		}
		return nil, fmt.Errorf("failed to query cargo: %w", err)
	}

	return &cargo, nil
}

// UpdateCargo updates an existing cargo item
func (d *Database) UpdateCargo(id int, name, status, destination, date string, weight float64) error {
	query := `
	UPDATE cargo 
	SET name = ?, status = ?, destination = ?, date = ?, weight = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	result, err := d.db.Exec(query, name, status, destination, date, weight, id)
	if err != nil {
		return fmt.Errorf("failed to update cargo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cargo not found")
	}

	return nil
}

// DeleteCargo removes a cargo item
func (d *Database) DeleteCargo(id int) error {
	query := `DELETE FROM cargo WHERE id = ?`

	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete cargo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cargo not found")
	}

	return nil
}

// GetCargoStats returns statistics about cargo
func (d *Database) GetCargoStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total cargo count
	var totalCount int
	err := d.db.QueryRow(`SELECT COUNT(*) FROM cargo`).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}
	stats["total_cargo"] = totalCount

	// Count by status
	statusCounts := make(map[string]int)
	rows, err := d.db.Query(`SELECT status, COUNT(*) FROM cargo GROUP BY status`)
	if err != nil {
		return nil, fmt.Errorf("failed to get status counts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan status counts: %w", err)
		}
		statusCounts[status] = count
	}
	stats["by_status"] = statusCounts

	// Total weight
	var totalWeight sql.NullFloat64
	err = d.db.QueryRow(`SELECT SUM(weight) FROM cargo`).Scan(&totalWeight)
	if err != nil {
		return nil, fmt.Errorf("failed to get total weight: %w", err)
	}
	stats["total_weight"] = totalWeight.Float64

	return stats, nil
}

// ============ SHIPMENT OPERATIONS ============

// Shipment represents a shipment
type Shipment struct {
	ID            int    `json:"id"`
	CargoID       int    `json:"cargo_id"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureDate string `json:"departure_date"`
	ArrivalDate   string `json:"arrival_date"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// CreateShipment creates a new shipment
func (d *Database) CreateShipment(cargoID int, origin, destination, departureDate, arrivalDate, status string) (*Shipment, error) {
	query := `
	INSERT INTO shipments (cargo_id, origin, destination, departure_date, arrival_date, status)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := d.db.Exec(query, cargoID, origin, destination, departureDate, arrivalDate, status)
	if err != nil {
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	shipment := &Shipment{
		ID:            int(id),
		CargoID:       cargoID,
		Origin:        origin,
		Destination:   destination,
		DepartureDate: departureDate,
		ArrivalDate:   arrivalDate,
		Status:        status,
	}

	return shipment, nil
}

// GetShipmentsByCargoID retrieves all shipments for a cargo
func (d *Database) GetShipmentsByCargoID(cargoID int) ([]Shipment, error) {
	query := `
	SELECT id, cargo_id, origin, destination, departure_date, arrival_date, status, created_at, updated_at 
	FROM shipments 
	WHERE cargo_id = ?
	`

	rows, err := d.db.Query(query, cargoID)
	if err != nil {
		return nil, fmt.Errorf("failed to query shipments: %w", err)
	}
	defer rows.Close()

	var shipments []Shipment
	for rows.Next() {
		var shipment Shipment
		if err := rows.Scan(&shipment.ID, &shipment.CargoID, &shipment.Origin, &shipment.Destination, &shipment.DepartureDate, &shipment.ArrivalDate, &shipment.Status, &shipment.CreatedAt, &shipment.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan shipment: %w", err)
		}
		shipments = append(shipments, shipment)
	}

	return shipments, nil
}

// UpdateShipmentStatus updates the status of a shipment
func (d *Database) UpdateShipmentStatus(id int, status string) error {
	query := `UPDATE shipments SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	result, err := d.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update shipment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("shipment not found")
	}

	return nil
}
