package main

import "time"

type User struct {
	ID         int64
	ChatID     int64
	Username   string
	Name       string
	LastCheck  time.Time
	UserLimit  int
	AlreadyGot int
	Lang       string
	Game       string
}

type Code struct {
	ID     int
	Code   string
	Used   bool
	UsedBy int64
}
