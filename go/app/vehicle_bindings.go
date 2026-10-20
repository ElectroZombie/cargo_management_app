package app

import "cargo_management_app/go/database"

// ---------- Vehicle methods ----------

// ListVehicles returns a filtered, sorted, paginated page of vehicles.
func (a *App) ListVehicles(filters database.VehicleFilters, opts database.QueryOptions) Response {
	vehicles, total, err := a.db.ListVehicles(filters, opts)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: vehicles, Total: total}
}

// GetVehicle returns a single vehicle by id.
func (a *App) GetVehicle(id int) Response {
	vehicle, err := a.db.GetVehicle(id)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: vehicle}
}

// CreateVehicle inserts a new vehicle.
func (a *App) CreateVehicle(dto database.VehicleDTO) Response {
	vehicle, err := a.db.CreateVehicle(dto)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: vehicle}
}

// UpdateVehicle modifies an existing vehicle.
func (a *App) UpdateVehicle(id int, dto database.VehicleDTO) Response {
	if err := a.db.UpdateVehicle(id, dto); err != nil {
		return failure(err)
	}
	vehicle, err := a.db.GetVehicle(id)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: vehicle}
}

// DeleteVehicle removes a vehicle.
func (a *App) DeleteVehicle(id int) Response {
	if err := a.db.DeleteVehicle(id); err != nil {
		return failure(err)
	}
	return Response{Success: true}
}

// SearchVehicles performs a fuzzy search across vehicle fields.
func (a *App) SearchVehicles(query string) Response {
	items, err := a.db.SearchVehicles(query)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: items, Total: int64(len(items))}
}
