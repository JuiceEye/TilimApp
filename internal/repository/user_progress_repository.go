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

func (r *UserProgressRepository) GetUserProgressByUserID(UserID int64) (*model.UserProgress, error) {
	rows, err := r.db.Query("SELECT * FROM app.user_progress WHERE user_id = $1::INTEGER", UserID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	up := new(model.UserProgress)

	for rows.Next() {
		up, err = scanRowIntoUserProgress(rows)

		if err != nil {
			return nil, err
		}
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
		(user_id, streak_days, xp_points, words_learned, lessons_done, last_lesson_completed_at, updated_at, last_lesson_completed_at)
		VALUES ($1::INTEGER, $2::INTEGER, $3::INTEGER, $4::INTEGER, $5::INTEGER, $6::TIMESTAMPTZ, $7::TIMESTAMPTZ, $8::DATE)`,
		up.UserID, up.StreakDays, up.XPPoints, up.WordsLearned, up.LessonsDone, up.LastLessonCompletedAt, up.UpdatedAt, up.LastStreakResetDate,
	)

	if err != nil {
		return nil, err
	}

	return up, nil
}

func scanRowIntoUserProgress(rows *sql.Rows) (*model.UserProgress, error) {
	userProgress := new(model.UserProgress)

	err := rows.Scan(
		&userProgress.UserID,
		&userProgress.StreakDays,
		&userProgress.XPPoints,
		&userProgress.WordsLearned,
		&userProgress.LessonsDone,
		&userProgress.LastLessonCompletedAt,
		&userProgress.UpdatedAt,
		&userProgress.LastStreakResetDate,
	)

	if err != nil {
		return nil, err
	}

	return userProgress, nil
}
