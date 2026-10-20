package dtos

// WellDTO is the write DTO for creating/updating a well.
type WellDTO struct {
	Name     string `json:"name"`
	Client   string `json:"client"`
	Location string `json:"location"`
}