package api

import "time"

type GetTweetsRequest struct {
	EntityID   string    `json:"entityId"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	MaxResults int       `json:"maxResults" validate:"gte=5,lte=100"`
}

type GetTweetsResponse struct {
	Error  string    `json:"error,omitempty"`
	Result Analytics `json:"result,omitempty"`
}

type Analytics struct {
	Total       int     `json:"total"`
	Authentic   float32 `json:"authentic"`
	Unauthentic float32 `json:"unauthentic"`
}
