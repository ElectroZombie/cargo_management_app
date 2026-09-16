package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// ============ DRIVER CRUD OPERATIONS ============

// CreateDriver adds a new driver to the database
func (d *Database) CreateDriver(name, license, phone, email, address string) (*Driver, error) {
	query := `
	INSERT INTO driver (name, license_number, phone, email, address, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := d.db.Exec(query, name, license, phone, email, address)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	driver := &Driver{
		ID:      int(id),
		Name:    name,
		License: license,
		Phone:   phone,
		Email:   email,
		Address: address,
	}

	log.Printf("Driver created successfully: %s (ID: %d)", name, id)
	return driver, nil
}

// GetAllDrivers retrieves all drivers from the database
func (d *Database) GetAllDrivers() ([]Driver, error) {
	query := `SELECT id, name, license_number, phone, email, address, created_at, updated_at FROM driver`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query drivers: %w", err)
	}
	defer rows.Close()

	var drivers []Driver
	for rows.Next() {
		var driver Driver
		if err := rows.Scan(&driver.ID, &driver.Name, &driver.License, &driver.Phone, &driver.Email, &driver.Address, &driver.CreatedAt, &driver.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan driver: %w", err)
		}
		drivers = append(drivers, driver)
	}

	return drivers, nil
}

// GetDriverByID retrieves a specific driver by ID
func (d *Database) GetDriverByID(id int) (*Driver, error) {
	query := `SELECT id, name, license_number, phone, email, address, created_at, updated_at FROM driver WHERE id = ?`

	var driver Driver
	err := d.db.QueryRow(query, id).Scan(&driver.ID, &driver.Name, &driver.License, &driver.Phone, &driver.Email, &driver.Address, &driver.CreatedAt, &driver.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("driver not found")
		}
		return nil, fmt.Errorf("failed to query driver: %w", err)
	}

	return &driver, nil
}

// UpdateDriver updates an existing driver
func (d *Database) UpdateDriver(id int, name, license, phone, email, address string) error {
	query := `
	UPDATE driver 
	SET name = ?, license_number = ?, phone = ?, email = ?, address = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	result, err := d.db.Exec(query, name, license, phone, email, address, id)
	if err != nil {
		return fmt.Errorf("failed to update driver: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("driver not found")
	}

	log.Printf("Driver updated successfully: ID %d", id)
	return nil
}

// DeleteDriver removes a driver from the database
func (d *Database) DeleteDriver(id int) error {
	query := `DELETE FROM driver WHERE id = ?`

	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete driver: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("driver not found")
	}

	log.Printf("Driver deleted successfully: ID %d", id)
	return nil
}

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

// ============ VEHICLE CRUD OPERATIONS ============

