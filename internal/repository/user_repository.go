package repository

import (
	"database/sql"
	"fmt"
	"tilimauth/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

type UserCredentials struct {
	ID             int64
	HashedPassword string
}

func NewUserRepo(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUserByID(UserID int64) (*model.User, error) {
	rows, err := r.db.Query("SELECT id, username, email, phone_number, registration_date FROM auth.users WHERE id = $1::INTEGER", UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u := new(model.User)
	for rows.Next() {
		u, err = scanRowIntoUsers(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, ErrNotFound
	}

	return u, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	rows, err := r.db.Query("SELECT id, username, email, phone_number, registration_date FROM auth.users WHERE email = $1::TEXT", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u := new(model.User)
	for rows.Next() {
		u, err = scanRowIntoUsers(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, ErrNotFound
	}

	return u, nil
}

func (r *UserRepository) GetUserByPhoneNumber(phoneNumber string) (*model.User, error) {
	rows, err := r.db.Query("SELECT id, username, email, phone_number, registration_date FROM auth.users WHERE phone_number = $1::TEXT", phoneNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u := new(model.User)
	for rows.Next() {
		u, err = scanRowIntoUsers(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, ErrNotFound
	}

	return u, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	rows, err := r.db.Query(
		"SELECT id, username, email, phone_number, registration_date FROM auth.users WHERE username = $1::TEXT",
		username,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u := new(model.User)
	for rows.Next() {
		u, err = scanRowIntoUsers(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, ErrNotFound
	}

	return u, nil
}

func (r *UserRepository) CreateUser(user *model.User) (*model.User, error) {
	err := r.db.QueryRow(
		`INSERT INTO auth.users (username, password, email, phone_number, image, registration_date) 
		VALUES ($1::TEXT, $2::TEXT, $3::TEXT, $4::TEXT, $5::TEXT, $6::TIMESTAMPTZ) RETURNING id`,
		user.Username,
		user.Password,
		user.Email,
		user.PhoneNumber,
		user.Image,
		user.RegistrationDate,
	).Scan(&user.ID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func scanRowIntoUsers(rows *sql.Rows) (*model.User, error) {
	user := new(model.User)

	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PhoneNumber,
		&user.RegistrationDate,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) getCredentials(field, value string) (*UserCredentials, error) {
	query := fmt.Sprintf(
		"SELECT id, password FROM auth.users WHERE %s = $1", field,
	)
	row := r.db.QueryRow(query, value)

	credentials := new(UserCredentials)
	if err := row.Scan(&credentials.ID, &credentials.HashedPassword); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return credentials, nil
}

func (r *UserRepository) GetCredentialsByUsername(username string) (*UserCredentials, error) {
	return r.getCredentials("username", username)
}

func (r *UserRepository) GetCredentialsByEmail(email string) (*UserCredentials, error) {
	return r.getCredentials("email", email)
}
