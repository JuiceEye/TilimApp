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
	subRepo *repository.SubscriptionRepository,
	userRepo *repository.UserRepository,
) *SubscriptionService {
	return &SubscriptionService{
		subRepo:  subRepo,
		userRepo: userRepo,
	}
}

func (s *SubscriptionService) BuySubscription(userID int64, expiresAt time.Time) (int64, error) {
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return 0, err
	}

	exists, err := s.subRepo.ExistsActive(userID)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, &BadRequestError{Msg: "у пользователя уже есть подписка"}
	}

	sub := &model.Subscription{
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().UTC().Truncate(24 * time.Hour),
	}

	return s.subRepo.AddUserSubscription(sub)
}
