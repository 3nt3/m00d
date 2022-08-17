package db

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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

// JWTWrapper wraps the signing key and the issuer
type JWTWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

// JWTClaim adds email as a claim to the token
type JWTClaim struct {
	Email string
	jwt.StandardClaims
}

var JwtWrapper JWTWrapper

// GenerateToken generates a jwt token
func (j *JWTWrapper) GenerateToken(email string) (signedToken string, err error) {
	claims := &JWTClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return
	}

	return
}

//ValidateToken validates the jwt token
func (j *JWTWrapper) ValidateToken(signedToken string) (claims *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("Couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}

	return
}
