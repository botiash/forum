package storage

import "database/sql"

type Storage struct {
	Auth
	PostIR
	User
	CommentIR
	ReactionIR
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Auth:       NewAuthStorage(db),
		PostIR:     NewPostStorage(db),
		User:       NewUserStorage(db),
		CommentIR:  newCommentStorage(db),
		ReactionIR: NewEmotionSQL(db),
	}
}
