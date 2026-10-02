package database

import (
	"database/sql"
	"fmt"
	"strings"
)

// QueryOptions supports pagination and sorting for listing operations.
type QueryOptions struct {
	Page    int
	PageSize int
	SortBy  string
	SortDir string
}

type DriverFilters struct {
	Name          string
	LicenseNumber string
	Phone         string
	Email         string
	Address       string
}

type LoaderFilters struct {
	Name    string
	Phone   string
	Email   string
	Address string
	Company string
}

type VehicleFilters struct {
	VecID             string
	LicenseID         string
	VIN               string
	LicenseExpiration string
	DriverID          int
	OverweightPermitID string
}

type WellFilters struct {
	Name     string
	Location string
}

type LoadFilters struct {
	Date         string
	LoadNumber   string
	LoaderID     int
	WellID       int
	TicketNumber string
	DriverID     int
	VecID        string
	Status       *bool
}

func normalizeQueryOptions(opts QueryOptions) QueryOptions {
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 {
		opts.PageSize = 25
	}
	if opts.PageSize > 100 {
		opts.PageSize = 100
	}
	if strings.TrimSpace(opts.SortDir) == "" {
		opts.SortDir = "ASC"
	}
	if strings.EqualFold(opts.SortDir, "desc") {
		opts.SortDir = "DESC"
	} else {
		opts.SortDir = "ASC"
	}
	if strings.TrimSpace(opts.SortBy) == "" {
		opts.SortBy = "id"
	}
	return opts
}

func pageQuery(limit, offset int) string {
	if limit <= 0 {
		limit = 25
	}
	if offset < 0 {
		offset = 0
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

func whereClauseAndArgs(base string, conditions []string, args []any) (string, []any) {
	if len(conditions) == 0 {
		return base, args
	}
	return base + " WHERE " + strings.Join(conditions, " AND "), args
}

func addLikeCondition(clauses *[]string, args *[]any, field, value string) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return
	}
	*clauses = append(*clauses, fmt.Sprintf("LOWER(%s) LIKE ?", field))
	*args = append(*args, "%"+strings.ToLower(trimmed)+"%")
}

func addExactCondition(clauses *[]string, args *[]any, field string, value any) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return
		}
		*clauses = append(*clauses, fmt.Sprintf("%s = ?", field))
		*args = append(*args, v)
	case int:
		if v <= 0 {
			return
		}
		*clauses = append(*clauses, fmt.Sprintf("%s = ?", field))
		*args = append(*args, v)
	case bool:
		*clauses = append(*clauses, fmt.Sprintf("%s = ?", field))
		*args = append(*args, v)
	}
}

func validSortColumn(column string, allowed map[string]struct{}) string {
	if _, ok := allowed[column]; ok {
		return column
	}
	return "id"
}

