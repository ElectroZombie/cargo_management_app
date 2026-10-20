package app

// ---------- Dashboard methods ----------

// GetDashboardStats returns the aggregated operational snapshot.
func (a *App) GetDashboardStats() Response {
	stats, err := a.db.GetDashboardStats()
	if err != nil {
		return failure(err)
	}
	return Response{Success: true, Data: stats}
}
