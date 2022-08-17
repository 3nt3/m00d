package routes

import (
	"encoding/json"
	"log"
	"m00d/db"
	"net/http"
	"strings"
)

func NewMood(w http.ResponseWriter, r *http.Request) {
}

func GetMoods(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	if header == "" {
		log.Printf("request not authorized.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	split := strings.Split(header, " ")
	if len(split) != 2 {
		log.Printf("authorization header malformed: %s\n", header)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, err := db.JwtWrapper.ValidateToken(split[1])
	if err != nil {
		log.Printf("error validating token: %v\n", err)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	u, err := db.GetUserByEmail(claims.Email)
	if err != nil {
		log.Printf("error reading from db: %v\n", err)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("user: %+v\n", u)

	moods, err := db.GetMoods(u.ID)
	if err != nil {
		log.Printf("error reading from db: %v\n", err)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(moods)
}
