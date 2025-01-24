package database

import (
	"database/sql"
	"errors"
	"log"
)

type sqliteDatabase struct {
	Db *sql.DB
}

var dbInstance *sqliteDatabase

func NewSqliteDatabase() Database {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}

	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS USERS (id INTEGER PRIMARY KEY, first_name TEXT, last_name TEXT, email TEXT, password TEXT);")
	statement.Exec()

	dbInstance = &sqliteDatabase{Db: db}
	return dbInstance
}

func (s *sqliteDatabase) GetDb() *sql.DB {
	return dbInstance.Db
}

func (s *sqliteDatabase) GetUsers() (*[]User, error) {
	users := []User{}
	rows, err := dbInstance.Db.Query("SELECT id, first_name, last_name, email FROM USERS;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var user User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &users, nil
}

func (s *sqliteDatabase) GetUser(u User) (*User, error) {
	var user User
	err := dbInstance.Db.QueryRow("SELECT id, email, first_name, last_name, password FROM USERS WHERE email = ?", u.Email).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password)
	if err == sql.ErrNoRows {
		log.Printf("%v", err)
		return nil, errors.New("User does not exist")
	}
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &user, nil
}

func (s *sqliteDatabase) CreateUser(u User) (*User, error) {
	result, err := dbInstance.Db.Exec("INSERT INTO USERS (first_name,last_name,email,password) VALUES (?,?,?,?)", u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		log.Printf("Err: %v", err)
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Err: %v", err)
		return nil, err
	}
	u.ID = id
	return &u, nil
}
