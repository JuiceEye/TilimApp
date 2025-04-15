package repository

import (
	"database/sql"
	"tilimauth/internal/model"
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepo(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

func (r *ProfileRepository) GetProfileByID(id int64) (*model.Profile, error) {
	rows, err := r.db.Query("SELECT id, username, registration_date FROM auth.users WHERE id = $1::INTEGER", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	p := new(model.Profile)
	for rows.Next() {
		p, err = scanRowIntoProfiles(rows)
		if err != nil {
			return nil, err
		}
	}

	if p.UserID == 0 {
		return nil, ErrNotFound
	}

	return p, nil
}

func scanRowIntoProfiles(rows *sql.Rows) (*model.Profile, error) {
	profile := new(model.Profile)

	err := rows.Scan(
		&profile.UserID,
		&profile.Username,
		&profile.RegistrationDate,
	)

	if err != nil {
		return nil, err
	}

	return profile, nil
}
