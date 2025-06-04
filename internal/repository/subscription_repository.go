package repository

import (
	"database/sql"
	"tilimauth/internal/model"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepo(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

func (r *SubscriptionRepository) AddUserSubscription(sub *model.Subscription) (int64, error) {
	err := r.db.QueryRow(
		`INSERT INTO app.user_subscriptions
		(user_id, expires_at, created_at)
		VALUES ($1::INTEGER, $2::TIMESTAMPTZ, $3::TIMESTAMPTZ) RETURNING id`,
		sub.UserID, sub.ExpiresAt, sub.CreatedAt,
	).Scan(&sub.ID)

	if err != nil {
		return 0, err
	}

	return sub.ID, nil
}

func (r *SubscriptionRepository) ExistsActive(userID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM app.user_subscriptions 
			WHERE user_id = $1 AND expires_at > NOW()
		)`, userID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
