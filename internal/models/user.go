package models

import "time"

type User struct {
	Id             int
	Email          string
	Username       string
	Password       string
	RepeatPassword string
	ExpiresAt      time.Time
	IsAuth         bool
}
