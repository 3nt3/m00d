package routes

import (
	"encoding/json"
	"log"
	"m00d/db"
	"net/http"
)

func NewMood(w http.ResponseWriter, r *http.Request) {
	u, status, err := authorizeRequest(r)
	if err != nil {
		log.Printf("error authorizing? (status %d): %v\n", status, err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var m db.Mood
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m.UserID = u.ID

	m, err = db.NewMood(m)
	if err != nil {
		log.Printf("error writing to db: %v\n", err)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("%+v\n", m)

	json.NewEncoder(w).Encode(m)
}

func GetMoods(w http.ResponseWriter, r *http.Request) {
	u, status, err := authorizeRequest(r)
	if err != nil {
		log.Printf("error authorizing? (status %d): %v\n", status, err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
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
