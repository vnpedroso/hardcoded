package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var cdb *characterDb

func init() {
	cdb = &characterDb{
		db: make(map[string]characterJSON),
	}
}

func WriteErrorJSON(rw http.ResponseWriter, error_msg string, status_code int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status_code)

	err := json.NewEncoder(rw).Encode(jsonError{
		Code:    status_code,
		Message: error_msg,
	})

	if err != nil {
		http.Error(rw, error_msg, status_code)
	}
}

func CharacterMultiplex(rw http.ResponseWriter, req *http.Request) {
	//multiplex handlerFunc for the /characters endpoint
	switch req.Method {
	case http.MethodGet:
		getAllChars(rw, req)
	case http.MethodPost:
		createChar(rw, req)
	default:
		WriteErrorJSON(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func CharacterIdMultiplex(rw http.ResponseWriter, req *http.Request) {
	//multiplexer handlerFunc for the /characters/{id} endpoint
	switch req.Method {
	case http.MethodGet:
		getCharById(rw, req)
	case http.MethodDelete:
		deleteChar(rw, req)
	case http.MethodPut:
		updateChar(rw, req)
	default:
		WriteErrorJSON(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.HandleFunc("/", index)
	http.HandleFunc("/docs", docs)
	http.HandleFunc("/characters", CharacterMultiplex)
	http.HandleFunc("/characters/{id}", CharacterIdMultiplex)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
