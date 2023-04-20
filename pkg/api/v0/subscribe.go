package v0

type SubscribeRequest struct {
	Email    string `json:"email" validate:"required,email"`
	EntityId string `json:"entityId" validate:"required"`
}

type SubscribeResponse struct {
	Error string `json:"error,omitempty"`
}
