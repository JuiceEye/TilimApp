package repository

import (
	"database/sql"
	"fmt"
	"tilimauth/internal/model"
)

type SectionRepository struct {
	db *sql.DB
}

func NewSectionRepo(db *sql.DB) *SectionRepository {
	return &SectionRepository{db: db}
}

func (r *SectionRepository) GetByModuleID(moduleID int64) ([]model.Section, error) {
	query := `
		SELECT id, title
		FROM app.sections
		WHERE module_id = $1
		ORDER BY id
	`

	rows, err := r.db.Query(query, moduleID)
	if err != nil {
		return nil, fmt.Errorf("error fetching sections: %w", err)
	}
	defer rows.Close()

	var sections []model.Section
	for rows.Next() {
		var section model.Section
		if err := rows.Scan(&section.ID, &section.Title); err != nil {
			return nil, fmt.Errorf("error scanning section row: %w", err)
		}
		sections = append(sections, section)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating section rows: %w", err)
	}

	return sections, nil
}
