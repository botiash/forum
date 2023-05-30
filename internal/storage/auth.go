package storage

import (
	"database/sql"
	"forum/internal/models"
	"time"
)

type Auth interface {
	CreateUser(user models.User) error
	GetUserByUsername(username string) (models.User, error)
	SaveToken(token string, expired time.Time, username string) error
	GetPasswordByUsername(username string) (string, error)
	DeleteToken(token string) error
}

type AuthStorage struct {
	db *sql.DB
}

func NewAuthStorage(db *sql.DB) *AuthStorage {
	return &AuthStorage{
		db: db,
	}
}

func (a *AuthStorage) CreateUser(user models.User) error {
	query := `INSERT INTO user(email, username, password) VALUES ($1, $2, $3);`
	_, err := a.db.Exec(query, user.Email, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthStorage) GetUserByUsername(username string) (models.User, error) {
	query := `SELECT id, email, username FROM user WHERE username = $1;`
	row := a.db.QueryRow(query, username)
	var user models.User
	if err := row.Scan(&user.Id, &user.Email, &user.Username); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (a *AuthStorage) SaveToken(token string, expired time.Time, username string) error {
	query := `UPDATE user SET session_token = $1, expiresAt = $2 WHERE username = $3;`
	if _, err := a.db.Exec(query, token, expired, username); err != nil {
		return err
	}
	return nil
}

func (a *AuthStorage) GetPasswordByUsername(username string) (string, error) {
	query := `SELECT password FROM user WHERE username = $1;`
	row := a.db.QueryRow(query, username)
	var password string
	if err := row.Scan(&password); err != nil {
		return password, err
	}
	return password, nil
}

func (a *AuthStorage) DeleteToken(token string) error {
	query := `UPDATE user SET session_token = NULL, expiresAt = NULL WHERE session_token = $1`
	if _, err := a.db.Exec(query, token); err != nil {
		return err
	}
	return nil
}
