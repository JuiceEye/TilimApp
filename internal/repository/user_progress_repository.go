package repository

import (
	"database/sql"
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

func (r *UserProgressRepository) CreateUserProgress(UserID int64) (*model.UserProgress, error) {
	up := &model.UserProgress{
		UserID:                UserID,
		StreakDays:            0,
		XPPoints:              0,
		WordsLearned:          0,
		LessonsDone:           0,
		LastLessonCompletedAt: nil,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	_, err := r.db.Exec(
		`INSERT INTO app.user_progress
		(user_id, streak_days, xp_points, words_learned, lessons_done, last_lesson_completed_at, created_at, updated_at)
		VALUES ($1::INTEGER, $2::INTEGER, $3::INTEGER, $4::INTEGER, $5::INTEGER, $6::TIMESTAMPTZ, $7::TIMESTAMPTZ, $8::TIMESTAMPTZ)`,
		up.UserID, up.StreakDays, up.XPPoints, up.WordsLearned, up.LessonsDone, up.LastLessonCompletedAt, up.CreatedAt, up.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return up, nil
}

func (r *UserProgressRepository) AddXP(userID int64, xp int64) error {
	_, err := r.db.Exec(
		`UPDATE app.user_progress
         SET xp_points     = xp_points + $1, lessons_done = lessons_done + 1, updated_at   = NOW()
         WHERE user_id = $2`,
		xp, userID,
	)
	return err
}
