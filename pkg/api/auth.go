package api

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone" validate:"required,numeric"`
}

type SignUpResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
