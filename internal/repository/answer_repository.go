package repository

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"tilimauth/internal/model"
)

type AnswerRepository struct {
	db *sql.DB
}

func NewAnswerRepo(db *sql.DB) *AnswerRepository {
	return &AnswerRepository{db: db}
}

func (r *AnswerRepository) GetByExerciseIDs(exerciseIDs []int64) (map[int64][]model.Answer, error) {
	query := `
		SELECT id, text, image, is_correct, exercise_id
		FROM app.answers
		WHERE exercise_id = ANY($1)
		ORDER BY id
	`

	rows, err := r.db.Query(query, pq.Array(exerciseIDs))
	if err != nil {
		return nil, fmt.Errorf("error fetching answers: %w", err)
	}
	defer rows.Close()

	answersByExercise := make(map[int64][]model.Answer)
	for rows.Next() {
		var answer model.Answer
		var exerciseID int64
		if err := rows.Scan(&answer.ID, &answer.Text, &answer.Image, &answer.IsCorrect, &exerciseID); err != nil {
			return nil, fmt.Errorf("error scanning answer row: %w", err)
		}
		answersByExercise[exerciseID] = append(answersByExercise[exerciseID], answer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating answer rows: %w", err)
	}

	return answersByExercise, nil
}
