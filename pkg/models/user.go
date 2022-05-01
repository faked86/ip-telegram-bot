package models

type User struct {
	ID       uint
	UserName string
	History  []string
	Admin    bool
}
