package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"tilimauth/internal/model"
)

type ModuleRepository struct {
	db *sql.DB
}

func NewModuleRepo(db *sql.DB) *ModuleRepository {
	return &ModuleRepository{
		db: db,
	}
}

func (r *ModuleRepository) GetByID(id int64) (*model.Module, error) {
	query := `
		SELECT id, title
		FROM app.modules
		WHERE id = $1
	`

	var module model.Module
	err := r.db.QueryRow(query, id).Scan(&module.ID, &module.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error fetching module: %w", err)
	}

	return &module, nil
}
