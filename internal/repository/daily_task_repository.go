package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"tilimauth/internal/model"
	"time"
)

// Define a struct to hold the data fetched from the SELECT query
type dailyTaskData struct {
	taskID   int64
	lessonID int64
	// Add other fields if needed for the insert, though your INSERT only uses taskID and lessonID from the SELECT
}

type DailyTaskRepository struct {
	db *sql.DB
}

func NewDailyTaskRepository(db *sql.DB) *DailyTaskRepository {
	return &DailyTaskRepository{
		db: db,
	}
}

// GetAllDailyTasks retrieves all daily tasks from the database
func (r *DailyTaskRepository) GetAllDailyTasks() ([]model.DailyTask, error) {
	query := `
		SELECT id, title, description, xp, lesson_id, created_at
		FROM app.daily_tasks
		ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching daily tasks: %w", err)
	}
	defer rows.Close()

	var tasks []model.DailyTask
	for rows.Next() {
		var task model.DailyTask
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.XP, &task.LessonID, &task.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning daily task row: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating daily task rows: %w", err)
	}

	return tasks, nil
}

// GetUserDailyTasks retrieves daily tasks assigned to a user for a specific date
func (r *DailyTaskRepository) GetUserDailyTasks(userID int64, date time.Time) ([]model.UserDailyTask, error) {
	query := `
		SELECT udt.id, udt.user_id, udt.daily_task_id, udt.lesson_id, dt.title, dt.description, dt.xp, 
		       udt.completed, udt.assigned_date, udt.completed_at
		FROM app.user_daily_tasks udt
		JOIN app.daily_tasks dt ON udt.daily_task_id = dt.id
		WHERE udt.user_id = $1 AND DATE(udt.assigned_date) = DATE($2)
		ORDER BY udt.id
	`

	rows, err := r.db.Query(query, userID, date)
	if err != nil {
		return nil, fmt.Errorf("error fetching user daily tasks: %w", err)
	}
	defer rows.Close()

	var tasks []model.UserDailyTask
	for rows.Next() {
		var task model.UserDailyTask
		if err := rows.Scan(
			&task.ID, &task.UserID, &task.DailyTaskID, &task.LessonID,
			&task.Title, &task.Description, &task.XP,
			&task.Completed, &task.AssignedDate, &task.CompletedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning user daily task row: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user daily task rows: %w", err)
	}

	return tasks, nil
}

func (r *DailyTaskRepository) AssignDailyTasksToUser(userID int64, taskIDs []int64, date time.Time) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback() // Ensure rollback on error or panic

	// First, get all the task details and store them
	query := `
		SELECT id, title, description, xp, lesson_id
		FROM app.daily_tasks
		WHERE id = ANY($1)
	`

	rows, err := tx.Query(query, pq.Array(taskIDs))
	if err != nil {
		return fmt.Errorf("error fetching daily tasks: %w", err)
	}
	// It's good practice to defer rows.Close(), but we'll close it explicitly sooner.
	// However, defer still catches returns from errors before explicit close.
	defer rows.Close()

	var tasksToAssign []dailyTaskData
	for rows.Next() {
		var fetchedTaskID, fetchedLessonID, xp int64 // Temporary vars for scanning if not using sql.Null types directly
		var title, description string

		// Assuming lesson_id IS NOT NULL. If it can be, use sql.NullInt64 and handle it.
		// For simplicity, sticking to your original int64 for lessonID, assuming it's NOT NULL.
		// If lesson_id in app.daily_tasks can be NULL, you should scan into sql.NullInt64
		// and then decide how to handle NULL for the INSERT (e.g. pass NULL, or a default value)
		if err := rows.Scan(&fetchedTaskID, &title, &description, &xp, &fetchedLessonID); err != nil {
			// No need to call rows.Close() here if defer rows.Close() is present and we are returning.
			return fmt.Errorf("error scanning daily task: %w", err)
		}
		tasksToAssign = append(tasksToAssign, dailyTaskData{taskID: fetchedTaskID, lessonID: fetchedLessonID})
	}

	if err := rows.Err(); err != nil {
		// No need to call rows.Close() here if defer rows.Close() is present and we are returning.
		return fmt.Errorf("error iterating daily tasks: %w", err)
	}

	// Explicitly close rows after we are done reading from it and before starting new operations.
	// This is crucial to prevent the "unexpected Parse response 'D'" error.
	if err := rows.Close(); err != nil {
		return fmt.Errorf("error closing rows after fetching daily tasks: %w", err)
	}

	// Now, insert each task for the user
	insertQuery := `
		INSERT INTO app.user_daily_tasks (user_id, daily_task_id, lesson_id, completed, assigned_date)
		VALUES ($1, $2, $3, false, $4)
	`

	for _, taskData := range tasksToAssign {
		// The fmt.Println should now reflect the taskData being used
		_, err = tx.Exec(insertQuery, userID, taskData.taskID, taskData.lessonID, date)
		if err != nil {
			// tx.Rollback() will be called by the defer statement
			return fmt.Errorf("error inserting user daily task (taskID %d): %w", taskData.taskID, err)
		}
	}

	return tx.Commit()
}

// MarkTaskCompleted marks a daily task as completed
func (r *DailyTaskRepository) MarkTaskCompletedTx(tx *sql.Tx, userID, lessonID int64, completedAt time.Time) error {
	query := `
		UPDATE app.user_daily_tasks
		SET completed = true, completed_at = $1
		WHERE user_id = $2 AND lesson_id = $3 AND DATE(assigned_date) = DATE($1) AND completed = false
	`

	result, err := tx.Exec(query, completedAt, userID, lessonID)
	if err != nil {
		return fmt.Errorf("error updating user daily task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		// No task was updated, which is fine - it might not be a daily task
		return nil
	}

	return nil
}

// HasUserDailyTasksForDate checks if a user has daily tasks assigned for a specific date
func (r *DailyTaskRepository) HasUserDailyTasksForDate(userID int64, date time.Time) (bool, error) {
	query := `
		SELECT 1
		FROM app.user_daily_tasks
		WHERE user_id = $1 AND DATE(assigned_date) = DATE($2)
		LIMIT 1
	`

	var dummy int
	err := r.db.QueryRow(query, userID, date).Scan(&dummy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error checking user daily tasks: %w", err)
	}

	return true, nil
}
