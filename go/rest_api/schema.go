package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

type characterJSON struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name"`
	Class      string `json:"class"`
	Race       string `json:"race"`
	Level      int    `json:"level"`
	MainWeapon string `json:"main_weapon,omitempty"`
}

type characterPayload []characterJSON

type jsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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

type characterDb struct {
	mtx sync.Mutex
	db  map[string]characterJSON
}

func (cdb *characterDb) getAll() characterPayload {
	cdb.mtx.Lock()
	defer cdb.mtx.Unlock()
	out := make(characterPayload, 0, len(cdb.db))
	for _, v := range cdb.db {
		out = append(out, v)
	}
	return out
}

func (cdb *characterDb) charExists(hexId string) bool {
	cdb.mtx.Lock()
	defer cdb.mtx.Unlock()
	_, ok := cdb.db[hexId]
	return ok
}

func (cdb *characterDb) create(hexId string, char characterJSON) error {
	ok := cdb.charExists(hexId)

	if !ok {
		cdb.mtx.Lock()
		cdb.db[hexId] = char
		cdb.mtx.Unlock()
		return nil
	}

	return errors.New("CharacterAlreadyExists: A character with the provided id already exists")
}

func (cdb *characterDb) delete(hexId string) error {
	if cdb.charExists(hexId) {
		cdb.mtx.Lock()
		delete(cdb.db, hexId)
		cdb.mtx.Unlock()
		return nil
	}

	return errors.New("CharacterDoesNotExist: A character with the provided id does not exist")
}

func (cdb *characterDb) get(hexId string) (characterJSON, error) {
	var char characterJSON
	cdb.mtx.Lock()
	char, ok := cdb.db[hexId]
	cdb.mtx.Unlock()
	if ok {
		return char, nil
	}
	return char, errors.New("CharacterDoesNotExist: A character with the provided id does not exist")
}
