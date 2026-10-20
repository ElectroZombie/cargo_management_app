package models

// Vehicle represents the read model for a vehicle entity returned by services.
type Vehicle struct {
	ID                 int    `json:"id"`
	VecID              string `json:"vec_id"`
	LicenseID          string `json:"license_id"`
	TrailerNumber      int    `json:"trailer_number"`
	VIN                string `json:"vin"`
	LicenseExpiration  string `json:"license_expiration"`
	DriverID           int    `json:"driver_id"`
	OverweightPermitID string `json:"overweight_permit_id"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}