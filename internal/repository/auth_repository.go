package repository

import (
	"database/sql"
	"errors"
	"net/http"
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

func (r *AuthRepository) GetUserByEmail(email string) (u *model.User, code int, err error) {
	rows, err := r.db.Query("SELECT * FROM auth.users WHERE email = $1::TEXT", email)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	u = new(model.User)
	for rows.Next() {
		u, err = scanRowIntoUsers(rows)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	if u.Id == 0 {
		return nil, http.StatusNotFound, errors.New("user not found")
	}

	return u, http.StatusOK, nil
}

func (r *AuthRepository) GetUserByPhoneNumber(phoneNumber string) (u *model.User, status int, err error) {
	rows, err := r.db.Query("SELECT * FROM auth.users WHERE phone_number = $1::TEXT", phoneNumber)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	u = new(model.User)
	for rows.Next() {
		u, err = scanRowIntoUsers(rows)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	if u.Id == 0 {
		return nil, http.StatusNotFound, errors.New("user not found")
	}

	return u, http.StatusOK, nil
}

func (r *AuthRepository) GetUserByUsername(username string) (u *model.User, status int, err error) {
	rows, err := r.db.Query("SELECT * FROM auth.users WHERE username = $1::TEXT", username)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	u = new(model.User)
	for rows.Next() {
		u, err = scanRowIntoUsers(rows)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	if u.Id == 0 {
		return nil, http.StatusBadRequest, errors.New("user not found")
	}

	return u, http.StatusOK, nil
}

func (r *AuthRepository) CreateUser(user *model.User) (*model.User, error) {
	err := r.db.QueryRow(
		"INSERT INTO auth.users (username, password, email, phone_number, image, registration_date) "+
			"VALUES ($1::TEXT, $2::TEXT, $3::TEXT, $4::TEXT, $5::TEXT, $6::TIMESTAMPTZ) RETURNING id",
		user.Username,
		user.Password,
		user.Email,
		user.PhoneNumber,
		user.Image,
		user.RegistrationDate,
	).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func scanRowIntoUsers(rows *sql.Rows) (*model.User, error) {
	user := new(model.User)

	err := rows.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.PhoneNumber,
		&user.Image,
		&user.RegistrationDate,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
