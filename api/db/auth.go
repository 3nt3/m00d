package db

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func CheckEmailExists(email string) (bool, error) {
	rows, err := DB.Query("SELECT 1 FROM users WHERE email = $1", email)
	if err != nil {
		return false, err
	}

	// if there is no user associated to the email, there will be no rows so rows.Next() will be false
	return rows.Next(), nil
}

func CreateUser(user User) (int, error) {
	row := DB.QueryRow("insert into users (email, name, created_at) values ($1, $2, $3) returning id", user.Email, user.Name, time.Now().UTC())
	if row.Err() != nil {
		return 0, row.Err()
	}

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func GetUserByEmail(email string) (User, error) {
	var u User

	row := DB.QueryRow("select id, email, name, created_at from users where email = $1", email)
	if row.Err() != nil {
		return u, row.Err()
	}

	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt)
	return u, err
}
