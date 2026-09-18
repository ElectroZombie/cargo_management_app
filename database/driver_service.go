package database

import("database/sql";"fmt")
const driverColumns=`id,name,license_number,phone,email,address,created_at,updated_at`
func(d *Database)CreateDriver(v DriverDTO)(*Driver,error){if e:=validateDriver(v);e!=nil{return nil,e};r,e:=d.db.Exec(`INSERT INTO driver(name,license_number,phone,email,address) VALUES(?,?,?,?,?)`,v.Name,v.LicenseNumber,v.Phone,v.Email,v.Address);if e!=nil{return nil,databaseError("create driver",e)};id,e:=r.LastInsertId();if e!=nil{return nil,databaseError("create driver",e)};return d.GetDriver(int(id))}
func(d *Database)GetDrivers()([]Driver,error){rows,e:=d.db.Query(`SELECT `+driverColumns+` FROM driver ORDER BY id`);if e!=nil{return nil,databaseError("read drivers",e)};defer rows.Close();var out []Driver;for rows.Next(){v,e:=scanDriver(rows);if e!=nil{return nil,databaseError("read driver",e)};out=append(out,*v)};return out,databaseError("read drivers",rows.Err())}
func(d *Database)GetDriver(id int)(*Driver,error){v,e:=scanDriver(d.db.QueryRow(`SELECT `+driverColumns+` FROM driver WHERE id=?`,id));if e==sql.ErrNoRows{return nil,fmt.Errorf("%w: driver",ErrNotFound)};if e!=nil{return nil,databaseError("read driver",e)};return v,nil}
func(d *Database)SearchDrivers(query string)([]Driver,error){items,e:=d.GetDrivers();if e!=nil{return nil,e};var out []Driver;for _,v:=range items{if containsSimilar(query,v.Name,v.LicenseNumber,v.Phone,v.Email,v.Address){out=append(out,v)}};return out,nil}
func(d *Database)UpdateDriver(id int,v DriverDTO)error{if e:=validateDriver(v);e!=nil{return e};r,e:=d.db.Exec(`UPDATE driver SET name=?,license_number=?,phone=?,email=?,address=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`,v.Name,v.LicenseNumber,v.Phone,v.Email,v.Address,id);if e!=nil{return databaseError("update driver",e)};return affected(r,"driver")}
func(d *Database)DeleteDriver(id int)error{r,e:=d.db.Exec(`DELETE FROM driver WHERE id=?`,id);if e!=nil{return databaseError("delete driver",e)};return affected(r,"driver")}
