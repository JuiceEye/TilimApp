package repository

import (
	"database/sql"
	"errors"
	"tilimauth/internal/model"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) GetUserByEmail(email string) (*model.User, error) {
	rows, err := r.db.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	u := new(model.User)
	for rows.Next() {
		u, err = scanRowIntoUsers(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.Id == 0 {
		return nil, errors.New("user not found")
	}

	return u, nil
}

func scanRowIntoUsers(rows *sql.Rows) (*model.User, error) {
	user := new(model.User)

	err := rows.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Phone,
		&user.Image,
		&user.RegistrationDate,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
