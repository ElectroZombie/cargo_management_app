package database

// schemaStatements is the single source of truth for the operational schema.
// Authentication and user-account tables are intentionally not part of this design.
func schemaStatements() []string {
    return []string{
        `CREATE TABLE IF NOT EXISTS driver (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, license_number TEXT NOT NULL UNIQUE, phone TEXT NOT NULL DEFAULT '', email TEXT NOT NULL DEFAULT '', address TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP)`,
        `CREATE TABLE IF NOT EXISTS loader (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, phone TEXT NOT NULL DEFAULT '', email TEXT NOT NULL DEFAULT '', address TEXT NOT NULL DEFAULT '', company TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP)`,
        `CREATE TABLE IF NOT EXISTS vehicle (id INTEGER PRIMARY KEY AUTOINCREMENT, vec_id TEXT NOT NULL UNIQUE, license_id TEXT NOT NULL UNIQUE, trailer_number INTEGER NOT NULL DEFAULT 0, vin TEXT NOT NULL UNIQUE, license_expiration TEXT NOT NULL, driver_id INTEGER NOT NULL, overweight_permit_id TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY(driver_id) REFERENCES driver(id) ON DELETE RESTRICT)`,
        `CREATE TABLE IF NOT EXISTS well (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, client TEXT NOT NULL, location TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, UNIQUE(name, client))`,
        `CREATE TABLE IF NOT EXISTS load (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT NOT NULL, load_number TEXT NOT NULL UNIQUE, loader_id INTEGER NOT NULL, well_id INTEGER NOT NULL, ticket_number TEXT NOT NULL UNIQUE, miles REAL NOT NULL DEFAULT 0, driver_id INTEGER NOT NULL, vec_id TEXT NOT NULL, trailer_number INTEGER NOT NULL DEFAULT 0, net_weight REAL NOT NULL, tons REAL NOT NULL, ton_rate REAL NOT NULL, total_value REAL NOT NULL, status INTEGER NOT NULL DEFAULT 1 CHECK(status IN (0,1)), ticket_id TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY(loader_id) REFERENCES loader(id) ON DELETE RESTRICT, FOREIGN KEY(well_id) REFERENCES well(id) ON DELETE RESTRICT, FOREIGN KEY(driver_id) REFERENCES driver(id) ON DELETE RESTRICT, FOREIGN KEY(vec_id) REFERENCES vehicle(vec_id) ON DELETE RESTRICT)`,
        `CREATE INDEX IF NOT EXISTS idx_vehicle_driver ON vehicle(driver_id)`,
        `CREATE INDEX IF NOT EXISTS idx_load_date ON load(date)`,
        `CREATE INDEX IF NOT EXISTS idx_load_driver ON load(driver_id)`,
        `CREATE INDEX IF NOT EXISTS idx_load_loader ON load(loader_id)`,
        `CREATE INDEX IF NOT EXISTS idx_load_well ON load(well_id)`,
        `CREATE INDEX IF NOT EXISTS idx_load_vehicle ON load(vec_id)`,
        `CREATE INDEX IF NOT EXISTS idx_load_status ON load(status)`,
    }
}
