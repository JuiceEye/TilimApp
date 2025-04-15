package repository

import (
	"database/sql"
	"tilimauth/internal/model"
)

type UserProgressRepository struct {
	db *sql.DB
}

func NewProfileRepo(db *sql.DB) *UserProgressRepository {
	return &UserProgressRepository{
		db: db,
	}
}

func (r *UserProgressRepository) GetUserProgressByUserID(UserID int64) (*model.UserProgress, error) {
	rows, err := r.db.Query("SELECT * FROM tilim.user_progress WHERE user_id = $1::INTEGER", UserID)
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

func scanRowIntoUserProgress(rows *sql.Rows) (*model.UserProgress, error) {
	userProgress := new(model.UserProgress)

	err := rows.Scan(
		&userProgress.UserID,
		&userProgress.Streak,
		&userProgress.XPPoints,
		&userProgress.WordsLearned,
		&userProgress.LessonsDone,
		&userProgress.LastLessonCompletedAt,
		&userProgress.CreatedAt,
		&userProgress.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return userProgress, nil
}
