package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Event struct {
	Key   string
	Value string
}

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

		if err := saveEventToLogFile(event, "myLog.bin"); err != nil {
			log.Print(err)
			return
		}

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// parse event from request body of the producer
func parseRequest(body io.Reader, event *Event) error {
	_, err := fmt.Fscanf(body, "Key=%s\nValue=%s", &event.Key, &event.Value)
	log.Print(event)
	return err
}

func saveEventToLogFile(event Event, filename string) error {
	// convert key value pair into binary
	keyByte := []byte(event.Key)
	valueByte := []byte(event.Value)
	// length of key and value bytes
	keyByteLen := uint16(len(keyByte))
	valueByteLen := uint16(len(valueByte))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0222)
	if err != nil {
		return err
	}
	defer file.Close()

	// write key value and len of both to the file
	if err := binary.Write(file, binary.LittleEndian, keyByteLen); err != nil {
		return err
	}
	if _, err := file.Write(keyByte); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, valueByteLen); err != nil {
		return err
	}
	if _, err := file.Write(valueByte); err != nil {
		return err
	}

	return nil
}

func main() {
	http.HandleFunc("/", Handler)

	// start the http server
	log.Println("Localhost running at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
