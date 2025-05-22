package domain

import "time"

const (
	Give   string = "give"
	Search string = "search"
)

type CreateListingRequest struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description"`
	Latitude    float64  `json:"latitude" validate:"required"`
	Longitude   float64  `json:"longitude" validate:"required"`
	PreviewURL  string   `json:"preview_url"`
	Status      string   `json:"status" validate:"required,oneof=search give"`
	Images      []string `json:"images"`
	City        string   `json:"city" validate:"required"`
	Street      string   `json:"street"`
}

type Owner struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Surname string  `json:"surname"`
	Phone   string  `json:"phone"`
	Rating  float64 `json:"rating"`
}

type Listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	DistanceM   float64   `json:"distance_m"`
	City        string    `json:"city"`
	Street      string    `json:"street"`
	PreviewURL  string    `json:"preview_url"`
	Images      []string  `json:"images"`
	Status      string    `json:"status"`
	Owner       Owner     `json:"owner"`
	CreatedAt   time.Time `json:"created_at"`
}
