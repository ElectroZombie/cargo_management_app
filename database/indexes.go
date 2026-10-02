package database

import (
	"fmt"
)

// indexStatements returns all index creation statements for query optimization.
// These indexes are designed to support:
// - Foreign key lookups and relationships
// - Frequent filtering and sorting operations
// - Status-based queries with partial indexes
// - Date range searches
func indexStatements() []string {
	return []string{
		// Driver indexes
		"CREATE INDEX IF NOT EXISTS idx_driver_license_number ON driver(license_number)",
		"CREATE INDEX IF NOT EXISTS idx_driver_name ON driver(name)",
		"CREATE INDEX IF NOT EXISTS idx_driver_created_at ON driver(created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_driver_updated_at ON driver(updated_at DESC)",

		// Vehicle indexes
		"CREATE INDEX IF NOT EXISTS idx_vehicle_driver_id ON vehicle(driver_id)",
		"CREATE INDEX IF NOT EXISTS idx_vehicle_vin ON vehicle(vin)",
		"CREATE INDEX IF NOT EXISTS idx_vehicle_license_id ON vehicle(license_id)",
		"CREATE INDEX IF NOT EXISTS idx_vehicle_vec_id ON vehicle(vec_id)",
		"CREATE INDEX IF NOT EXISTS idx_vehicle_created_at ON vehicle(created_at DESC)",

		// Well indexes
		"CREATE INDEX IF NOT EXISTS idx_well_name ON well(name)",
		"CREATE INDEX IF NOT EXISTS idx_well_location ON well(location)",
		"CREATE INDEX IF NOT EXISTS idx_well_created_at ON well(created_at DESC)",

		// Loader indexes
		"CREATE INDEX IF NOT EXISTS idx_loader_name ON loader(name)",
		"CREATE INDEX IF NOT EXISTS idx_loader_company ON loader(company)",
		"CREATE INDEX IF NOT EXISTS idx_loader_created_at ON loader(created_at DESC)",

		// Load indexes (most complex due to multiple filter combinations)
		"CREATE INDEX IF NOT EXISTS idx_load_driver_id ON load(driver_id)",
		"CREATE INDEX IF NOT EXISTS idx_load_well_id ON load(well_id)",
		"CREATE INDEX IF NOT EXISTS idx_load_loader_id ON load(loader_id)",
		"CREATE INDEX IF NOT EXISTS idx_load_status ON load(status)",
		"CREATE INDEX IF NOT EXISTS idx_load_date ON load(date DESC)",
		"CREATE INDEX IF NOT EXISTS idx_load_load_number ON load(load_number)",
		"CREATE INDEX IF NOT EXISTS idx_load_ticket_number ON load(ticket_number)",

		// Composite indexes for common query patterns
		"CREATE INDEX IF NOT EXISTS idx_load_driver_status ON load(driver_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_load_well_status ON load(well_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_load_loader_status ON load(loader_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_load_date_status ON load(date, status)",

		// Partial indexes for active loads (common filter)
		"CREATE INDEX IF NOT EXISTS idx_load_active ON load(driver_id, date DESC) WHERE status = 1",
		"CREATE INDEX IF NOT EXISTS idx_load_inactive ON load(driver_id, date DESC) WHERE status = 0",

		// Date range search optimization
		"CREATE INDEX IF NOT EXISTS idx_load_date_created ON load(date, created_at DESC)",
	}
}

// createIndexes executes all index creation statements.
func (d *Database) createIndexes() error {
	for _, statement := range indexStatements() {
		if _, err := d.db.Exec(statement); err != nil {
			return fmt.Errorf("create index: %w", err)
		}
	}
	return nil
}
