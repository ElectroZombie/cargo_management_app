package dtos

// LoaderDTO is the write DTO for creating/updating a loader.
type LoaderDTO struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Company string `json:"company"`
}