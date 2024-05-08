package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

type Event struct {
	Key   string
	Value string
}

var buf bytes.Buffer

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		// Respond with allowed methods information
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Methods allowed: POST\n")
		return
	} else if r.Method == http.MethodPost {
		var event Event

		if err := parseRequest(r.Body, &event); err != nil {
			log.Print(err)
			http.Error(w, "Error parsing request", http.StatusBadRequest)
			return
		}

		if err := saveEventToBuffer(event); err != nil {
			log.Print(err)
			http.Error(w, "Error saving event to buffer", http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	http.HandleFunc("/", Handler)

	// start the http server
	log.Println("Localhost running at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
