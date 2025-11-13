package main

import (
	"log"
	"net/http"
	"time"
)

var cdb *characterDb

func init() {
	cdb = &characterDb{
		db: make(map[string]characterJSON),
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()

		next.ServeHTTP(rw, req)

		log.Printf("%s %s %s [%v]",
			req.Method,
			req.URL.Path,
			req.RemoteAddr,
			time.Since(start),
		)
	})
}

func main() {

	log.SetPrefix("myGoServer: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	mux := http.NewServeMux()

	mux.Handle("/favicon.ico", http.NotFoundHandler())

	mux.HandleFunc("/", index)
	mux.HandleFunc("/docs", docs)
	mux.HandleFunc("/characters", CharacterMultiplex)
	mux.HandleFunc("/characters/{id}", CharacterIdMultiplex)

	server := loggingMiddleware(mux)

	err := http.ListenAndServe(":8080", server)
	if err != nil {
		log.Fatalln(err)
	}
}
