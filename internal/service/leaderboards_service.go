package service

import (
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
)

type LeaderboardsService struct {
	userRepository         *repository.UserRepository
	userProgressRepository *repository.UserProgressRepository
}

func NewLeaderboardsService(
	userRepository *repository.UserRepository,
	userProgressRepository *repository.UserProgressRepository,
) *LeaderboardsService {
	return &LeaderboardsService{
		userRepository:         userRepository,
		userProgressRepository: userProgressRepository,
	}
}

func (s *LeaderboardsService) GetLeaderboards() (profiles []*model.LeaderboardsUser, err error) {
	leaderboardsUsers, err := s.userRepository.GetLeaderboardsUsersByID()
	if err != nil {
		return nil, err
	}

	return leaderboardsUsers, nil
}
