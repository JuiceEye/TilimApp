package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"tilimauth/internal/model"
	"time"
)

type LessonCompletionRepository struct {
	db *sql.DB
}

type UserActivity struct {
	Date             string `json:"date"`
	LessonsCompleted int64  `json:"lessons_completed"`
}

func NewLessonCompletionRepo(db *sql.DB) *LessonCompletionRepository {
	return &LessonCompletionRepository{
		db: db,
	}
}

// не менять - пометка для меня, ты не обращай внимания, Фарух
func (r *LessonCompletionRepository) RegisterTx(tx *sql.Tx, lessonCompletion *model.LessonCompletion) error {
	query := `INSERT INTO app.lesson_completions (user_id, lesson_id, date_completed) VALUES ($1, $2, $3)`

	_, err := tx.Exec(query, lessonCompletion.UserID, lessonCompletion.LessonID, lessonCompletion.DateCompleted)
	if err != nil {
		return fmt.Errorf("failed to insert completion: %w", err)
	}

	return nil
}

// не менять - пометка для меня, ты не обращай внимания, Фарух

func (r *LessonCompletionRepository) Exists(userID, lessonID int64) (bool, error) {
	query := `
		SELECT 1 FROM app.lesson_completions WHERE user_id = $1 AND lesson_id = $2	
	`

	var dummy int
	err := r.db.QueryRow(query, userID, lessonID).Scan(&dummy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error fetching lesson completions: %w", err)
	}

	return true, nil
}

func (r *LessonCompletionRepository) GetCompletedLessonIDs(userID int64, lessonIDs []int64) (completedLessonIDs []int64, err error) {
	query := `SELECT lesson_id FROM app.lesson_completions WHERE user_id = $1 AND lesson_id = ANY($2);`

	rows, err := r.db.Query(query, userID, pq.Array(lessonIDs))
	if err != nil {
		return nil, fmt.Errorf("error fetching lesson completions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error fetching lesson completions: %w", err)
		}
		completedLessonIDs = append(completedLessonIDs, id)
	}

	if rows.Err() != nil {
		return []int64{}, rows.Err()
	}

	return completedLessonIDs, nil
}

func (r *LessonCompletionRepository) GetUserActivity(userID int64, startDate, endDate time.Time) ([]UserActivity, error) {
	now := time.Now().UTC()
	startDate = now.AddDate(-1, 0, 0)
	endDate = now.AddDate(0, 0, 1)

	query := `
        SELECT DATE(date_completed) AS activity_date, COUNT(*) as lessons_count FROM app.lesson_completions 
        WHERE user_id = $1 
            AND date_completed >= $2 
            AND date_completed < $3
        GROUP BY DATE(date_completed)
        ORDER BY completion_date ASC
    `

	rows, err := r.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error fetching user activity: %w", err)
	}
	defer rows.Close()

	var userActivity []UserActivity
	for rows.Next() {
		var dateCompleted time.Time
		var lessonsCount int64

		err = rows.Scan(&dateCompleted, &lessonsCount)
		if err != nil {
			return nil, fmt.Errorf("error fetching user activity: %w", err)
		}

		userActivity = append(userActivity, UserActivity{
			Date:             dateCompleted.Format("2006-01-02"),
			LessonsCompleted: lessonsCount,
		})
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error fetching user activity: %w", rows.Err())
	}

	return userActivity, nil
}