// CreateVehicle adds a new vehicle to the database
func (d *Database) CreateVehicle(vecID, licenseID string, trailerNumber int, vin string, licenseExpiration time.Time, driverID int, overweightPermitID string) (*Vehicle, error) {
	query := `
	INSERT INTO vehicle (vec_id, license_id, trailer_number, vin, license_expiration, driver_id, overweight_permit_id, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := d.db.Exec(query, vecID, licenseID, trailerNumber, vin, licenseExpiration, driverID, overweightPermitID)
	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	vehicle := &Vehicle{
		ID:                 int(id),
		VecID:              vecID,
		LicenseID:          licenseID,
		TrailerNumber:      trailerNumber,
		VIN:                vin,
		LicenseExpiration:  licenseExpiration,
		DriverID:           driverID,
		OverweightPermitID: overweightPermitID,
	}

	log.Printf("Vehicle created successfully: %s (ID: %d)", vecID, id)
	return vehicle, nil
}

// GetAllVehicles retrieves all vehicles from the database
func (d *Database) GetAllVehicles() ([]Vehicle, error) {
	query := `SELECT id, vec_id, license_id, trailer_number, vin, license_expiration, driver_id, overweight_permit_id, created_at, updated_at FROM vehicle`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var vehicle Vehicle
		if err := rows.Scan(&vehicle.ID, &vehicle.VecID, &vehicle.LicenseID, &vehicle.TrailerNumber, &vehicle.VIN, &vehicle.LicenseExpiration, &vehicle.DriverID, &vehicle.OverweightPermitID, &vehicle.CreatedAt, &vehicle.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}
		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

// GetVehicleByID retrieves a specific vehicle by ID
func (d *Database) GetVehicleByID(id int) (*Vehicle, error) {
	query := `SELECT id, vec_id, license_id, trailer_number, vin, license_expiration, driver_id, overweight_permit_id, created_at, updated_at FROM vehicle WHERE id = ?`

	var vehicle Vehicle
	err := d.db.QueryRow(query, id).Scan(&vehicle.ID, &vehicle.VecID, &vehicle.LicenseID, &vehicle.TrailerNumber, &vehicle.VIN, &vehicle.LicenseExpiration, &vehicle.DriverID, &vehicle.OverweightPermitID, &vehicle.CreatedAt, &vehicle.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle not found")
		}
		return nil, fmt.Errorf("failed to query vehicle: %w", err)
	}

	return &vehicle, nil
}

// GetVehicleByVecID retrieves a specific vehicle by vec_id
func (d *Database) GetVehicleByVecID(vecID string) (*Vehicle, error) {
	query := `SELECT id, vec_id, license_id, trailer_number, vin, license_expiration, driver_id, overweight_permit_id, created_at, updated_at FROM vehicle WHERE vec_id = ?`

	var vehicle Vehicle
	err := d.db.QueryRow(query, vecID).Scan(&vehicle.ID, &vehicle.VecID, &vehicle.LicenseID, &vehicle.TrailerNumber, &vehicle.VIN, &vehicle.LicenseExpiration, &vehicle.DriverID, &vehicle.OverweightPermitID, &vehicle.CreatedAt, &vehicle.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle not found")
		}
		return nil, fmt.Errorf("failed to query vehicle: %w", err)
	}

	return &vehicle, nil
}

// UpdateVehicle updates an existing vehicle
func (d *Database) UpdateVehicle(id int, vecID, licenseID string, trailerNumber int, vin string, licenseExpiration time.Time, driverID int, overweightPermitID string) error {
	query := `
	UPDATE vehicle 
	SET vec_id = ?, license_id = ?, trailer_number = ?, vin = ?, license_expiration = ?, driver_id = ?, overweight_permit_id = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	result, err := d.db.Exec(query, vecID, licenseID, trailerNumber, vin, licenseExpiration, driverID, overweightPermitID, id)
	if err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle not found")
	}

	log.Printf("Vehicle updated successfully: ID %d", id)
	return nil
}

// DeleteVehicle removes a vehicle from the database
func (d *Database) DeleteVehicle(id int) error {
	query := `DELETE FROM vehicle WHERE id = ?`

	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle not found")
	}

	log.Printf("Vehicle deleted successfully: ID %d", id)
	return nil
}

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

// ============ LOAD CRUD OPERATIONS ============

