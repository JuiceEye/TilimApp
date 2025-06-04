package request

import (
	"fmt"
	"time"
)

type SubscriptionPurchaseRequest struct {
	ExpiresAt time.Time `json:"expires_at"`
}

func (req *SubscriptionPurchaseRequest) ValidateRequest() (err error) {
	if req.ExpiresAt.Before(time.Now().UTC()) {
		return fmt.Errorf("дата окончания подписки не может быть раньше текущей даты")
	}

	return nil
}
