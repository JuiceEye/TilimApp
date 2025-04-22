package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"tilimauth/internal/model"
)

type ExerciseRepository struct {
	db *sql.DB
}

func NewExerciseRepo(db *sql.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (r *ExerciseRepository) GetByLessonID(lessonID int64) (*model.Exercise, error) {
	query := `
		SELECT id, text, image, question_text
		FROM app.exercises
		WHERE lesson_id = $1
	`

	var exercise model.Exercise
	err := r.db.QueryRow(query, lessonID).Scan(&exercise.ID, &exercise.Text, &exercise.Image, &exercise.QuestionText)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error fetching exercise: %w", err)
	}

	return &exercise, nil
}
