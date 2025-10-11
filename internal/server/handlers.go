package web

import (
	"html/template"
	"log"
	"net/http"
)

type Game struct {
	Player1    string
	Player2    string
	Difficulty string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	game := Game{
		Player1:    r.FormValue("player1"),
		Player2:    r.FormValue("player2"),
		Difficulty: r.FormValue("difficulty"),
	}

	if game.Player1 == "" || game.Player2 == "" || game.Difficulty == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	tmpl, err := template.ParseFiles("templates/welcome.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, game)
}
