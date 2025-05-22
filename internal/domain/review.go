package domain

import "time"

type Review struct {
	ID         string    `json:"id"`
	AuthorID   string    `json:"author_id"`
	AuthorName string    `json:"author_name"`
	TargetID   string    `json:"target_id"`
	Rating     float64   `json:"rating"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateReviewRequest struct {
	TargetID string `json:"target_id" validate:"required"`
	Rating   int    `json:"rating" validate:"required,min=1,max=5"`
	Comment  string `json:"comment" validate:"required"`
}
