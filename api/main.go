package main

import (
	"log"
	"m00d/db"
	"m00d/routes"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// initialize db
	if err := db.Init(); err != nil {
		log.Panicf("error initiazing db: %v\n", err)
	}
	log.Printf("successfully initialized db!")

	// routes
	r := mux.NewRouter()

	r.HandleFunc("/moods", routes.NewMood).Methods("POST")
	r.HandleFunc("/moods", routes.GetMoods).Methods("GET")
	r.HandleFunc("/login", routes.Login).Methods("POST")

	log.Printf("listening on :8080")
	http.ListenAndServe(":8080", r)

	db.DB.Close()
}
