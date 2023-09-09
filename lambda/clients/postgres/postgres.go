package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type UserStorageDB interface {
	GetUser(username string) (*User, error)
	ValidateRefreshToken(username, refreshToken string) bool
	AddUserToDB(username, password, hashedToken string) error
	UpdateUserToken(username, token string) error
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RDSClient struct {
	db *sql.DB
}

const (
	// Demo connection string
	// connectionString = "postgresql://postgres:postgres@localhost:5431/postgres?schema=public"
	connectionString = "ragstackcdkstack-rdsdatabaseda351f35-pzt9quiob2l7.ca9nnjlv85uj.us-west-2.rds.amazonaws.com"
)

// Should this be a RDS Client?
// Or should it be postgres client?
func NewRDSClient() (*RDSClient, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return &RDSClient{db: db}, nil
}

func (r *RDSClient) Close() {
	r.db.Close()
}

func (db *RDSClient) GetUser(username string) (*User, error) {
	var user User
	row := db.db.QueryRow("SELECT username, password FROM users WHERE username = $1", username)

	if err := row.Scan(&user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No users found in db")
		}
		return nil, err
	}
	return &user, nil
}
