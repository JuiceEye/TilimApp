package request

import (
	"fmt"
	"time"
)

type SubscriptionPurchaseRequest struct {
	ExpiresAt time.Time `json:"expires_at"`
}

func (req *SubscriptionPurchaseRequest) ValidateRequest() (err error) {
	t := req.ExpiresAt

	if t.Before(time.Now().UTC()) {
		return fmt.Errorf("дата окончания подписки не может быть раньше текущей даты")
	}

	req.ExpiresAt = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	return nil
}
