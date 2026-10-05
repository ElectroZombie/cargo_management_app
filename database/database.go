package database

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

type Database struct { db *sql.DB }

func NewDatabase(path string) (*Database, error) {
    db, err := sql.Open("sqlite3", path)
    if err != nil { return nil, fmt.Errorf("open database: %w", err) }
    db.SetMaxOpenConns(1)
    if err := db.Ping(); err != nil { db.Close(); return nil, fmt.Errorf("ping database: %w", err) }
    d := &Database{db: db}
    if err := d.createSchema(); err != nil { db.Close(); return nil, err }
    return d, nil
}

func (d *Database) Close() error { if d == nil || d.db == nil { return nil }; return d.db.Close() }

func (d *Database) createSchema() error {
    if _, err := d.db.Exec(`PRAGMA foreign_keys = ON`); err != nil { return fmt.Errorf("enable foreign keys: %w", err) }
    for _, statement := range schemaStatements() {
        if _, err := d.db.Exec(statement); err != nil { return fmt.Errorf("create schema: %w", err) }
    }
    return nil
}
