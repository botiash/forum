package storage

import (
	"database/sql"

	"forum/internal/models"
)

type User interface {
	GetUserByToken(token string) (models.User, error)
}

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (u *UserStorage) GetUserByToken(token string) (models.User, error) {
	query := `SELECT id, email, username, expiresAt FROM user WHERE session_token = $1;`
	row := u.db.QueryRow(query, token)
	var user models.User
	if err := row.Scan(&user.Id, &user.Email, &user.Username, &user.ExpiresAt); err != nil {
		return models.User{}, err
	}
	return user, nil
}
