package model

type User struct {
	ID   int64  `db:"id"`
	User string `db:"user"`
	Pass string `db:"pass"`
}
