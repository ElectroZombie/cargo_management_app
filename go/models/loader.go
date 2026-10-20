package models

// Loader represents the read model for a loader entity returned by services.
type Loader struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Company   string `json:"company"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}