func ListDriversQuery(db *sql.DB, filters DriverFilters, opts QueryOptions) ([]Driver, int64, error) {
	if db == nil {
		return nil, 0, fmt.Errorf("%w: database is nil", ErrValidation)
	}
	opts = normalizeQueryOptions(opts)
	allowed := map[string]struct{}{ "id":{}, "name":{}, "license_number":{}, "phone":{}, "email":{}, "address":{}, "created_at":{}, "updated_at":{} }
	ord := validSortColumn(opts.SortBy, allowed)
	if strings.EqualFold(opts.SortDir, "DESC") {
		ord += " DESC"
	} else {
		ord += " ASC"
	}

	clauses := make([]string, 0, 6)
	args := make([]any, 0, 6)
	addLikeCondition(&clauses, &args, "name", filters.Name)
	addLikeCondition(&clauses, &args, "license_number", filters.LicenseNumber)
	addLikeCondition(&clauses, &args, "phone", filters.Phone)
	addLikeCondition(&clauses, &args, "email", filters.Email)
	addLikeCondition(&clauses, &args, "address", filters.Address)
	base := "FROM driver"
	base, args = whereClauseAndArgs(base, clauses, args)

	countQuery := "SELECT COUNT(*) " + base
	var total int64
	if err := db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT id,name,license_number,phone,email,address,meta,created_at,updated_at %s ORDER BY %s %s", base, ord, pageQuery(opts.PageSize, (opts.Page-1)*opts.PageSize))
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]Driver, 0, 16)
	for rows.Next() {
		v, err := scanDriver(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *v)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (d *Database) ListDrivers(filters DriverFilters, opts QueryOptions) ([]Driver, int64, error) {
	return ListDriversQuery(d.db, filters, opts)
}

func ListLoadersQuery(db *sql.DB, filters LoaderFilters, opts QueryOptions) ([]Loader, int64, error) {
	if db == nil {
		return nil, 0, fmt.Errorf("%w: database is nil", ErrValidation)
	}
	opts = normalizeQueryOptions(opts)
	allowed := map[string]struct{}{ "id":{}, "name":{}, "phone":{}, "email":{}, "address":{}, "company":{}, "created_at":{}, "updated_at":{} }
	ord := validSortColumn(opts.SortBy, allowed)
	if strings.EqualFold(opts.SortDir, "DESC") {
		ord += " DESC"
	} else {
		ord += " ASC"
	}

	clauses := make([]string, 0, 6)
	args := make([]any, 0, 6)
	addLikeCondition(&clauses, &args, "name", filters.Name)
	addLikeCondition(&clauses, &args, "phone", filters.Phone)
	addLikeCondition(&clauses, &args, "email", filters.Email)
	addLikeCondition(&clauses, &args, "address", filters.Address)
	addLikeCondition(&clauses, &args, "company", filters.Company)
	base := "FROM loader"
	base, args = whereClauseAndArgs(base, clauses, args)

	countQuery := "SELECT COUNT(*) " + base
	var total int64
	if err := db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT id,name,phone,email,address,company,created_at,updated_at %s ORDER BY %s %s", base, ord, pageQuery(opts.PageSize, (opts.Page-1)*opts.PageSize))
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]Loader, 0, 16)
	for rows.Next() {
		v, err := scanLoader(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *v)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (d *Database) ListLoaders(filters LoaderFilters, opts QueryOptions) ([]Loader, int64, error) {
	return ListLoadersQuery(d.db, filters, opts)
}

func ListVehiclesQuery(db *sql.DB, filters VehicleFilters, opts QueryOptions) ([]Vehicle, int64, error) {
	if db == nil {
		return nil, 0, fmt.Errorf("%w: database is nil", ErrValidation)
	}
	opts = normalizeQueryOptions(opts)
	allowed := map[string]struct{}{ "id":{}, "vec_id":{}, "license_id":{}, "trailer_number":{}, "vin":{}, "license_expiration":{}, "driver_id":{}, "overweight_permit_id":{}, "created_at":{}, "updated_at":{} }
	ord := validSortColumn(opts.SortBy, allowed)
	if strings.EqualFold(opts.SortDir, "DESC") {
		ord += " DESC"
	} else {
		ord += " ASC"
	}

	clauses := make([]string, 0, 8)
	args := make([]any, 0, 8)
	addLikeCondition(&clauses, &args, "vec_id", filters.VecID)
	addLikeCondition(&clauses, &args, "license_id", filters.LicenseID)
	addLikeCondition(&clauses, &args, "vin", filters.VIN)
	addLikeCondition(&clauses, &args, "license_expiration", filters.LicenseExpiration)
	addLikeCondition(&clauses, &args, "overweight_permit_id", filters.OverweightPermitID)
	addExactCondition(&clauses, &args, "driver_id", filters.DriverID)
	base := "FROM vehicle"
	base, args = whereClauseAndArgs(base, clauses, args)

	countQuery := "SELECT COUNT(*) " + base
	var total int64
	if err := db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT id,vec_id,license_id,trailer_number,vin,license_expiration,driver_id,overweight_permit_id,created_at,updated_at %s ORDER BY %s %s", base, ord, pageQuery(opts.PageSize, (opts.Page-1)*opts.PageSize))
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]Vehicle, 0, 16)
	for rows.Next() {
		v, err := scanVehicle(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *v)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (d *Database) ListVehicles(filters VehicleFilters, opts QueryOptions) ([]Vehicle, int64, error) {
	return ListVehiclesQuery(d.db, filters, opts)
}

func ListWellsQuery(db *sql.DB, filters WellFilters, opts QueryOptions) ([]Well, int64, error) {
	if db == nil {
		return nil, 0, fmt.Errorf("%w: database is nil", ErrValidation)
	}
	opts = normalizeQueryOptions(opts)
	allowed := map[string]struct{}{ "id":{}, "name":{}, "location":{}, "created_at":{}, "updated_at":{} }
	ord := validSortColumn(opts.SortBy, allowed)
	if strings.EqualFold(opts.SortDir, "DESC") {
		ord += " DESC"
	} else {
		ord += " ASC"
	}

	clauses := make([]string, 0, 4)
	args := make([]any, 0, 4)
	addLikeCondition(&clauses, &args, "name", filters.Name)
	addLikeCondition(&clauses, &args, "location", filters.Location)
	base := "FROM well"
	base, args = whereClauseAndArgs(base, clauses, args)

	countQuery := "SELECT COUNT(*) " + base
	var total int64
	if err := db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT id,name,location,created_at,updated_at %s ORDER BY %s %s", base, ord, pageQuery(opts.PageSize, (opts.Page-1)*opts.PageSize))
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]Well, 0, 16)
	for rows.Next() {
		v, err := scanWell(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *v)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (d *Database) ListWells(filters WellFilters, opts QueryOptions) ([]Well, int64, error) {
	return ListWellsQuery(d.db, filters, opts)
}

func ListLoadsQuery(db *sql.DB, filters LoadFilters, opts QueryOptions) ([]Load, int64, error) {
	if db == nil {
		return nil, 0, fmt.Errorf("%w: database is nil", ErrValidation)
	}
	opts = normalizeQueryOptions(opts)
	allowed := map[string]struct{}{ "id":{}, "date":{}, "load_number":{}, "loader_id":{}, "well_id":{}, "ticket_number":{}, "miles":{}, "driver_id":{}, "vec_id":{}, "trailer_number":{}, "net_weight":{}, "tons":{}, "ton_rate":{}, "total_value":{}, "status":{}, "ticket_id":{}, "created_at":{}, "updated_at":{} }
	ord := validSortColumn(opts.SortBy, allowed)
	if strings.EqualFold(opts.SortDir, "DESC") {
		ord += " DESC"
	} else {
		ord += " ASC"
	}

	clauses := make([]string, 0, 10)
	args := make([]any, 0, 10)
	addLikeCondition(&clauses, &args, "date", filters.Date)
	addLikeCondition(&clauses, &args, "load_number", filters.LoadNumber)
	addLikeCondition(&clauses, &args, "ticket_number", filters.TicketNumber)
	addLikeCondition(&clauses, &args, "vec_id", filters.VecID)
	addExactCondition(&clauses, &args, "loader_id", filters.LoaderID)
	addExactCondition(&clauses, &args, "well_id", filters.WellID)
	addExactCondition(&clauses, &args, "driver_id", filters.DriverID)
	if filters.Status != nil {
		clauses = append(clauses, "status = ?")
		args = append(args, *filters.Status)
	}
	base := "FROM load"
	base, args = whereClauseAndArgs(base, clauses, args)

	countQuery := "SELECT COUNT(*) " + base
	var total int64
	if err := db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT id,date,load_number,loader_id,well_id,ticket_number,miles,driver_id,vec_id,trailer_number,net_weight,tons,ton_rate,total_value,status,ticket_id,created_at,updated_at %s ORDER BY %s %s", base, ord, pageQuery(opts.PageSize, (opts.Page-1)*opts.PageSize))
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]Load, 0, 16)
	for rows.Next() {
		v, err := scanLoad(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *v)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (d *Database) ListLoads(filters LoadFilters, opts QueryOptions) ([]Load, int64, error) {
	return ListLoadsQuery(d.db, filters, opts)
}

// Convenience wrappers keep older list methods working while adding richer paging and sorting support.
func (d *Database) GetDriversPage(filters DriverFilters, opts QueryOptions) ([]Driver, int64, error) {
	return d.ListDrivers(filters, opts)
}

func (d *Database) GetLoadersPage(filters LoaderFilters, opts QueryOptions) ([]Loader, int64, error) {
	return d.ListLoaders(filters, opts)
}

func (d *Database) GetVehiclesPage(filters VehicleFilters, opts QueryOptions) ([]Vehicle, int64, error) {
	return d.ListVehicles(filters, opts)
}

func (d *Database) GetWellsPage(filters WellFilters, opts QueryOptions) ([]Well, int64, error) {
	return d.ListWells(filters, opts)
}

func (d *Database) GetLoadsPage(filters LoadFilters, opts QueryOptions) ([]Load, int64, error) {
	return d.ListLoads(filters, opts)
}
