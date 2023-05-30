package models

type Like struct {
	UserID       int
	PostID       int
	Islike       int
	CommentID    int
	CountLike    int
	Countdislike int
}