// CreateLoad adds a new load to the database
func (d *Database) CreateLoad(date time.Time, loadNumber string, loaderID int, wellID int, ticketNumber string, miles float64, driverID int, vecID string, trailerNumber int, netWeight, tons, tonRate, totalValue float64, status bool, ticketID string) (*Load, error) {
	query := `
	INSERT INTO load (date, load_number, loader_id, well_id, ticket_number, miles, driver_id, vec_id, trailer_number, net_weight, tons, ton_rate, total_value, status, ticket_id, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := d.db.Exec(query, date, loadNumber, loaderID, wellID, ticketNumber, miles, driverID, vecID, trailerNumber, netWeight, tons, tonRate, totalValue, status, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to create load: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	load := &Load{
		ID:            int(id),
		Date:          date,
		LoadNumber:    loadNumber,
		LoaderID:      loaderID,
		WellID:        wellID,
		TicketNumber:  ticketNumber,
		Miles:         miles,
		DriverID:      driverID,
		VecID:         vecID,
		TrailerNumber: trailerNumber,
		NetWeight:     netWeight,
		Tons:          tons,
		TonRate:       tonRate,
		TotalValue:    totalValue,
		Status:        status,
		TicketID:      ticketID,
	}

	log.Printf("Load created successfully: %s (ID: %d)", loadNumber, id)
	return load, nil
}

// GetAllLoads retrieves all loads from the database
func (d *Database) GetAllLoads() ([]Load, error) {
	query := `SELECT id, date, load_number, loader_id, well_id, ticket_number, miles, driver_id, vec_id, trailer_number, net_weight, tons, ton_rate, total_value, status, ticket_id, created_at, updated_at FROM load`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query loads: %w", err)
	}
	defer rows.Close()

	var loads []Load
	for rows.Next() {
		var load Load
		if err := rows.Scan(&load.ID, &load.Date, &load.LoadNumber, &load.LoaderID, &load.WellID, &load.TicketNumber, &load.Miles, &load.DriverID, &load.VecID, &load.TrailerNumber, &load.NetWeight, &load.Tons, &load.TonRate, &load.TotalValue, &load.Status, &load.TicketID, &load.CreatedAt, &load.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan load: %w", err)
		}
		loads = append(loads, load)
	}

	return loads, nil
}

// GetLoadByID retrieves a specific load by ID
func (d *Database) GetLoadByID(id int) (*Load, error) {
	query := `SELECT id, date, load_number, loader_id, well_id, ticket_number, miles, driver_id, vec_id, trailer_number, net_weight, tons, ton_rate, total_value, status, ticket_id, created_at, updated_at FROM load WHERE id = ?`

	var load Load
	err := d.db.QueryRow(query, id).Scan(&load.ID, &load.Date, &load.LoadNumber, &load.LoaderID, &load.WellID, &load.TicketNumber, &load.Miles, &load.DriverID, &load.VecID, &load.TrailerNumber, &load.NetWeight, &load.Tons, &load.TonRate, &load.TotalValue, &load.Status, &load.TicketID, &load.CreatedAt, &load.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("load not found")
		}
		return nil, fmt.Errorf("failed to query load: %w", err)
	}

	return &load, nil
}

// UpdateLoad updates an existing load
func (d *Database) UpdateLoad(id int, date time.Time, loadNumber string, loaderID int, wellID int, ticketNumber string, miles float64, driverID int, vecID string, trailerNumber int, netWeight, tons, tonRate, totalValue float64, status bool, ticketID string) error {
	query := `
	UPDATE load 
	SET date = ?, load_number = ?, loader_id = ?, well_id = ?, ticket_number = ?, miles = ?, driver_id = ?, vec_id = ?, trailer_number = ?, net_weight = ?, tons = ?, ton_rate = ?, total_value = ?, status = ?, ticket_id = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	result, err := d.db.Exec(query, date, loadNumber, loaderID, wellID, ticketNumber, miles, driverID, vecID, trailerNumber, netWeight, tons, tonRate, totalValue, status, ticketID, id)
	if err != nil {
		return fmt.Errorf("failed to update load: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("load not found")
	}

	log.Printf("Load updated successfully: ID %d", id)
	return nil
}

// DeleteLoad removes a load from the database
func (d *Database) DeleteLoad(id int) error {
	query := `DELETE FROM load WHERE id = ?`

	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete load: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("load not found")
	}

	log.Printf("Load deleted successfully: ID %d", id)
	return nil
}

// GetLoadsByDriverID retrieves all loads for a specific driver
func (d *Database) GetLoadsByDriverID(driverID int) ([]Load, error) {
	query := `SELECT id, date, load_number, loader_id, well_id, ticket_number, miles, driver_id, vec_id, trailer_number, net_weight, tons, ton_rate, total_value, status, ticket_id, created_at, updated_at FROM load WHERE driver_id = ?`

	rows, err := d.db.Query(query, driverID)
	if err != nil {
		return nil, fmt.Errorf("failed to query loads: %w", err)
	}
	defer rows.Close()

	var loads []Load
	for rows.Next() {
		var load Load
		if err := rows.Scan(&load.ID, &load.Date, &load.LoadNumber, &load.LoaderID, &load.WellID, &load.TicketNumber, &load.Miles, &load.DriverID, &load.VecID, &load.TrailerNumber, &load.NetWeight, &load.Tons, &load.TonRate, &load.TotalValue, &load.Status, &load.TicketID, &load.CreatedAt, &load.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan load: %w", err)
		}
		loads = append(loads, load)
	}

	return loads, nil
}

// GetLoadsByWellID retrieves all loads for a specific well
func (d *Database) GetLoadsByWellID(wellID int) ([]Load, error) {
	query := `SELECT id, date, load_number, loader_id, well_id, ticket_number, miles, driver_id, vec_id, trailer_number, net_weight, tons, ton_rate, total_value, status, ticket_id, created_at, updated_at FROM load WHERE well_id = ?`

	rows, err := d.db.Query(query, wellID)
	if err != nil {
		return nil, fmt.Errorf("failed to query loads: %w", err)
	}
	defer rows.Close()

	var loads []Load
	for rows.Next() {
		var load Load
		if err := rows.Scan(&load.ID, &load.Date, &load.LoadNumber, &load.LoaderID, &load.WellID, &load.TicketNumber, &load.Miles, &load.DriverID, &load.VecID, &load.TrailerNumber, &load.NetWeight, &load.Tons, &load.TonRate, &load.TotalValue, &load.Status, &load.TicketID, &load.CreatedAt, &load.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan load: %w", err)
		}
		loads = append(loads, load)
	}

	return loads, nil
}
