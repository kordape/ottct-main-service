package api

type UpdateSubscriptionRequest struct {
	Subscribe bool `json:"subscribe" validate:"required"`
}
