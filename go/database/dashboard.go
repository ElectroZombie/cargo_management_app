package database

// DashboardStats is the aggregated operational snapshot shown on the dashboard.
type DashboardStats struct {
	TotalLoads     int     `json:"total_loads"`
	ActiveLoads    int     `json:"active_loads"`
	TotalNetWeight float64 `json:"total_net_weight"`
	TotalTons      float64 `json:"total_tons"`
	TotalValue     float64 `json:"total_value"`
	DriverCount    int     `json:"driver_count"`
	LoaderCount    int     `json:"loader_count"`
	VehicleCount   int     `json:"vehicle_count"`
	WellCount      int     `json:"well_count"`
	RecentLoads    []Load  `json:"recent_loads"`
}

// GetDashboardStats aggregates counts, load totals, and the most recent loads.
func (d *Database) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{RecentLoads: []Load{}}

	if err := d.runner.QueryRow(`
		SELECT COUNT(*),
			COALESCE(SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(net_weight), 0),
			COALESCE(SUM(tons), 0),
			COALESCE(SUM(total_value), 0)
		FROM load`).Scan(
		&stats.TotalLoads,
		&stats.ActiveLoads,
		&stats.TotalNetWeight,
		&stats.TotalTons,
		&stats.TotalValue,
	); err != nil {
		return nil, databaseError("read dashboard load totals", err)
	}

	counts := []struct {
		table string
		dest  *int
	}{
		{"driver", &stats.DriverCount},
		{"loader", &stats.LoaderCount},
		{"vehicle", &stats.VehicleCount},
		{"well", &stats.WellCount},
	}
	for _, c := range counts {
		if err := d.runner.QueryRow("SELECT COUNT(*) FROM " + c.table).Scan(c.dest); err != nil {
			return nil, databaseError("read dashboard count for "+c.table, err)
		}
	}

	recent, _, err := d.ListLoads(LoadFilters{}, QueryOptions{Page: 1, PageSize: 5, SortBy: "date", SortDir: "DESC"})
	if err != nil {
		return nil, err
	}
	if recent != nil {
		stats.RecentLoads = recent
	}

	return stats, nil
}
