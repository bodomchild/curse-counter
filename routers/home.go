package routers

import (
	"bufio"
	"curse-count/db"
	"curse-count/models"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var wsURL = wsUrl()

func Home(w http.ResponseWriter, _ *http.Request) {
	t, err := template.ParseFiles("templates/home.gohtml")
	if err != nil {
		notFound(w)
		return
	}

	people, err := db.GetAll()

	data := struct {
		Personas     []*models.Person
		URLWebSocket string
	}{
		Personas:     people,
		URLWebSocket: wsURL,
	}

	if err != nil {
		notFound(w)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		notFound(w)
		return
	}
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	_ = json.NewEncoder(w).Encode("PAGE NOT FOUND")
}

func wsUrl() string {
	var url string
	f, err := os.Open("./props.txt")
	if err != nil {
		log.Fatalf("home: %v", err)
	}
	input := bufio.NewScanner(f)
	for input.Scan() {
		lines := strings.SplitN(input.Text(), "=", 2)
		if lines[0] == "ws" {
			url = strings.Trim(lines[1], " \"")
		}
	}
	_ = f.Close()
	return url
}
