package domain

type Draft struct {
	ID          string   `json:"id" bson:"_id"`
	Title       string   `json:"title" bson:"title"`
	OwnerID     string   `json:"owner_id" bson:"owner_id"`
	Description string   `json:"description" bson:"description"`
	Latitude    float64  `json:"latitude" bson:"latitude"`
	Longitude   float64  `json:"longitude" bson:"longitude"`
	PreviewURL  string   `json:"preview_url" bson:"preview_url"`
	Status      string   `json:"status" bson:"status"`
	Images      []string `json:"images" bson:"images"`
	City        string   `json:"city" bson:"city"`
	Street      string   `json:"street" bson:"street"`
}
