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
		SELECT e.id, et.code AS type_code, COALESCE(e.text, ''), COALESCE(e.image, ''), COALESCE(e.question_text, ''), COALESCE(e.audio_uuid, '')
		FROM app.exercises e
		INNER JOIN dict.exercise_types et ON e.type_id = et.id
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
		var audioUUID sql.NullString
		if err := rows.Scan(&exercise.ID, &exercise.TypeCode, &exercise.Text, &exercise.Image, &exercise.QuestionText, &audioUUID); err != nil {
			return nil, fmt.Errorf("error scanning section row: %w", err)
		}
		if audioUUID.Valid && audioUUID.String != "" {
			exercise.Audio = &model.File{UUID: audioUUID.String}
		} else {
			exercise.Audio = nil
		}
		exercises = append(exercises, exercise)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating exercise rows: %w", err)
	}

	return exercises, nil
}
