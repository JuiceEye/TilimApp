package service

import (
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
	"time"
)

type SubscriptionService struct {
	userRepo *repository.UserRepository
	subRepo  *repository.SubscriptionRepository
}

func NewSubscriptionService(
	userRepo *repository.UserRepository,
	subRepo *repository.SubscriptionRepository,
) *SubscriptionService {
	return &SubscriptionService{
		userRepo: userRepo,
		subRepo:  subRepo,
	}
}

func (s *SubscriptionService) BuySubscription(userID int64, expiresAt time.Time) (int64, error) {
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return 0, err
	}

	sub := &model.Subscription{
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().UTC().Truncate(24 * time.Hour),
	}

	return s.subRepo.AddUserSubscription(sub)
}
