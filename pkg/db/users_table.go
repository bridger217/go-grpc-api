package db

import (
	"database/sql"
	"log"
	"time"

	v1 "github.com/bridger217/go-grpc-api/pkg/api/v1"
	_ "github.com/go-sql-driver/mysql"
)

// Internal storage type.
type User struct {
	Id        string
	Username  string
	FirstName string
	LastName  string
	Created   time.Time // Ignored on Create/Update
	Updated   time.Time // Ignored on Create/Update
}

// Convert to API type.
func (u *User) ToExternal() *v1.User {
	return &v1.User{Id: u.Id, Username: u.Username, FirstName: u.FirstName, LastName: u.LastName}
}

type UsersTable struct {
	db *sql.DB
}

func NewUsersTable(db *sql.DB) *UsersTable {
	return &UsersTable{db: db}
}

func (t *UsersTable) InsertUser(id string, username string, firstName string, lastName string) error {
	_, err := t.db.Exec(
		`INSERT INTO Users
		(
			Id, Username, FirstName, LastName
		)
		VALUES
		(
			?, ?, ?, ?
		)`, id, username, firstName, lastName)

	if err != nil {
		log.Printf("Error %s when inserting row into Users table", err)
		return err
	}

	return nil
}

func (t *UsersTable) GetUser(id string) (*User, error) {
	var u User
	err := t.db.QueryRow(
		"SELECT * FROM Users where Id = ?", id).Scan(
		&u.Id, &u.Username, &u.FirstName,
		&u.LastName, &u.Created, &u.Updated)

	if err != nil {
		log.Printf("Error %s when reading row into Users table", err)
		return nil, err
	}

	return &u, nil
}
