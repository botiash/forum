package models

import "time"

type Post struct {
	Id          int
	Title       string
	Description string
	Category    []string
	Author      string
	Likes       int
	Dislikes    int
	CreateAt    time.Time
}
