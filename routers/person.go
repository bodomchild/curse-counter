package routers

import (
	"curse-count/db"
	"curse-count/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func Person(w http.ResponseWriter, r *http.Request) {
	var person models.Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		notFound(w)
		return
	}

	id, err := db.InsertPerson(person)
	if err != nil {
		notFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(id)
}

func Count(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	countType := r.URL.Query()["t"]
	qt := 0

	switch countType[0] {
	case "a":
		qt = 1
	case "s":
		qt = -1
	}

	person, ok := db.Count(id, qt)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(person)
}
