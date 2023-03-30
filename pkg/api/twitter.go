package api

import "time"

type GetTweetsRequest struct {
	EntityID   string    `json:"entityId"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	MaxResults int       `json:"maxResults" validate:"gte=5,lte=100"`
}

type GetTweetsResponse struct {
	Error  string  `json:"error,omitempty"`
	Tweets []Tweet `json:"tweets"`
}

type Tweet struct {
	ID            string    `json:"id"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"createdAt"`
	RealnessScore float32   `json:"realnessScore"`
}
