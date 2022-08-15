package main

import (
	"m00d/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/moods", routes.NewMood).Methods("POST")
	r.HandleFunc("/moods", routes.GetMoods).Methods("GET")

	r.HandleFunc("/login", routes.Login).Methods("POST")
}
