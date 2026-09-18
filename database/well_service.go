package database

import("database/sql";"fmt")
const wellColumns=`id,name,client,location,created_at,updated_at`
func(d *Database)CreateWell(v WellDTO)(*Well,error){if e:=required(v.Name,"name");e!=nil{return nil,e};if e:=required(v.Client,"client");e!=nil{return nil,e};r,e:=d.db.Exec(`INSERT INTO well(name,client,location) VALUES(?,?,?)`,v.Name,v.Client,v.Location);if e!=nil{return nil,fmt.Errorf("create well: %w",e)};id,_:=r.LastInsertId();return d.GetWell(int(id))}
func(d *Database)GetWells()([]Well,error){rows,e:=d.db.Query(`SELECT `+wellColumns+` FROM well ORDER BY id`);if e!=nil{return nil,e};defer rows.Close();var out []Well;for rows.Next(){v,e:=scanWell(rows);if e!=nil{return nil,e};out=append(out,*v)};return out,rows.Err()}
func(d *Database)GetWell(id int)(*Well,error){v,e:=scanWell(d.db.QueryRow(`SELECT `+wellColumns+` FROM well WHERE id=?`,id));if e==sql.ErrNoRows{return nil,fmt.Errorf("well not found")};return v,e}
func(d *Database)UpdateWell(id int,v WellDTO)error{r,e:=d.db.Exec(`UPDATE well SET name=?,client=?,location=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`,v.Name,v.Client,v.Location,id);if e!=nil{return e};return affected(r,"well")}
func(d *Database)DeleteWell(id int)error{r,e:=d.db.Exec(`DELETE FROM well WHERE id=?`,id);if e!=nil{return e};return affected(r,"well")}
