package app

import "cargo_management_app/go/database"

// ---------- Driver methods ----------

// ListDrivers returns a filtered, sorted, paginated page of drivers.
func (a *App) ListDrivers(filters database.DriverFilters, opts database.QueryOptions) Response {
	drivers, total, err := a.db.ListDrivers(filters, opts)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: drivers, Total: total}
}

// GetDriver returns a single driver by id.
func (a *App) GetDriver(id int) Response {
	driver, err := a.db.GetDriver(id)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: driver}
}

// CreateDriver inserts a new driver.
func (a *App) CreateDriver(dto database.DriverDTO) Response {
	driver, err := a.db.CreateDriver(dto)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: driver}
}

// UpdateDriver modifies an existing driver.
func (a *App) UpdateDriver(id int, dto database.DriverDTO) Response {
	if err := a.db.UpdateDriver(id, dto); err != nil {
		return failure(err)
	}
	driver, err := a.db.GetDriver(id)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: driver}
}

// DeleteDriver removes a driver.
func (a *App) DeleteDriver(id int) Response {
	if err := a.db.DeleteDriver(id); err != nil {
		return failure(err)
	}
	return Response{Success: true}
}

// SearchDrivers performs a fuzzy search across driver fields.
func (a *App) SearchDrivers(query string) Response {
	items, err := a.db.SearchDrivers(query)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: items, Total: int64(len(items))}
}
