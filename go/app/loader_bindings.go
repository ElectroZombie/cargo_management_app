package app

import "cargo_management_app/go/database"

// ---------- Loader methods ----------

// ListLoaders returns a filtered, sorted, paginated page of loaders.
func (a *App) ListLoaders(filters database.LoaderFilters, opts database.QueryOptions) Response {
	loaders, total, err := a.db.ListLoaders(filters, opts)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: loaders, Total: total}
}

// GetLoader returns a single loader by id.
func (a *App) GetLoader(id int) Response {
	loader, err := a.db.GetLoader(id)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: loader}
}

// CreateLoader inserts a new loader.
func (a *App) CreateLoader(dto database.LoaderDTO) Response {
	loader, err := a.db.CreateLoader(dto)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: loader}
}

// UpdateLoader modifies an existing loader.
func (a *App) UpdateLoader(id int, dto database.LoaderDTO) Response {
	if err := a.db.UpdateLoader(id, dto); err != nil {
		return failure(err)
	}
	loader, err := a.db.GetLoader(id)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: loader}
}

// DeleteLoader removes a loader.
func (a *App) DeleteLoader(id int) Response {
	if err := a.db.DeleteLoader(id); err != nil {
		return failure(err)
	}
	return Response{Success: true}
}

// SearchLoaders performs a fuzzy search across loader fields.
func (a *App) SearchLoaders(query string) Response {
	items, err := a.db.SearchLoaders(query)
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: items, Total: int64(len(items))}
}
