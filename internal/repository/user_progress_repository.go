package repository

import (
	"database/sql"
	"tilimauth/internal/model"
	"time"
)

type UserProgressRepository struct {
	db *sql.DB
}

type DBExecutor interface {
	QueryRow(query string, args ...any) *sql.Row
}

func NewUserProgressRepo(db *sql.DB) *UserProgressRepository {
	return &UserProgressRepository{
		db: db,
	}
}

func (r *UserProgressRepository) GetUserProgressByUserID(UserID int64) (*model.UserProgress, error) {
	return r.getUserProgressByUserID(r.db, UserID)
}

func (r *UserProgressRepository) GetUserProgressByUserIDTx(tx *sql.Tx, UserID int64) (*model.UserProgress, error) {
	return r.getUserProgressByUserID(tx, UserID)
}

func (r *UserProgressRepository) getUserProgressByUserID(executor DBExecutor, UserID int64) (*model.UserProgress, error) {
	up := &model.UserProgress{}

	err := executor.QueryRow("SELECT * FROM app.user_progress WHERE user_id = $1::INTEGER", UserID).Scan(
		&up.UserID,
		&up.StreakDays,
		&up.XPPoints,
		&up.WordsLearned,
		&up.LessonsDone,
		&up.LastLessonCompletedAt,
		&up.UpdatedAt,
		&up.LastStreakResetDate,
	)

	if err != nil {
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
		UpdatedAt:             time.Now(),
		LastStreakResetDate:   nil,
	}

	_, err := r.db.Exec(
		`INSERT INTO app.user_progress
		(user_id, streak_days, xp_points, words_learned, lessons_done, last_lesson_completed_at, updated_at, last_streak_reset_date)
		VALUES ($1::INTEGER, $2::INTEGER, $3::INTEGER, $4::INTEGER, $5::INTEGER, $6::TIMESTAMPTZ, $7::TIMESTAMPTZ, $8::DATE)`,
		up.UserID, up.StreakDays, up.XPPoints, up.WordsLearned, up.LessonsDone, up.LastLessonCompletedAt, up.UpdatedAt, up.LastStreakResetDate,
	)

	if err != nil {
		return nil, err
	}

	return up, nil
}
