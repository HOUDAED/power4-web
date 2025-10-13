package web

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var currentGame *Game

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
		return
	}

	player1 := r.FormValue("player1")
	player2 := r.FormValue("player2")
	difficulty := r.FormValue("difficulty")

	if player1 == "" || player2 == "" || difficulty == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	game := Game{
		Player1:    player1,
		Player2:    player2,
		Difficulty: difficulty,
	}

	tmpl, err := template.ParseFiles("templates/welcome.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, game)
}


// Démarrage d'une nouvelle partie
// -------------------------
func GameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	player1 := r.FormValue("player1")
	player2 := r.FormValue("player2")
	difficulty := r.FormValue("difficulty")

	currentGame = NewGame(player1, player2, difficulty)

	tmpl, err := template.ParseFiles("templates/game.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, currentGame)
}

func PlayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	colStr := r.FormValue("col")
	colIndex, err := strconv.Atoi(colStr)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if currentGame == nil {
		player1 := r.FormValue("player1")
		player2 := r.FormValue("player2")
		difficulty := r.FormValue("difficulty")
		currentGame = NewGame(player1, player2, difficulty)
	}

	// Si la partie est terminée → juste réafficher le plateau
	if currentGame.Winner != 0 || currentGame.Draw {
		tmpl, _ := template.ParseFiles("templates/game.html")
		tmpl.Execute(w, currentGame)
		return
	}

	// Si la colonne est pleine, on ne redirige plus : on réaffiche le plateau avec un message optionnel
	if !currentGame.PlaceToken(colIndex) {
		// Optionnel : message d'erreur (colonne pleine)
		tmpl, _ := template.ParseFiles("templates/game.html")
		tmpl.Execute(w, currentGame)
		return
	}

	// Vérification victoire ou match nul
	winner := currentGame.CheckWin()
	if winner != 0 {
		currentGame.Winner = winner
	} else if currentGame.CheckDraw() {
		currentGame.Draw = true
	} else {
		if currentGame.CurrentPlayer == 1 {
			currentGame.CurrentPlayer = 2
		} else {
			currentGame.CurrentPlayer = 1
		}
	}

	// Réaffichage du plateau
	tmpl, err := template.ParseFiles("templates/game.html")
	if err != nil {
		log.Println("Erreur template:", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, currentGame)
}
