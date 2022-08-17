package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type Mood struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Mood      int       `json:"mood"`
	CreatedAt time.Time `json:"created_at"`
}

type Config struct {
	SecretKey string `toml:"secret_key"`
}

var DB *sql.DB

func Init(config Config) error {
	JwtWrapper.SecretKey = config.SecretKey
	JwtWrapper.ExpirationHours = 24
	JwtWrapper.Issuer = "???"

	const (
		host     = "localhost"
		port     = 5435
		user     = "m00d"
		password = "m00d"
		dbname   = "m00d"
	)

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return err
	}

	DB = db

	// init moods table
	if _, err = db.Exec("create table if not exists moods (id serial primary key, user_id text, mood int, created_at timestamptz)"); err != nil {
		return err
	}

	// init users table
	if _, err = db.Exec("create table if not exists users (id serial primary key, email text, name text, created_at timestamptz)"); err != nil {
		return err
	}

	return nil
}
