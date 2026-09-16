package main

import "time"

// ============ DRIVER MODEL ============
type Driver struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	License   string    `json:"license_number"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relationships
	Vehicles []Vehicle `json:"vehicles,omitempty"` // One driver can have many vehicles
	Loads    []Load    `json:"loads,omitempty"`    // One driver can have many loads
}

// ============ LOADER MODEL ============
type Loader struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	Company   string    `json:"company"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relationships
	Loads []Load `json:"loads,omitempty"` // One loader can have many loads
}

// ============ VEHICLE MODEL ============
type Vehicle struct {
	ID                   int       `json:"id"`
	VecID                string    `json:"vec_id"`
	LicenseID            string    `json:"license_id"`
	TrailerNumber        int       `json:"trailer_number"`
	VIN                  string    `json:"vin"`
	LicenseExpiration    time.Time `json:"license_expiration"`
	DriverID             int       `json:"driver_id"` // Foreign Key
	OverweightPermitID   string    `json:"overweight_permit_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	
	// Relationships
	Driver Driver `json:"driver,omitempty"` // Many vehicles belong to one driver
	Loads  []Load `json:"loads,omitempty"` // One vehicle can have many loads
}

// ============ WELL MODEL ============
type Well struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Client    string    `json:"client"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relationships
	Loads []Load `json:"loads,omitempty"` // One well can have many loads
}

// ============ LOAD MODEL ============
type Load struct {
	ID            int       `json:"id"`
	Date          time.Time `json:"date"`
	LoadNumber    string    `json:"load_number"`
	LoaderID      int       `json:"loader_id"`      // Foreign Key
	WellID        int       `json:"well_id"`        // Foreign Key
	TicketNumber  string    `json:"ticket_number"`
	Miles         float64   `json:"miles"`
	DriverID      int       `json:"driver_id"`      // Foreign Key
	VecID         string    `json:"vec_id"`         // Foreign Key
	TrailerNumber int       `json:"trailer_number"`
	NetWeight     float64   `json:"net_weight"`
	Tons          float64   `json:"tons"`
	TonRate       float64   `json:"ton_rate"`
	TotalValue    float64   `json:"total_value"`
	Status        bool      `json:"status"` // true = active, false = inactive
	TicketID      string    `json:"ticket_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	
	// Relationships
	Loader  Loader  `json:"loader,omitempty"`  // Many loads belong to one loader
	Well    Well    `json:"well,omitempty"`    // Many loads belong to one well
	Driver  Driver  `json:"driver,omitempty"`  // Many loads belong to one driver
	Vehicle Vehicle `json:"vehicle,omitempty"` // Many loads belong to one vehicle
}

// ============ RELATIONSHIP DIAGRAM ============
/*
Relationships Overview:

DRIVER (1) -------- (M) VEHICLE
  |
  +---------- (M) LOAD

LOADER (1) -------- (M) LOAD

VEHICLE (1) -------- (M) LOAD

WELL (1) -------- (M) LOAD

Where:
- 1 = One
- M = Many

DRIVER has a 1-to-Many relationship with VEHICLE (one driver can operate many vehicles)
DRIVER has a 1-to-Many relationship with LOAD (one driver can have many loads assigned)

LOADER has a 1-to-Many relationship with LOAD (one loader can load many loads)

VEHICLE has a 1-to-Many relationship with LOAD (one vehicle can carry many loads)

WELL has a 1-to-Many relationship with LOAD (one well can have many loads)

LOAD is a junction table that connects:
- DRIVER (who drives)
- LOADER (who loads)
- VEHICLE (what vehicle is used)
- WELL (where the load comes from)
*/
