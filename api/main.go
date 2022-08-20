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
	r.Use(corsMiddleware)

	log.Printf("listening on :8080")
	err = http.ListenAndServe("0.0.0.0:8080", r)
	if err != nil {
		log.Printf("error listening: %v\n", err)
	}

	db.DB.Close()
}

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			handlePreflight(w, r)
		}

		// do normal stuff
		h.ServeHTTP(w, r)
	})
}

func handlePreflight(w http.ResponseWriter, r *http.Request) {
	// set cors headers
	// very secure lol
	log.Printf("%+v\n", r.Method)
	w.Header().Add("access-control-allow-origin", "*")
	w.Header().Add("access-control-allow-credentials", "true")
	w.Header().Add("access-control-allow-headers", "authorization")
	w.Header().Add("access-control-allow-methods", "get,put,post,delete,options")
	w.Header().Add("access-control-expose-headers", "origin,authorization,Authorization,*")
}
