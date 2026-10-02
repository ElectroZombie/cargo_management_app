package database

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "strings"
)

// SeedData sets up repeatable test and development fixtures for the operational database.
// These records are intentionally lightweight and safe to re-run because they are idempotent.
func SeedData(db *sql.DB) error {
    if db == nil {
        return fmt.Errorf("%w: database is nil", ErrValidation)
    }

    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("begin seed transaction: %w", err)
    }
    defer func() {
        if err != nil {
            _ = tx.Rollback()
        }
    }()

    if _, err = tx.Exec(`INSERT OR IGNORE INTO driver (id, name, license_number, phone, email, address) VALUES
        (1, 'John Doe', 'DL-1001', '555-0101', 'john.doe@example.com', '123 Main St'),
        (2, 'Maria Gomez', 'DL-1002', '555-0102', 'maria.gomez@example.com', '456 Oak Ave')`); err != nil {
        return fmt.Errorf("seed drivers: %w", err)
    }

    if _, err = tx.Exec(`INSERT OR IGNORE INTO loader (id, name, phone, email, address, company) VALUES
        (1, 'Carlos Ruiz', '555-2001', 'carlos.ruiz@example.com', '78 Pine Rd', 'North Line Logistics'),
        (2, 'Alicia Stone', '555-2002', 'alicia.stone@example.com', '12 Cedar St', 'Blue Mesa Haulage')`); err != nil {
        return fmt.Errorf("seed loaders: %w", err)
    }

    if _, err = tx.Exec(`INSERT OR IGNORE INTO well (id, name, location) VALUES
        (1, 'North Basin', 'West Field'),
        (2, 'South Ridge', 'East Field')`); err != nil {
        return fmt.Errorf("seed wells: %w", err)
    }

    if _, err = tx.Exec(`INSERT OR IGNORE INTO vehicle (id, vec_id, license_id, trailer_number, vin, license_expiration, driver_id, overweight_permit_id) VALUES
        (1, 'VEC-1001', 'LIC-1001', 24, '1HGBH41JXMN109186', '2027-01-31', 1, 'OP-1001'),
        (2, 'VEC-1002', 'LIC-1002', 18, '2T2BK1BA4KC123456', '2026-12-15', 2, 'OP-1002')`); err != nil {
        return fmt.Errorf("seed vehicles: %w", err)
    }

    if _, err = tx.Exec(`INSERT OR IGNORE INTO load (id, date, load_number, loader_id, well_id, ticket_number, miles, driver_id, vec_id, trailer_number, net_weight, tons, ton_rate, total_value, status, ticket_id) VALUES
        (1, '2026-09-15', 'LD-1001', 1, 1, 'TKT-1001', 142.5, 1, 'VEC-1001', 24, 21000, 10.5, 185.25, 1946.62, 1, 'TICKET-1001'),
        (2, '2026-09-18', 'LD-1002', 2, 2, 'TKT-1002', 138.0, 2, 'VEC-1002', 18, 19850, 9.925, 180.50, 1791.00, 0, 'TICKET-1002')`); err != nil {
        return fmt.Errorf("seed loads: %w", err)
    }

    if _, err = tx.Exec(`INSERT OR IGNORE INTO category (id, type) VALUES
        (1, 'Fuel'),
        (2, 'Maintenance'),
        (3, 'Operations')`); err != nil {
        return fmt.Errorf("seed categories: %w", err)
    }

    if _, err = tx.Exec(`INSERT OR IGNORE INTO expense (id, date, concept, category_id, amount, quantity, total, receipt_id) VALUES
        (1, '2026-09-10', 'Diesel fill', 1, 120.00, 40.00, 120.00, 1),
        (2, '2026-09-12', 'Brake service', 2, 450.00, 1.00, 450.00, 2)`); err != nil {
        return fmt.Errorf("seed expenses: %w", err)
    }

    if _, err = tx.Exec(`INSERT OR IGNORE INTO income (id, date, concept, category_id, amount, quantity, total, receipt_id) VALUES
        (1, '2026-09-15', 'Load payment', 3, 1946.62, 1.00, 1946.62, 1001),
        (2, '2026-09-18', 'Load payment', 3, 1791.00, 1.00, 1791.00, 1002)`); err != nil {
        return fmt.Errorf("seed incomes: %w", err)
    }

    if _, err = tx.Exec(`INSERT OR IGNORE INTO fuel_expense (id, id_expense, id_driver) VALUES
        (1, 1, 1),
        (2, 2, 2)`); err != nil {
        return fmt.Errorf("seed fuel expenses: %w", err)
    }

    if _, err = tx.Exec(`INSERT OR IGNORE INTO stats (id, number_week, date_range_init, date_range_end, expense_total, income_total, net_total) VALUES
        (1, 37, '2026-09-13', '2026-09-19', 570.00, 3737.62, 3167.62)`); err != nil {
        return fmt.Errorf("seed stats: %w", err)
    }

    if err = tx.Commit(); err != nil {
        return fmt.Errorf("commit seed transaction: %w", err)
    }

    return nil
}

// SeedFixtures loads a JSON fixture payload into the database for testing or local development.
func SeedFixtures(db *sql.DB, payload string) error {
    if db == nil {
        return fmt.Errorf("%w: database is nil", ErrValidation)
    }
    if strings.TrimSpace(payload) == "" {
        return fmt.Errorf("%w: fixture payload is empty", ErrValidation)
    }

    var fixtures map[string]json.RawMessage
    if err := json.Unmarshal([]byte(payload), &fixtures); err != nil {
        return fmt.Errorf("parse fixture payload: %w", err)
    }

    // Seed the most common fixture structures in a neutral, idempotent way.
    for table, items := range fixtures {
        switch table {
        case "drivers":
            if err := seedFixtureRows(db, table, items, "driver", []string{"id", "name", "license_number", "phone", "email", "address"}); err != nil {
                return err
            }
        case "loaders":
            if err := seedFixtureRows(db, table, items, "loader", []string{"id", "name", "phone", "email", "address", "company"}); err != nil {
                return err
            }
        case "vehicles":
            if err := seedFixtureRows(db, table, items, "vehicle", []string{"id", "vec_id", "license_id", "trailer_number", "vin", "license_expiration", "driver_id", "overweight_permit_id"}); err != nil {
                return err
            }
        case "wells":
            if err := seedFixtureRows(db, table, items, "well", []string{"id", "name", "location"}); err != nil {
                return err
            }
        case "loads":
            if err := seedFixtureRows(db, table, items, "load", []string{"id", "date", "load_number", "loader_id", "well_id", "ticket_number", "miles", "driver_id", "vec_id", "trailer_number", "net_weight", "tons", "ton_rate", "total_value", "status", "ticket_id"}); err != nil {
                return err
            }
        }
    }

    return nil
}

func seedFixtureRows(db *sql.DB, fixtureName string, payload json.RawMessage, table string, columns []string) error {
    var rows []map[string]any
    if err := json.Unmarshal(payload, &rows); err != nil {
        return fmt.Errorf("parse fixture %s: %w", fixtureName, err)
    }

    for _, row := range rows {
        vals := make([]any, 0, len(columns))
        placeholders := make([]string, 0, len(columns))
        for _, col := range columns {
            vals = append(vals, row[col])
            placeholders = append(placeholders, "?")
        }

        insertSQL := fmt.Sprintf(`INSERT OR IGNORE INTO %s (%s) VALUES (%s)`, table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))
        if _, err := db.Exec(insertSQL, vals...); err != nil {
            return fmt.Errorf("insert fixture row into %s: %w", table, err)
        }
    }

    return nil
}
