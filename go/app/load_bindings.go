package app

import "cargo_management_app/go/database"

// ---------- Load methods ----------

// ListLoads returns a filtered, sorted, paginated page of loads.
func (a *App) ListLoads(filters database.LoadFilters, opts database.QueryOptions) Response {
	loads, total, err := a.db.ListLoads(filters, opts)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: loads, Total: total}
}

// GetLoad returns a single load by id.
func (a *App) GetLoad(id int) Response {
	load, err := a.db.GetLoad(id)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: load}
}

// CreateLoad inserts a new load.
func (a *App) CreateLoad(dto database.LoadDTO) Response {
	load, err := a.db.CreateLoad(dto)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: load}
}

// UpdateLoad modifies an existing load.
func (a *App) UpdateLoad(id int, dto database.LoadDTO) Response {
	if err := a.db.UpdateLoad(id, dto); err != nil {
		return failure(err)
	}
	load, err := a.db.GetLoad(id)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: load}
}

// DeleteLoad removes a load.
func (a *App) DeleteLoad(id int) Response {
	if err := a.db.DeleteLoad(id); err != nil {
		return failure(err)
	}
	return Response{Success: true}
}

// SearchLoads performs a fuzzy search across load fields.
func (a *App) SearchLoads(query string) Response {
	items, err := a.db.SearchLoads(query)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: items, Total: int64(len(items))}
}

// GetLoadsByDriver returns every load assigned to a driver.
func (a *App) GetLoadsByDriver(driverID int) Response {
	loads, err := a.db.GetLoadsByDriver(driverID)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: loads, Total: int64(len(loads))}
}

// GetLoadsByWell returns every load assigned to a well.
func (a *App) GetLoadsByWell(wellID int) Response {
	loads, err := a.db.GetLoadsByWell(wellID)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: loads, Total: int64(len(loads))}
}

// GetLoadsByStatus returns loads filtered by active/inactive status.
func (a *App) GetLoadsByStatus(status bool) Response {
	loads, err := a.db.GetLoadsByStatus(status)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: loads, Total: int64(len(loads))}
}
