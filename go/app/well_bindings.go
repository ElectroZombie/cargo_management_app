package app

import "cargo_management_app/go/database"

// ---------- Well methods ----------

// ListWells returns a filtered, sorted, paginated page of wells.
func (a *App) ListWells(filters database.WellFilters, opts database.QueryOptions) Response {
	wells, total, err := a.db.ListWells(filters, opts)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: wells, Total: total}
}

// GetWell returns a single well by id.
func (a *App) GetWell(id int) Response {
	well, err := a.db.GetWell(id)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: well}
}

// CreateWell inserts a new well.
func (a *App) CreateWell(dto database.WellDTO) Response {
	well, err := a.db.CreateWell(dto)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: well}
}

// UpdateWell modifies an existing well.
func (a *App) UpdateWell(id int, dto database.WellDTO) Response {
	if err := a.db.UpdateWell(id, dto); err != nil {
		return failure(err)
	}
	well, err := a.db.GetWell(id)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: well}
}

// DeleteWell removes a well.
func (a *App) DeleteWell(id int) Response {
	if err := a.db.DeleteWell(id); err != nil {
		return failure(err)
	}
	return Response{Success: true}
}

// SearchWells performs a fuzzy search across well fields.
func (a *App) SearchWells(query string) Response {
	items, err := a.db.SearchWells(query)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: items, Total: int64(len(items))}
}
