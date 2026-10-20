package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// querier is the subset of database/sql used by the service layer.
// Both *sql.DB and *sql.Tx satisfy it, which lets transaction helpers
// run service calls inside a transaction without duplicating logic.
type querier interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type Database struct {
	db     *sql.DB
	runner querier
}

func NewDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	d := &Database{db: db, runner: db}
	if err := d.createSchema(); err != nil {
		db.Close()
		return nil, err
	}
	return d, nil
}

func (d *Database) Close() error {
	if d == nil || d.db == nil {
		return nil
	}
	return d.db.Close()
}

func (d *Database) createSchema() error {
	if _, err := d.runner.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		return fmt.Errorf("enable foreign keys: %w", err)
	}
	for _, statement := range schemaStatements() {
		if _, err := d.runner.Exec(statement); err != nil {
			return fmt.Errorf("create schema: %w", err)
		}
	}
	if err := d.createIndexes(); err != nil {
		return err
	}
	return nil
}