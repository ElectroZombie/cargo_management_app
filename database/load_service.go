package database

import (
	"database/sql"
	"fmt"
)

const loadColumns = `id,date,load_number,loader_id,well_id,ticket_number,miles,driver_id,vec_id,trailer_number,net_weight,tons,ton_rate,total_value,status,ticket_id,created_at,updated_at`

func (d *Database) CreateLoad(v LoadDTO) (*Load, error) {
	if err := validateLoad(v); err != nil {
		return nil, err
	}
	if err := d.validateLoadReferences(v); err != nil {
		return nil, err
	}
	r, err := d.db.Exec(`INSERT INTO load(date,load_number,loader_id,well_id,ticket_number,miles,driver_id,vec_id,trailer_number,net_weight,tons,ton_rate,total_value,status,ticket_id) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, v.Date, v.LoadNumber, v.LoaderID, v.WellID, v.TicketNumber, v.Miles, v.DriverID, v.VecID, v.TrailerNumber, v.NetWeight, v.Tons, v.TonRate, v.TotalValue, v.Status, v.TicketID)
	if err != nil {
		return nil, databaseError("create load", err)
	}
	id, err := r.LastInsertId()
	if err != nil {
		return nil, databaseError("create load", err)
	}
	return d.GetLoad(int(id))
}

func (d *Database) GetLoads() ([]Load, error) {
	rows, err := d.db.Query(`SELECT `+loadColumns+` FROM load ORDER BY date DESC,id DESC`)
	if err != nil {
		return nil, databaseError("read loads", err)
	}
	defer rows.Close()
	var out []Load
	for rows.Next() {
		v, err := scanLoad(rows)
		if err != nil {
			return nil, databaseError("read load", err)
		}
		out = append(out, *v)
	}
	return out, databaseError("read loads", rows.Err())
}

func (d *Database) GetLoad(id int) (*Load, error) {
	v, err := scanLoad(d.db.QueryRow(`SELECT `+loadColumns+` FROM load WHERE id=?`, id))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: load", ErrNotFound)
	}
	if err != nil {
		return nil, databaseError("read load", err)
	}
	return v, nil
}

func (d *Database) SearchLoads(query string) ([]Load, error) {
	items, err := d.GetLoads()
	if err != nil {
		return nil, err
	}
	var out []Load
	for _, v := range items {
		if containsSimilar(query, v.Date, v.LoadNumber, v.TicketNumber, v.VecID, v.TicketID) {
			out = append(out, v)
		}
	}
	return out, nil
}

func (d *Database) UpdateLoad(id int, v LoadDTO) error {
	if err := validateLoad(v); err != nil { return err }
	if err := d.validateLoadReferences(v); err != nil { return err }
	r, err := d.db.Exec(`UPDATE load SET date=?,load_number=?,loader_id=?,well_id=?,ticket_number=?,miles=?,driver_id=?,vec_id=?,trailer_number=?,net_weight=?,tons=?,ton_rate=?,total_value=?,status=?,ticket_id=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`, v.Date, v.LoadNumber, v.LoaderID, v.WellID, v.TicketNumber, v.Miles, v.DriverID, v.VecID, v.TrailerNumber, v.NetWeight, v.Tons, v.TonRate, v.TotalValue, v.Status, v.TicketID, id)
	if err != nil {
		return databaseError("update load", err)
	}
	return affected(r, "load")
}

func (d *Database) DeleteLoad(id int) error {
	r, err := d.db.Exec(`DELETE FROM load WHERE id=?`, id)
	if err != nil {
		return databaseError("delete load", err)
	}
	return affected(r, "load")
}

func (d *Database) GetLoadsByDriver(driverID int) ([]Load, error) {
	if err := positiveID(driverID, "driver_id"); err != nil { return nil, err }
	return d.loadsWhere(`driver_id=?`, driverID)
}

func (d *Database) GetLoadsByWell(wellID int) ([]Load, error) {
	if err := positiveID(wellID, "well_id"); err != nil { return nil, err }
	return d.loadsWhere(`well_id=?`, wellID)
}

func (d *Database) GetLoadsByStatus(status bool) ([]Load, error) {
	return d.loadsWhere(`status=?`, status)
}

func (d *Database) loadsWhere(where string, arg any) ([]Load, error) {
	rows, err := d.db.Query(`SELECT `+loadColumns+` FROM load WHERE `+where+` ORDER BY date DESC,id DESC`, arg)
	if err != nil {
		return nil, databaseError("read loads", err)
	}
	defer rows.Close()
	var out []Load
	for rows.Next() {
		v, err := scanLoad(rows)
		if err != nil {
			return nil, databaseError("read load", err)
		}
		out = append(out, *v)
	}
	return out, databaseError("read loads", rows.Err())
}
