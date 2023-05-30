package models

import "time"

type Comment struct {
	Id         int
	PostId     int
	Creator    string
	Text       string
	Likes      int
	Dislikes   int
	IsAuth     bool
	Created_at time.Time
}
