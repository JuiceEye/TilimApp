package request

import (
	"fmt"
	"strings"
	"time"
)

type SubscriptionPurchaseRequest struct {
	ExpiresAt CustomDate `json:"expires_at"`
}

func (req *SubscriptionPurchaseRequest) ValidateRequest() (err error) {
	t := req.ExpiresAt

	if t.Before(time.Now().UTC()) {
		return fmt.Errorf("дата окончания подписки не может быть раньше текущей даты")
	}

	req.ExpiresAt = CustomDate{
		Time: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()),
	}
	return nil
}

type CustomDate struct {
	time.Time
}

func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)
	if str == "" {
		return nil
	}

	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return fmt.Errorf("неправильный формат даты: %w", err)
	}

	cd.Time = t
	return nil
}

func (cd *CustomDate) ToTime() time.Time {
	return cd.Time
}
