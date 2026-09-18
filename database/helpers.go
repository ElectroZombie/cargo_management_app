package database

import (
    "database/sql"
    "fmt"
)

func required(value, field string) error { if value == "" { return fmt.Errorf("%s is required", field) }; return nil }
func affected(result sql.Result, entity string) error { n, err := result.RowsAffected(); if err != nil { return err }; if n == 0 { return fmt.Errorf("%s not found", entity) }; return nil }
func scanDriver(row interface{ Scan(...any) error }) (*Driver, error) { var v Driver; err := row.Scan(&v.ID,&v.Name,&v.LicenseNumber,&v.Phone,&v.Email,&v.Address,&v.CreatedAt,&v.UpdatedAt); return &v, err }
func scanLoader(row interface{ Scan(...any) error }) (*Loader, error) { var v Loader; err := row.Scan(&v.ID,&v.Name,&v.Phone,&v.Email,&v.Address,&v.Company,&v.CreatedAt,&v.UpdatedAt); return &v, err }
func scanVehicle(row interface{ Scan(...any) error }) (*Vehicle, error) { var v Vehicle; err := row.Scan(&v.ID,&v.VecID,&v.LicenseID,&v.TrailerNumber,&v.VIN,&v.LicenseExpiration,&v.DriverID,&v.OverweightPermitID,&v.CreatedAt,&v.UpdatedAt); return &v, err }
func scanWell(row interface{ Scan(...any) error }) (*Well, error) { var v Well; err := row.Scan(&v.ID,&v.Name,&v.Client,&v.Location,&v.CreatedAt,&v.UpdatedAt); return &v, err }
func scanLoad(row interface{ Scan(...any) error }) (*Load, error) { var v Load; err := row.Scan(&v.ID,&v.Date,&v.LoadNumber,&v.LoaderID,&v.WellID,&v.TicketNumber,&v.Miles,&v.DriverID,&v.VecID,&v.TrailerNumber,&v.NetWeight,&v.Tons,&v.TonRate,&v.TotalValue,&v.Status,&v.TicketID,&v.CreatedAt,&v.UpdatedAt); return &v, err }
