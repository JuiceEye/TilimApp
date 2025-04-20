package repository

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"tilimauth/internal/model"
)

type LessonRepository struct {
	db *sql.DB
}

func NewLessonRepository(db *sql.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

func (r *LessonRepository) GetBySectionIDs(sectionIDs []int64) (map[int64][]model.Lesson, error) {
	query := `
		SELECT id, title, xp, section_id
		FROM app.lessons
		WHERE section_id = ANY($1)
		ORDER BY id
	`

	rows, err := r.db.Query(query, pq.Array(sectionIDs))
	if err != nil {
		return nil, fmt.Errorf("error fetching lessons: %w", err)
	}
	defer rows.Close()

	lessonsBySection := make(map[int64][]model.Lesson)
	for rows.Next() {
		var lesson model.Lesson
		var sectionID int64
		if err := rows.Scan(&lesson.ID, &lesson.Title, &lesson.XP, &sectionID); err != nil {
			return nil, fmt.Errorf("error scanning lesson row: %w", err)
		}
		lessonsBySection[sectionID] = append(lessonsBySection[sectionID], lesson)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating lesson rows: %w", err)
	}

	return lessonsBySection, nil
}
