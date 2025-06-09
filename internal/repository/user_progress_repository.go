package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"tilimauth/internal/model"
	"time"
)

type UserProgressRepository struct {
	db *sql.DB
}

func NewUserProgressRepo(db *sql.DB) *UserProgressRepository {
	return &UserProgressRepository{
		db: db,
	}
}

type DBExecutor interface {
	QueryRow(query string, args ...any) *sql.Row
}

func truncateOrNil(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.Truncate(24 * time.Hour)
}

func (r *UserProgressRepository) GetUserProgressByUserID(UserID int64) (*model.UserProgress, error) {
	return r.getUserProgressByUserID(r.db, UserID)
}

func (r *UserProgressRepository) GetUserProgressByUserIDTx(tx *sql.Tx, UserID int64) (*model.UserProgress, error) {
	return r.getUserProgressByUserID(tx, UserID)
}

func (r *UserProgressRepository) getUserProgressByUserID(executor DBExecutor, UserID int64) (*model.UserProgress, error) {
	up := &model.UserProgress{}

	err := executor.QueryRow(
		"SELECT user_id, streak_days, xp_points, words_learned, lessons_done, "+
			"last_lesson_completed_at, updated_at, last_streak_reset_date "+
			"FROM app.user_progress "+
			"WHERE user_id = $1::INTEGER",
		UserID,
	).Scan(
		&up.UserID,
		&up.StreakDays,
		&up.XPPoints,
		&up.WordsLearned,
		&up.LessonsDone,
		&up.LastLessonCompletedAt,
		&up.UpdatedAt,
		&up.LastStreakResetDate,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if up.UserID == 0 {
		return nil, ErrNotFound
	}

	return up, nil
}

func (r *UserProgressRepository) CreateUserProgress(UserID int64) (*model.UserProgress, error) {
	up := &model.UserProgress{
		UserID:                UserID,
		StreakDays:            0,
		XPPoints:              0,
		WordsLearned:          0,
		LessonsDone:           0,
		LastLessonCompletedAt: nil,
		UpdatedAt:             time.Now().UTC(),
		LastStreakResetDate:   nil,
	}

	_, err := r.db.Exec(
		`INSERT INTO app.user_progress
		(user_id, streak_days, xp_points, words_learned, lessons_done, last_lesson_completed_at, updated_at, last_streak_reset_date)
		VALUES ($1::INTEGER, $2::INTEGER, $3::INTEGER, $4::INTEGER, $5::INTEGER, $6::TIMESTAMPTZ, $7::TIMESTAMPTZ, $8::DATE)`,
		up.UserID, up.StreakDays, up.XPPoints, up.WordsLearned, up.LessonsDone, truncateOrNil(up.LastLessonCompletedAt), up.UpdatedAt, truncateOrNil(up.LastStreakResetDate),
	)

	if err != nil {
		return nil, err
	}

	return up, nil
}

func (r *UserProgressRepository) SaveStreakTx(tx *sql.Tx, userID int64, up *model.UserProgress) error {
	query := `
		UPDATE app.user_progress SET 
        	streak_days = $1,
			last_lesson_completed_at = $2
		WHERE user_id = $3
	`

	_, err := tx.Exec(query, up.StreakDays, up.LastLessonCompletedAt, userID)

	if err != nil {
		return fmt.Errorf("failed to save user streak: %w", err)
	}

	return nil
}

func (r *UserProgressRepository) GetUserActivity(userID int64, startDate, endDate time.Time) ([]UserActivity, error) {
	now := time.Now().UTC()
	startDate = now.AddDate(-1, 0, 0)
	endDate = now.AddDate(0, 0, 1)

	query := `
        SELECT DATE(last_lesson_completed_at) AS activity_date, COUNT(*) as lessons_count FROM app.user_progress 
        WHERE user_id = $1 
            AND last_lesson_completed_at >= $2 
            AND last_lesson_completed_at < $3
        GROUP BY DATE(last_lesson_completed_at)
        ORDER BY activity_date ASC
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
