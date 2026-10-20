package models

// Load represents the read model for a load entity returned by services.
type Load struct {
	ID            int     `json:"id"`
	Date          string  `json:"date"`
	LoadNumber    string  `json:"load_number"`
	LoaderID      int     `json:"loader_id"`
	WellID        int     `json:"well_id"`
	TicketNumber  string  `json:"ticket_number"`
	Miles         float64 `json:"miles"`
	DriverID      int     `json:"driver_id"`
	VecID         string  `json:"vec_id"`
	TrailerNumber int     `json:"trailer_number"`
	NetWeight     float64 `json:"net_weight"`
	Tons          float64 `json:"tons"`
	TonRate       float64 `json:"ton_rate"`
	TotalValue    float64 `json:"total_value"`
	Status        bool    `json:"status"`
	TicketID      string  `json:"ticket_id"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}