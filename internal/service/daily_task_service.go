package service

import (
	"database/sql"
	"math/rand"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
	"time"
)

const (
	DailyTasksCount = 3
)

type DailyTaskService struct {
	dailyTaskRepo *repository.DailyTaskRepository
}

func NewDailyTaskService(dailyTaskRepo *repository.DailyTaskRepository) *DailyTaskService {
	return &DailyTaskService{
		dailyTaskRepo: dailyTaskRepo,
	}
}

func (s *DailyTaskService) GetUserDailyTasks(userID int64) ([]model.UserDailyTask, error) {
	today := time.Now().Truncate(24 * time.Hour)

	hasTasks, err := s.dailyTaskRepo.HasUserDailyTasksForDate(userID, today)
	if err != nil {
		return nil, err
	}

	if !hasTasks {
		if err := s.assignDailyTasksToUser(userID); err != nil {
			return nil, err
		}
	}

	return s.dailyTaskRepo.GetUserDailyTasks(userID, today)
}

func (s *DailyTaskService) assignDailyTasksToUser(userID int64) error {
	allTasks, err := s.dailyTaskRepo.GetAllDailyTasks()
	if err != nil {
		return err
	}

	if len(allTasks) == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allTasks), func(i, j int) {
		allTasks[i], allTasks[j] = allTasks[j], allTasks[i]
	})

	count := DailyTasksCount
	if len(allTasks) < count {
		count = len(allTasks)
	}

	taskIDs := make([]int64, count)
	for i := 0; i < count; i++ {
		taskIDs[i] = allTasks[i].ID
	}

	return s.dailyTaskRepo.AssignDailyTasksToUser(userID, taskIDs, time.Now())
}

func (s *DailyTaskService) CheckAndMarkTaskCompletedTx(tx *sql.Tx, userID, lessonID int64, completedAt time.Time) error {
	return s.dailyTaskRepo.MarkTaskCompletedTx(tx, userID, lessonID, completedAt)
}
