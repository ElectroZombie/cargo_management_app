package database

import "database/sql"

const loaderColumns = `id,name,phone,email,address,company,created_at,updated_at`

func (d *Database) CreateLoader(v LoaderDTO) (*Loader, error) {
	if e := validateLoader(v); e != nil {
		return nil, e
	}
	r, e := d.runner.Exec(`INSERT INTO loader(name,phone,email,address,company) VALUES(?,?,?,?,?)`, v.Name, v.Phone, v.Email, v.Address, v.Company)
	if e != nil {
		return nil, databaseError("create loader", e)
	}
	id, e := r.LastInsertId()
	if e != nil {
		return nil, databaseError("create loader", e)
	}
	return d.GetLoader(int(id))
}
func (d *Database) GetLoaders() ([]Loader, error) {
	rows, e := d.runner.Query(`SELECT ` + loaderColumns + ` FROM loader ORDER BY id`)
	if e != nil {
		return nil, databaseError("read loaders", e)
	}
	defer rows.Close()
	var out []Loader
	for rows.Next() {
		v, e := scanLoader(rows)
		if e != nil {
			return nil, databaseError("read loader", e)
		}
		out = append(out, *v)
	}
	return out, databaseError("read loaders", rows.Err())
}
func (d *Database) GetLoader(id int) (*Loader, error) {
	v, e := scanLoader(d.runner.QueryRow(`SELECT `+loaderColumns+` FROM loader WHERE id=?`, id))
	if e == sql.ErrNoRows {
		return nil, notFound("loader")
	}
	if e != nil {
		return nil, databaseError("read loader", e)
	}
	return v, nil
}
func (d *Database) SearchLoaders(query string) ([]Loader, error) {
	items, e := d.GetLoaders()
	if e != nil {
		return nil, e
	}
	var out []Loader
	for _, v := range items {
		if containsSimilar(query, v.Name, v.Phone, v.Email, v.Address, v.Company) {
			out = append(out, v)
		}
	}
	return out, nil
}
func (d *Database) UpdateLoader(id int, v LoaderDTO) error {
	if e := validateLoader(v); e != nil {
		return e
	}
	r, e := d.runner.Exec(`UPDATE loader SET name=?,phone=?,email=?,address=?,company=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`, v.Name, v.Phone, v.Email, v.Address, v.Company, id)
	if e != nil {
		return databaseError("update loader", e)
	}
	return affected(r, "loader")
}
func (d *Database) DeleteLoader(id int) error {
	r, e := d.runner.Exec(`DELETE FROM loader WHERE id=?`, id)
	if e != nil {
		return deleteError("delete loader", "loader", e)
	}
	return affected(r, "loader")
}
