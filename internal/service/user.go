package service

import (
	"forum/internal/models"
	"forum/internal/storage"
)

type User interface {
	GetUserByToken(token string) (models.User, error)
}

type UserService struct {
	storage *storage.Storage
}

func NewUserService(storage *storage.Storage) *UserService {
	return &UserService{
		storage: storage,
	}
}

func (u *UserService) GetUserByToken(token string) (models.User, error) {
	return u.storage.User.GetUserByToken(token)
}
