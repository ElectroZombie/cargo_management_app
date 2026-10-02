package database

import (
	"database/sql"
	"fmt"
	"sort"
	"time"
)

const schemaMigrationsTableSQL = `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
`

// Migration represents a single versioned database migration.
type Migration struct {
	Version string
	Name    string
	Up      func(*sql.Tx) error
	Down    func(*sql.Tx) error
}

// MigrationManager tracks and applies schema changes with rollback support.
type MigrationManager struct {
	db *sql.DB
}

func NewMigrationManager(db *sql.DB) *MigrationManager {
	return &MigrationManager{db: db}
}

// DefaultMigrations defines the current versioned migration set for the app.
func DefaultMigrations() []Migration {
	return []Migration{
		{
			Version: "2026-10-02_001_initial_schema",
			Name:    "initial_schema",
			Up: func(tx *sql.Tx) error {
				for _, statement := range schemaStatements() {
					if _, err := tx.Exec(statement); err != nil {
						return fmt.Errorf("execute schema statement: %w", err)
					}
				}
				return nil
			},
			Down: func(tx *sql.Tx) error {
				for _, statement := range reverseSchemaDropStatements() {
					if _, err := tx.Exec(statement); err != nil {
						return fmt.Errorf("rollback schema statement: %w", err)
					}
				}
				return nil
			},
		},
	}
}

func reverseSchemaDropStatements() []string {
	return []string{
		`DROP TABLE IF EXISTS fuel_expense`,
		`DROP TABLE IF EXISTS expense`,
		`DROP TABLE IF EXISTS income`,
		`DROP TABLE IF EXISTS category`,
		`DROP TABLE IF EXISTS trip`,
		`DROP TABLE IF EXISTS trailer`,
		`DROP TABLE IF EXISTS truck`,
		`DROP TABLE IF EXISTS well`,
		`DROP TABLE IF EXISTS driver`,
	}
}

func ApplyMigrations(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("%w: database is nil", ErrValidation)
	}
	return NewMigrationManager(db).Apply(DefaultMigrations())
}

func RollbackMigration(db *sql.DB, version string) error {
	if db == nil {
		return fmt.Errorf("%w: database is nil", ErrValidation)
	}
	return NewMigrationManager(db).Rollback(version, DefaultMigrations())
}

func (m *MigrationManager) Apply(migrations []Migration) error {
	if m == nil || m.db == nil {
		return fmt.Errorf("%w: migration manager has no database", ErrValidation)
	}
	if err := m.ensureMigrationsTable(); err != nil {
		return fmt.Errorf("ensure migration table: %w", err)
	}

	ordered := append([]Migration(nil), migrations...)
	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].Version < ordered[j].Version
	})

	for _, migration := range ordered {
		applied, err := m.isApplied(migration.Version)
		if err != nil {
			return fmt.Errorf("check migration %s: %w", migration.Version, err)
		}
		if applied {
			continue
		}

		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("begin migration %s: %w", migration.Version, err)
		}

		if migration.Up != nil {
			if err := migration.Up(tx); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("apply migration %s (%s): %w", migration.Version, migration.Name, err)
			}
		}

		_, err = tx.Exec(
			`INSERT INTO schema_migrations(version, name, applied_at) VALUES (?, ?, ?)`,
			migration.Version,
			migration.Name,
			time.Now().UTC().Format(time.RFC3339),
		)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", migration.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", migration.Version, err)
		}
	}

	return nil
}

func (m *MigrationManager) Rollback(version string, migrations []Migration) error {
	if m == nil || m.db == nil {
		return fmt.Errorf("%w: migration manager has no database", ErrValidation)
	}
	if err := m.ensureMigrationsTable(); err != nil {
		return fmt.Errorf("ensure migration table: %w", err)
	}

	var target *Migration
	for i := range migrations {
		if migrations[i].Version == version {
			target = &migrations[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("%w: migration %s not found", ErrNotFound, version)
	}

	applied, err := m.isApplied(version)
	if err != nil {
		return fmt.Errorf("check migration %s: %w", version, err)
	}
	if !applied {
		return fmt.Errorf("%w: migration %s is not applied", ErrNotFound, version)
	}

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("begin rollback for %s: %w", version, err)
	}

	if target.Down != nil {
		if err := target.Down(tx); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("rollback migration %s (%s): %w", version, target.Name, err)
		}
	}

	if _, err := tx.Exec(`DELETE FROM schema_migrations WHERE version = ?`, version); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("remove migration record %s: %w", version, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit rollback for %s: %w", version, err)
	}

	return nil
}

func (m *MigrationManager) ensureMigrationsTable() error {
	_, err := m.db.Exec(schemaMigrationsTableSQL)
	return err
}

func (m *MigrationManager) isApplied(version string) (bool, error) {
	var present int
	err := m.db.QueryRow(`SELECT 1 FROM schema_migrations WHERE version = ? LIMIT 1`, version).Scan(&present)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
