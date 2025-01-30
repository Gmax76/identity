package database

import (
	"database/sql"

	"github.com/Gmax76/identity/entity"
)

type Database interface {
	GetDb() *sql.DB
	GetUsers() (*[]entity.User, error)
	GetUser(entity.User) (*entity.User, error)
	CreateUser(entity.User) (*entity.User, error)
}
