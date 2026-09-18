package database

import("database/sql";"fmt")
const loaderColumns=`id,name,phone,email,address,company,created_at,updated_at`
func(d *Database)CreateLoader(v LoaderDTO)(*Loader,error){if e:=required(v.Name,"name");e!=nil{return nil,e};r,e:=d.db.Exec(`INSERT INTO loader(name,phone,email,address,company) VALUES(?,?,?,?,?)`,v.Name,v.Phone,v.Email,v.Address,v.Company);if e!=nil{return nil,fmt.Errorf("create loader: %w",e)};id,_:=r.LastInsertId();return d.GetLoader(int(id))}
func(d *Database)GetLoaders()([]Loader,error){rows,e:=d.db.Query(`SELECT `+loaderColumns+` FROM loader ORDER BY id`);if e!=nil{return nil,e};defer rows.Close();var out []Loader;for rows.Next(){v,e:=scanLoader(rows);if e!=nil{return nil,e};out=append(out,*v)};return out,rows.Err()}
func(d *Database)GetLoader(id int)(*Loader,error){v,e:=scanLoader(d.db.QueryRow(`SELECT `+loaderColumns+` FROM loader WHERE id=?`,id));if e==sql.ErrNoRows{return nil,fmt.Errorf("loader not found")};return v,e}
func(d *Database)UpdateLoader(id int,v LoaderDTO)error{r,e:=d.db.Exec(`UPDATE loader SET name=?,phone=?,email=?,address=?,company=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`,v.Name,v.Phone,v.Email,v.Address,v.Company,id);if e!=nil{return e};return affected(r,"loader")}
func(d *Database)DeleteLoader(id int)error{r,e:=d.db.Exec(`DELETE FROM loader WHERE id=?`,id);if e!=nil{return e};return affected(r,"loader")}
