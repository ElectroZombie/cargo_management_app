package models

// Well represents the read model for a well entity returned by services.
type Well struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Client    string `json:"client"`
	Location  string `json:"location"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}