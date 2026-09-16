package main

import (
	"database/sql"
	"fmt"
	"log"
)

// ============ WELL CRUD OPERATIONS ============

// CreateWell adds a new well to the database
func (d *Database) CreateWell(name, client, location string) (*Well, error) {
	query := `
	INSERT INTO well (name, client, location, created_at, updated_at)
	VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := d.db.Exec(query, name, client, location)
	if err != nil {
		return nil, fmt.Errorf("failed to create well: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	well := &Well{
		ID:       int(id),
		Name:     name,
		Client:   client,
		Location: location,
	}

	log.Printf("Well created successfully: %s (ID: %d)", name, id)
	return well, nil
}

// GetAllWells retrieves all wells from the database
func (d *Database) GetAllWells() ([]Well, error) {
	query := `SELECT id, name, client, location, created_at, updated_at FROM well`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query wells: %w", err)
	}
	defer rows.Close()

	var wells []Well
	for rows.Next() {
		var well Well
		if err := rows.Scan(&well.ID, &well.Name, &well.Client, &well.Location, &well.CreatedAt, &well.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan well: %w", err)
		}
		wells = append(wells, well)
	}

	return wells, nil
}

// GetWellByID retrieves a specific well by ID
func (d *Database) GetWellByID(id int) (*Well, error) {
	query := `SELECT id, name, client, location, created_at, updated_at FROM well WHERE id = ?`

	var well Well
	err := d.db.QueryRow(query, id).Scan(&well.ID, &well.Name, &well.Client, &well.Location, &well.CreatedAt, &well.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("well not found")
		}
		return nil, fmt.Errorf("failed to query well: %w", err)
	}

	return &well, nil
}

// UpdateWell updates an existing well
func (d *Database) UpdateWell(id int, name, client, location string) error {
	query := `
	UPDATE well 
	SET name = ?, client = ?, location = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	result, err := d.db.Exec(query, name, client, location, id)
	if err != nil {
		return fmt.Errorf("failed to update well: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("well not found")
	}

	log.Printf("Well updated successfully: ID %d", id)
	return nil
}

// DeleteWell removes a well from the database
func (d *Database) DeleteWell(id int) error {
	query := `DELETE FROM well WHERE id = ?`

	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete well: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("well not found")
	}

	log.Printf("Well deleted successfully: ID %d", id)
	return nil
}
