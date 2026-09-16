package main

import (
	"fmt"
	"log"
)

// createTablesWithNewSchema creates all database tables including vehicle, well, and load tables
func (d *Database) createTablesWithNewSchema() error {
	// Driver table - stores driver information
	createDriverTable := `
	CREATE TABLE IF NOT EXISTS driver (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		license_number TEXT UNIQUE NOT NULL,
		phone TEXT,
		email TEXT,
		address TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	// Loader table - stores loader information
	createLoaderTable := `
	CREATE TABLE IF NOT EXISTS loader (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		phone TEXT,
		email TEXT,
		address TEXT,
		company TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	// Vehicle table - stores vehicle information
	createVehicleTable := `
	CREATE TABLE IF NOT EXISTS vehicle (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		vec_id TEXT UNIQUE NOT NULL,
		license_id TEXT UNIQUE NOT NULL,
		trailer_number INTEGER,
		vin TEXT UNIQUE NOT NULL,
		license_expiration DATE NOT NULL,
		driver_id INTEGER NOT NULL,
		overweight_permit_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(driver_id) REFERENCES driver(id) ON DELETE RESTRICT
	);
	`

	// Well table - stores well/location information
	createWellTable := `
	CREATE TABLE IF NOT EXISTS well (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		client TEXT NOT NULL,
		location TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(name, client)
	);
	`

	// Load table - stores load/shipment information
	createLoadTable := `
	CREATE TABLE IF NOT EXISTS load (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date DATE NOT NULL,
		load_number TEXT UNIQUE NOT NULL,
		loader_id INTEGER NOT NULL,
		well_id INTEGER NOT NULL,
		ticket_number TEXT UNIQUE NOT NULL,
		miles REAL,
		driver_id INTEGER NOT NULL,
		vec_id TEXT NOT NULL,
		trailer_number INTEGER,
		net_weight REAL NOT NULL,
		tons REAL NOT NULL,
		ton_rate REAL NOT NULL,
		total_value REAL NOT NULL,
		status BOOLEAN DEFAULT 1,
		ticket_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(loader_id) REFERENCES loader(id) ON DELETE RESTRICT,
		FOREIGN KEY(well_id) REFERENCES well(id) ON DELETE RESTRICT,
		FOREIGN KEY(driver_id) REFERENCES driver(id) ON DELETE RESTRICT,
		FOREIGN KEY(vec_id) REFERENCES vehicle(vec_id) ON DELETE RESTRICT
	);
	`

	// Create indexes for better query performance
	createIndexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_vehicle_driver_id ON vehicle(driver_id);`,
		`CREATE INDEX IF NOT EXISTS idx_vehicle_license_id ON vehicle(license_id);`,
		`CREATE INDEX IF NOT EXISTS idx_well_client ON well(client);`,
		`CREATE INDEX IF NOT EXISTS idx_load_date ON load(date);`,
		`CREATE INDEX IF NOT EXISTS idx_load_driver_id ON load(driver_id);`,
		`CREATE INDEX IF NOT EXISTS idx_load_well_id ON load(well_id);`,
		`CREATE INDEX IF NOT EXISTS idx_load_loader_id ON load(loader_id);`,
		`CREATE INDEX IF NOT EXISTS idx_load_status ON load(status);`,
		`CREATE INDEX IF NOT EXISTS idx_load_load_number ON load(load_number);`,
	}

	// Execute table creation queries
	queries := []string{
		createDriverTable,
		createLoaderTable,
		createVehicleTable,
		createWellTable,
		createLoadTable,
	}

	for _, query := range queries {
		if _, err := d.db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
		log.Printf("Table created successfully: %s", query[:50])
	}

	// Execute index creation queries
	for _, query := range createIndexes {
		if _, err := d.db.Exec(query); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	log.Println("All indexes created successfully")
	return nil
}
