package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"m00d/db"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	var idToken string
	var v map[string]string

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("%+v\n", v)

	idToken, ok := v["id_token"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	googleClaims, err := validateGoogleJWT(idToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		log.Printf("error: %v\n", err)
		return
	}

	emailExists, err := db.CheckEmailExists(googleClaims.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error reading db: %v\n", err)
		return
	}

	var u db.User
	if !emailExists {
		log.Printf("user doesn't exist yet.")
		u = db.User{Email: googleClaims.Email, Name: googleClaims.FirstName + " " + googleClaims.LastName}
		u.ID, err = db.CreateUser(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error writing to db: %v\n", err)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		log.Printf("created new user: %+v\n", u)
	} else {
		log.Printf("user exists.")
		u, err = db.GetUserByEmail(googleClaims.Email)
		if err != nil {
			log.Printf("error reading db: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		log.Printf("user logged in: %+v\n", u)
	}

	signedToken, err := db.JwtWrapper.GenerateToken(u.Email)
	if err != nil {
		log.Printf("error signing token: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(map[string]string{"token": signedToken})
}

// this is stolen
func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}

// I stole all of this
// https://medium.com/bootdotdev/how-to-implement-sign-in-with-google-in-golang-962052ac5b95
func validateGoogleJWT(tokenString string) (GoogleClaims, error) {
	claimsStruct := GoogleClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
			if err != nil {
				return nil, err
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			if err != nil {
				return nil, err
			}
			return key, nil
		},
	)
	if err != nil {
		return GoogleClaims{}, err
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok {
		return GoogleClaims{}, errors.New("Invalid Google JWT")
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return GoogleClaims{}, errors.New("iss is invalid")
	}

	log.Printf("%v\n", claims.Audience)
	if claims.Audience != "82145806916-vocueu5na49d2lgusnotbrjdd7ne77mp.apps.googleusercontent.com" {
		return GoogleClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return GoogleClaims{}, errors.New("JWT is expired")
	}

	return *claims, nil
}

func authorizeRequest(r *http.Request) (db.User, int, error) {
	var u db.User

	header := r.Header.Get("authorization")
	if header == "" {
		return u, http.StatusUnauthorized, errors.New("request not authorized.")
	}

	split := strings.Split(header, " ")
	if len(split) != 2 {
		return u, http.StatusUnauthorized, errors.New(fmt.Sprintf("authorization header malformed: %s", header))
	}

	claims, err := db.JwtWrapper.ValidateToken(split[1])
	if err != nil {
		return u, http.StatusInternalServerError, errors.New(fmt.Sprintf("error validating token: %v", err))
	}

	u, err = db.GetUserByEmail(claims.Email)
	if err != nil {
		return u, http.StatusInternalServerError, errors.New(fmt.Sprintf("error reading from db: %v", err))
	}

	return u, http.StatusOK, nil
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	u, status, err := authorizeRequest(r)
	if err != nil {
		log.Printf("error authorizing? (status %d): %v\n", status, err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	signedToken, err := db.JwtWrapper.GenerateToken(u.Email)
	if err != nil {
		log.Printf("error signing token: %v\n", err)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": signedToken})
}
