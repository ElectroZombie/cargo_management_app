package database

import (
    "database/sql"
    "fmt"
    "strings"
)

// Txn executes a callback inside a database transaction and automatically rolls back on error.
func Txn(db *sql.DB, fn func(*sql.Tx) error) error {
    if db == nil {
        return fmt.Errorf("%w: database is nil", ErrValidation)
    }

    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }

    if err := fn(tx); err != nil {
        _ = tx.Rollback()
        return err
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit transaction: %w", err)
    }

    return nil
}

// ExecInTransaction wraps a callback and ensures rollback-on-error semantics for multi-step database work.
func ExecInTransaction(db *sql.DB, fn func(*sql.Tx) error) error {
    return Txn(db, fn)
}

// WithTransaction is a convenience alias used by service code that needs an explicit transaction boundary.
func WithTransaction(db *sql.DB, fn func(*sql.Tx) error) error {
    return Txn(db, fn)
}

// Transaction helpers for multi-entity operations.
func (d *Database) CreateDriverWithLoader(db *sql.DB, driver DriverDTO, loader LoaderDTO) (*Driver, *Loader, error) {
    if db == nil {
        return nil, nil, fmt.Errorf("%w: database is nil", ErrValidation)
    }

    var createdDriver *Driver
    var createdLoader *Loader

    err := Txn(db, func(tx *sql.Tx) error {
        driverSvc := &Database{db: tx}

        if err := validateDriver(driver); err != nil {
            return err
        }
        if err := validateLoader(loader); err != nil {
            return err
        }

        r, err := tx.Exec(`INSERT INTO driver(name,license_number,phone,email,address) VALUES(?,?,?,?,?)`, driver.Name, driver.LicenseNumber, driver.Phone, driver.Email, driver.Address)
        if err != nil {
            return databaseError("create driver", err)
        }
        driverID, err := r.LastInsertId()
        if err != nil {
            return databaseError("create driver", err)
        }

        createdDriver, err = driverSvc.GetDriver(int(driverID))
        if err != nil {
            return err
        }

        r, err = tx.Exec(`INSERT INTO loader(name,phone,email,address,company) VALUES(?,?,?,?,?)`, loader.Name, loader.Phone, loader.Email, loader.Address, loader.Company)
        if err != nil {
            return databaseError("create loader", err)
        }
        loaderID, err := r.LastInsertId()
        if err != nil {
            return databaseError("create loader", err)
        }

        createdLoader, err = driverSvc.GetLoader(int(loaderID))
        if err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        return nil, nil, err
    }

    return createdDriver, createdLoader, nil
}

func (d *Database) CreateLoadWithVehicle(db *sql.DB, load LoadDTO, vehicle VehicleDTO) (*Load, *Vehicle, error) {
    if db == nil {
        return nil, nil, fmt.Errorf("%w: database is nil", ErrValidation)
    }

    var createdLoad *Load
    var createdVehicle *Vehicle

    err := Txn(db, func(tx *sql.Tx) error {
        txnDB := &Database{db: tx}

        if err := validateLoad(load); err != nil {
            return err
        }
        if err := validateVehicle(vehicle); err != nil {
            return err
        }
        if err := txnDB.validateLoadReferences(load); err != nil {
            return err
        }
        if err := txnDB.validateVehicleReferences(vehicle); err != nil {
            return err
        }

        r, err := tx.Exec(`INSERT INTO vehicle(vec_id,license_id,trailer_number,vin,license_expiration,driver_id,overweight_permit_id) VALUES(?,?,?,?,?,?,?)`, vehicle.VecID, vehicle.LicenseID, vehicle.TrailerNumber, vehicle.VIN, vehicle.LicenseExpiration, vehicle.DriverID, vehicle.OverweightPermitID)
        if err != nil {
            return databaseError("create vehicle", err)
        }
        vehicleID, err := r.LastInsertId()
        if err != nil {
            return databaseError("create vehicle", err)
        }
        createdVehicle, err = txnDB.GetVehicle(int(vehicleID))
        if err != nil {
            return err
        }

        r, err = tx.Exec(`INSERT INTO load(date,load_number,loader_id,well_id,ticket_number,miles,driver_id,vec_id,trailer_number,net_weight,tons,ton_rate,total_value,status,ticket_id) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, load.Date, load.LoadNumber, load.LoaderID, load.WellID, load.TicketNumber, load.Miles, load.DriverID, load.VecID, load.TrailerNumber, load.NetWeight, load.Tons, load.TonRate, load.TotalValue, load.Status, load.TicketID)
        if err != nil {
            return databaseError("create load", err)
        }
        loadID, err := r.LastInsertId()
        if err != nil {
            return databaseError("create load", err)
        }

        createdLoad, err = txnDB.GetLoad(int(loadID))
        if err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        return nil, nil, err
    }

    return createdLoad, createdVehicle, nil
}

func (d *Database) DeleteDriverAndLoads(db *sql.DB, driverID int) error {
    if db == nil {
        return fmt.Errorf("%w: database is nil", ErrValidation)
    }

    return Txn(db, func(tx *sql.Tx) error {
        txnDB := &Database{db: tx}

        if _, err := tx.Exec(`DELETE FROM load WHERE driver_id = ?`, driverID); err != nil {
            return databaseError("delete loads for driver", err)
        }
        if _, err := tx.Exec(`DELETE FROM driver WHERE id = ?`, driverID); err != nil {
            return databaseError("delete driver", err)
        }

        if _, err := txnDB.GetDriver(driverID); err == nil {
            return fmt.Errorf("%w: driver still exists after delete", ErrConflict)
        } else if !strings.Contains(err.Error(), ErrNotFound.Error()) {
            return err
        }

        return nil
    })
}
