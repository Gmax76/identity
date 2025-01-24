package database

import "database/sql"

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  []byte `json:"password"`
}

type Database interface {
	GetDb() *sql.DB
	GetUsers() (*[]User, error)
	GetUser(User) (*User, error)
	CreateUser(User) (*User, error)
}
