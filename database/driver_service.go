package database

import (
    "database/sql"
    "fmt"
)

const driverColumns = `id,name,license_number,phone,email,address,created_at,updated_at`
func (d *Database) CreateDriver(v DriverDTO) (*Driver,error) { if err:=required(v.Name,"name");err!=nil{return nil,err}; if err:=required(v.LicenseNumber,"license_number");err!=nil{return nil,err}; r,e:=d.db.Exec(`INSERT INTO driver(name,license_number,phone,email,address) VALUES(?,?,?,?,?)`,v.Name,v.LicenseNumber,v.Phone,v.Email,v.Address);if e!=nil{return nil,fmt.Errorf("create driver: %w",e)};id,_:=r.LastInsertId();return d.GetDriver(int(id)) }
func (d *Database) GetDrivers()([]Driver,error){rows,e:=d.db.Query(`SELECT `+driverColumns+` FROM driver ORDER BY id`);if e!=nil{return nil,e};defer rows.Close();var out []Driver;for rows.Next(){v,e:=scanDriver(rows);if e!=nil{return nil,e};out=append(out,*v)};return out,rows.Err()}
func (d *Database) GetDriver(id int)(*Driver,error){v,e:=scanDriver(d.db.QueryRow(`SELECT `+driverColumns+` FROM driver WHERE id=?`,id));if e==sql.ErrNoRows{return nil,fmt.Errorf("driver not found")};return v,e}
func (d *Database) UpdateDriver(id int,v DriverDTO)error{if err:=required(v.Name,"name");err!=nil{return err};r,e:=d.db.Exec(`UPDATE driver SET name=?,license_number=?,phone=?,email=?,address=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`,v.Name,v.LicenseNumber,v.Phone,v.Email,v.Address,id);if e!=nil{return e};return affected(r,"driver")}
func (d *Database) DeleteDriver(id int)error{r,e:=d.db.Exec(`DELETE FROM driver WHERE id=?`,id);if e!=nil{return e};return affected(r,"driver")}
