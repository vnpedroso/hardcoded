package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func index(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		rw.Header().Set("Allow", http.MethodGet)
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rw.Header().Set("Content-Type", "text/html;charset=utf-8")
	io.WriteString(rw, `
	<h1>rpg characters API</h1>
	<p>Welcome! Please check our documentation</p>
	<a href="/docs">docs</a>
	`)
}

func docs(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		rw.Header().Set("Allow", http.MethodGet)
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rw.Header().Set("Content-Type", "text/html;charset=utf-8")
	io.WriteString(rw, `<p>to be written...</p>`)
}

func getAllChars(rw http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	chars := cdb.getAll()

	err := json.NewEncoder(&buf).Encode(chars)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(buf.Bytes())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

}

func getCharById(rw http.ResponseWriter, req *http.Request) {
	charId := req.PathValue("id")
	char, err := cdb.get(charId)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(char)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(buf.Bytes())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createChar(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var char characterJSON

	err := json.NewDecoder(req.Body).Decode(&char)
	if err != nil {
		http.Error(rw, "Invalid JSON payload containing character data", http.StatusBadRequest)
		return
	}

	charNameHex := hex.EncodeToString([]byte(strings.ToLower(char.Name)))
	charID := string(charNameHex)
	char.Id = charID

	if !cdb.charExists(charID) {

		var buf bytes.Buffer

		err = json.NewEncoder(&buf).Encode(char)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = cdb.create(charID, char)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
		_, err = rw.Write(buf.Bytes())
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	http.Error(rw, "character already exists!", http.StatusConflict) // status 409
}

func updateChar(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	charId := req.PathValue("id")

	if cdb.charExists(charId) {
		var char characterJSON
		err := json.NewDecoder(req.Body).Decode(&char)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = cdb.delete(charId)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = cdb.create(charId, char)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(char)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-type", "application/json")
		_, err = rw.Write(buf.Bytes())
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	http.Error(rw, "character with the provided id does not exist", http.StatusNotFound) //404

}

func deleteChar(rw http.ResponseWriter, req *http.Request) {
	charId := req.PathValue("id")

	if cdb.charExists(charId) {

		char := cdb.db[charId]
		var buf bytes.Buffer

		err := json.NewEncoder(&buf).Encode(char)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = cdb.delete(charId)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusNoContent)

		return
	}

	http.Error(rw, "character with the provided id does not exist", http.StatusNotFound) //404
}
