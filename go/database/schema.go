package database

// schemaStatements is the single source of truth for the operational schema.
// Authentication and user-account tables are intentionally not part of this design.
// Indexes live in indexes.go and are applied after the tables are created.
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

		`CREATE TABLE IF NOT EXISTS loader (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			phone TEXT NOT NULL DEFAULT '',
			email TEXT NOT NULL DEFAULT '',
			address TEXT NOT NULL DEFAULT '',
			company TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS vehicle (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			vec_id TEXT NOT NULL UNIQUE,
			license_id TEXT NOT NULL,
			trailer_number INTEGER NOT NULL DEFAULT 0,
			vin TEXT NOT NULL UNIQUE,
			license_expiration TEXT NOT NULL,
			driver_id INTEGER NOT NULL,
			overweight_permit_id TEXT,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(driver_id) REFERENCES driver(id)
		)`,

		`CREATE TABLE IF NOT EXISTS well (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			client TEXT NOT NULL DEFAULT '',
			location TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS load (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			load_number TEXT NOT NULL UNIQUE,
			loader_id INTEGER NOT NULL,
			well_id INTEGER NOT NULL,
			ticket_number TEXT NOT NULL,
			miles REAL NOT NULL DEFAULT 0,
			driver_id INTEGER NOT NULL,
			vec_id TEXT NOT NULL,
			trailer_number INTEGER NOT NULL DEFAULT 0,
			net_weight REAL NOT NULL DEFAULT 0,
			tons REAL NOT NULL DEFAULT 0,
			ton_rate REAL NOT NULL DEFAULT 0,
			total_value REAL NOT NULL DEFAULT 0,
			status BOOLEAN NOT NULL DEFAULT 0,
			ticket_id TEXT,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(loader_id) REFERENCES loader(id),
			FOREIGN KEY(well_id) REFERENCES well(id),
			FOREIGN KEY(driver_id) REFERENCES driver(id),
			FOREIGN KEY(vec_id) REFERENCES vehicle(vec_id)
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
	}
}
