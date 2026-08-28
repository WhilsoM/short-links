package dto

import "time"

type Link struct {
	ID        int       `json:"id"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateLinkRequest struct {
	Link string `json:"link"`
}

type CreateLinkResponse struct {
	ID        int       `json:"id"`
	ShortLink string    `json:"short_link"`
	CreatedAt time.Time `json:"created_at"`
}
