package service

import (
	"database/sql"
	"math/rand"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
	"time"
)

const (
	DailyTasksCount = 3 // Number of daily tasks to assign to a user each day
)

type DailyTaskService struct {
	dailyTaskRepo *repository.DailyTaskRepository
}

func NewDailyTaskService(dailyTaskRepo *repository.DailyTaskRepository) *DailyTaskService {
	return &DailyTaskService{
		dailyTaskRepo: dailyTaskRepo,
	}
}

// GetUserDailyTasks retrieves daily tasks for a user for the current day
// If the user doesn't have daily tasks assigned for today, it assigns new ones
func (s *DailyTaskService) GetUserDailyTasks(userID int64) ([]model.UserDailyTask, error) {
	today := time.Now().Truncate(24 * time.Hour)

	// Check if user already has daily tasks for today
	hasTasks, err := s.dailyTaskRepo.HasUserDailyTasksForDate(userID, today)
	if err != nil {
		return nil, err
	}

	// If user doesn't have daily tasks for today, assign new ones
	if !hasTasks {
		if err := s.assignDailyTasksToUser(userID); err != nil {
			return nil, err
		}
	}

	// Get user's daily tasks for today
	return s.dailyTaskRepo.GetUserDailyTasks(userID, today)
}

// assignDailyTasksToUser assigns random daily tasks to a user
func (s *DailyTaskService) assignDailyTasksToUser(userID int64) error {
	// Get all available daily tasks
	allTasks, err := s.dailyTaskRepo.GetAllDailyTasks()
	if err != nil {
		return err
	}

	// If there are no tasks, return
	if len(allTasks) == 0 {
		return nil
	}

	// Shuffle the tasks
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allTasks), func(i, j int) {
		allTasks[i], allTasks[j] = allTasks[j], allTasks[i]
	})

	// Select the first DailyTasksCount tasks (or all if there are fewer)
	count := DailyTasksCount
	if len(allTasks) < count {
		count = len(allTasks)
	}

	// Get the IDs of the selected tasks
	taskIDs := make([]int64, count)
	for i := 0; i < count; i++ {
		taskIDs[i] = allTasks[i].ID
	}

	// Assign the tasks to the user
	return s.dailyTaskRepo.AssignDailyTasksToUser(userID, taskIDs, time.Now())
}

// CheckAndMarkTaskCompleted checks if a completed lesson is a daily task and marks it as completed
func (s *DailyTaskService) CheckAndMarkTaskCompletedTx(tx *sql.Tx, userID, lessonID int64, completedAt time.Time) error {
	// This method is called from the lesson completion service
	// It checks if the completed lesson is a daily task and marks it as completed
	return s.dailyTaskRepo.MarkTaskCompletedTx(tx, userID, lessonID, completedAt)
}
