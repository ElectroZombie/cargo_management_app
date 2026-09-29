package database

// schemaStatements is the single source of truth for the operational schema.
// Authentication and user-account tables are intentionally not part of this design.
func schemaStatements() []string {
	return []string{
		// Core operational entities
		`CREATE TABLE IF NOT EXISTS driver (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			license_number TEXT NOT NULL UNIQUE,
			phone TEXT NOT NULL DEFAULT '',
			email TEXT NOT NULL DEFAULT '',
			address TEXT NOT NULL DEFAULT '',
			meta TEXT,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS truck (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			vin TEXT NOT NULL UNIQUE,
			license_expiration TEXT NOT NULL,
			driver_id INTEGER NOT NULL,
			overweight_permit_id TEXT,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(driver_id) REFERENCES driver(id)
		)`,

		`CREATE TABLE IF NOT EXISTS trailer (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			number INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS well (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			location TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS trip (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			load_number TEXT NOT NULL UNIQUE,
			well_id INTEGER NOT NULL,
			ticket_number TEXT NOT NULL,
			driver_id INTEGER NOT NULL,
			truck_id TEXT NOT NULL,
			id_trailer INTEGER NOT NULL,
			meta TEXT,
			status BOOLEAN NOT NULL DEFAULT 0,
			ticket_id TEXT,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(well_id) REFERENCES well(id),
			FOREIGN KEY(driver_id) REFERENCES driver(id),
			FOREIGN KEY(id_trailer) REFERENCES trailer(id)
		)`,

		// Financial entities
		`CREATE TABLE IF NOT EXISTS category (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL UNIQUE,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS expense (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			concept TEXT NOT NULL,
			category_id INTEGER NOT NULL,
			amount REAL NOT NULL DEFAULT 0,
			quantity REAL NOT NULL DEFAULT 0,
			total REAL NOT NULL DEFAULT 0,
			receipt_id INTEGER,
			meta TEXT,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(category_id) REFERENCES category(id)
		)`,

		`CREATE TABLE IF NOT EXISTS income (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			concept TEXT NOT NULL,
			category_id INTEGER NOT NULL,
			amount REAL NOT NULL DEFAULT 0,
			quantity REAL NOT NULL DEFAULT 0,
			total REAL NOT NULL DEFAULT 0,
			receipt_id INTEGER,
			meta TEXT,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(category_id) REFERENCES category(id)
		)`,

		`CREATE TABLE IF NOT EXISTS fuel_expense (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			id_expense INTEGER NOT NULL,
			id_driver INTEGER NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(id_expense) REFERENCES expense(id),
			FOREIGN KEY(id_driver) REFERENCES driver(id)
		)`,

		`CREATE TABLE IF NOT EXISTS stats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			number_week INTEGER NOT NULL,
			date_range_init TEXT NOT NULL,
			date_range_end TEXT NOT NULL,
			expense_total REAL NOT NULL DEFAULT 0,
			income_total REAL NOT NULL DEFAULT 0,
			net_total REAL NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,

		// Indexes for performance
		`CREATE INDEX IF NOT EXISTS idx_truck_driver ON truck(driver_id)`,
		`CREATE INDEX IF NOT EXISTS idx_trip_date ON trip(date)`,
		`CREATE INDEX IF NOT EXISTS idx_trip_driver ON trip(driver_id)`,
		`CREATE INDEX IF NOT EXISTS idx_trip_well ON trip(well_id)`,
		`CREATE INDEX IF NOT EXISTS idx_trip_truck ON trip(truck_id)`,
		`CREATE INDEX IF NOT EXISTS idx_trip_trailer ON trip(id_trailer)`,
		`CREATE INDEX IF NOT EXISTS idx_trip_status ON trip(status)`,
		`CREATE INDEX IF NOT EXISTS idx_expense_date ON expense(date)`,
		`CREATE INDEX IF NOT EXISTS idx_expense_category ON expense(category_id)`,
		`CREATE INDEX IF NOT EXISTS idx_income_date ON income(date)`,
		`CREATE INDEX IF NOT EXISTS idx_income_category ON income(category_id)`,
		`CREATE INDEX IF NOT EXISTS idx_fuel_expense_driver ON fuel_expense(id_driver)`,
		`CREATE INDEX IF NOT EXISTS idx_fuel_expense_expense ON fuel_expense(id_expense)`,
		`CREATE INDEX IF NOT EXISTS idx_stats_week ON stats(number_week)`,
	}
}
