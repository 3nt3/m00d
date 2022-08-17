package main

import (
	"io/ioutil"
	"log"
	"m00d/db"
	"m00d/routes"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
)

func main() {
	// read config
	var config db.Config

	content, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Printf("error reading config.toml: %v\n", err)
		return
	}

	if _, err = toml.Decode(string(content), &config); err != nil {
		log.Printf("error decoding config.toml: %v\n", err)
		return
	}

	// initialize db
	if err := db.Init(config); err != nil {
		log.Panicf("error initiazing db: %v\n", err)
	}
	log.Printf("successfully initialized db!")

	// routes
	r := mux.NewRouter()

	r.HandleFunc("/moods", routes.NewMood).Methods("POST")
	r.HandleFunc("/moods", routes.GetMoods).Methods("GET")
	r.HandleFunc("/login", routes.Login).Methods("POST")
	r.HandleFunc("/refresh-token", routes.Refresh).Methods("POST")

	log.Printf("listening on :8080")
	http.ListenAndServe(":8080", r)

	db.DB.Close()
}
