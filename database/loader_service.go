package main

import (
	"database/sql"
	"fmt"
	"log"
)

// ============ LOADER CRUD OPERATIONS ============

// CreateLoader adds a new loader to the database
func (d *Database) CreateLoader(name, phone, email, address, company string) (*Loader, error) {
	query := `
	INSERT INTO loader (name, phone, email, address, company, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := d.db.Exec(query, name, phone, email, address, company)
	if err != nil {
		return nil, fmt.Errorf("failed to create loader: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	loader := &Loader{
		ID:      int(id),
		Name:    name,
		Phone:   phone,
		Email:   email,
		Address: address,
		Company: company,
	}

	log.Printf("Loader created successfully: %s (ID: %d)", name, id)
	return loader, nil
}

// GetAllLoaders retrieves all loaders from the database
func (d *Database) GetAllLoaders() ([]Loader, error) {
	query := `SELECT id, name, phone, email, address, company, created_at, updated_at FROM loader`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query loaders: %w", err)
	}
	defer rows.Close()

	var loaders []Loader
	for rows.Next() {
		var loader Loader
		if err := rows.Scan(&loader.ID, &loader.Name, &loader.Phone, &loader.Email, &loader.Address, &loader.Company, &loader.CreatedAt, &loader.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan loader: %w", err)
		}
		loaders = append(loaders, loader)
	}

	return loaders, nil
}

// GetLoaderByID retrieves a specific loader by ID
func (d *Database) GetLoaderByID(id int) (*Loader, error) {
	query := `SELECT id, name, phone, email, address, company, created_at, updated_at FROM loader WHERE id = ?`

	var loader Loader
	err := d.db.QueryRow(query, id).Scan(&loader.ID, &loader.Name, &loader.Phone, &loader.Email, &loader.Address, &loader.Company, &loader.CreatedAt, &loader.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("loader not found")
		}
		return nil, fmt.Errorf("failed to query loader: %w", err)
	}

	return &loader, nil
}

// UpdateLoader updates an existing loader
func (d *Database) UpdateLoader(id int, name, phone, email, address, company string) error {
	query := `
	UPDATE loader 
	SET name = ?, phone = ?, email = ?, address = ?, company = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	result, err := d.db.Exec(query, name, phone, email, address, company, id)
	if err != nil {
		return fmt.Errorf("failed to update loader: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("loader not found")
	}

	log.Printf("Loader updated successfully: ID %d", id)
	return nil
}

// DeleteLoader removes a loader from the database
func (d *Database) DeleteLoader(id int) error {
	query := `DELETE FROM loader WHERE id = ?`

	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete loader: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("loader not found")
	}

	log.Printf("Loader deleted successfully: ID %d", id)
	return nil
}
