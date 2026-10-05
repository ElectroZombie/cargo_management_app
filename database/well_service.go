package database

import("database/sql";"fmt")
const wellColumns=`id,name,client,location,created_at,updated_at`
func(d *Database)CreateWell(v WellDTO)(*Well,error){if e:=validateWell(v);e!=nil{return nil,e};r,e:=d.db.Exec(`INSERT INTO well(name,client,location) VALUES(?,?,?)`,v.Name,v.Client,v.Location);if e!=nil{return nil,databaseError("create well",e)};id,e:=r.LastInsertId();if e!=nil{return nil,databaseError("create well",e)};return d.GetWell(int(id))}
func(d *Database)GetWells()([]Well,error){rows,e:=d.db.Query(`SELECT `+wellColumns+` FROM well ORDER BY id`);if e!=nil{return nil,databaseError("read wells",e)};defer rows.Close();var out []Well;for rows.Next(){v,e:=scanWell(rows);if e!=nil{return nil,databaseError("read well",e)};out=append(out,*v)};return out,databaseError("read wells",rows.Err())}
func(d *Database)GetWell(id int)(*Well,error){v,e:=scanWell(d.db.QueryRow(`SELECT `+wellColumns+` FROM well WHERE id=?`,id));if e==sql.ErrNoRows{return nil,fmt.Errorf("%w: well",ErrNotFound)};if e!=nil{return nil,databaseError("read well",e)};return v,nil}
func(d *Database)SearchWells(query string)([]Well,error){items,e:=d.GetWells();if e!=nil{return nil,e};var out []Well;for _,v:=range items{if containsSimilar(query,v.Name,v.Client,v.Location){out=append(out,v)}};return out,nil}
func(d *Database)UpdateWell(id int,v WellDTO)error{if e:=validateWell(v);e!=nil{return e};r,e:=d.db.Exec(`UPDATE well SET name=?,client=?,location=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`,v.Name,v.Client,v.Location,id);if e!=nil{return databaseError("update well",e)};return affected(r,"well")}
func(d *Database)DeleteWell(id int)error{r,e:=d.db.Exec(`DELETE FROM well WHERE id=?`,id);if e!=nil{return databaseError("delete well",e)};return affected(r,"well")}
