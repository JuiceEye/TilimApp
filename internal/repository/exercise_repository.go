package repository

import (
	"database/sql"
	"fmt"
	"tilimauth/internal/model"
)

type ExerciseRepository struct {
	db *sql.DB
}

func NewExerciseRepo(db *sql.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (r *ExerciseRepository) GetByLessonID(lessonID int64) ([]model.Exercise, error) {
	query := `
		SELECT id, COALESCE(text, ''), COALESCE(image, ''), question_text
		FROM app.exercises
		WHERE lesson_id = $1
	`

	rows, err := r.db.Query(query, lessonID)
	if err != nil {
		return nil, fmt.Errorf("error fetching exercises: %w", err)
	}
	defer rows.Close()

	var exercises []model.Exercise
	for rows.Next() {
		var exercise model.Exercise
		if err := rows.Scan(&exercise.ID, &exercise.Text, &exercise.Image, &exercise.QuestionText); err != nil {
			return nil, fmt.Errorf("error scanning section row: %w", err)
		}
		exercises = append(exercises, exercise)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating exercise rows: %w", err)
	}

	return exercises, nil
}
