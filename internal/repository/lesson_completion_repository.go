package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"tilimauth/internal/model"
	"time"
)

type LessonCompletionRepository struct {
	db *sql.DB
}

func NewLessonCompletionRepo(db *sql.DB) *LessonCompletionRepository {
	return &LessonCompletionRepository{
		db: db,
	}
}

func (r *LessonCompletionRepository) Register(lessonCompletion *model.LessonCompletion) error {
	query := `INSERT INTO app.lesson_completions (user_id, lesson_id, date_completed) VALUES ($1, $2, $3)`

	_, err := r.db.Exec(query, lessonCompletion.UserID, lessonCompletion.LessonID, lessonCompletion.DateCompleted)
	if err != nil {
		return fmt.Errorf("failed to insert completion: %w", err)
	}

	return nil
}

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
