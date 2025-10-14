package web

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Variable globale pour la partie en cours
var currentGame *Game

// Structure intermédiaire pour le rendu HTML
type GameView struct {
	Player1       string
	Player2       string
	Difficulty    string
	CurrentPlayer int
	Winner        int
	Draw          bool
	ColRange      []int
	GravityDown   bool
	Grid          [][]Cell // grille enrichie avec IsWinning
}

// Convertit une Game en GameView pour le template
func ToGameView(g *Game) GameView {
	return GameView{
		Player1:       g.Player1,
		Player2:       g.Player2,
		Difficulty:    g.Difficulty,
		CurrentPlayer: g.CurrentPlayer,
		Winner:        g.Winner,
		Draw:          g.Draw,
		ColRange:      g.ColRange,
		GravityDown:   g.GravityDown,
		Grid:          g.RenderedGrid(),
	}
}

// Page d'accueil
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

// Page de bienvenue après saisie des joueurs
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
	tmpl.Execute(w, ToGameView(currentGame))
}

// Jouer un coup
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
		tmpl.Execute(w, ToGameView(currentGame))
		return
	}

	// Placer le jeton en tenant compte de la gravité
	if !currentGame.PlaceToken(colIndex) {
		// Colonne pleine → réafficher le plateau
		tmpl, _ := template.ParseFiles("templates/game.html")
		tmpl.Execute(w, ToGameView(currentGame))
		return
	}

	// Vérification victoire ou match nul
	winner := currentGame.CheckWin()
	if winner != 0 {
		currentGame.Winner = winner
	} else if currentGame.CheckDraw() {
		currentGame.Draw = true
	} else {
		// Changement de joueur
		if currentGame.CurrentPlayer == 1 {
			currentGame.CurrentPlayer = 2
		} else {
			currentGame.CurrentPlayer = 1
		}
	}

	tmpl, err := template.ParseFiles("templates/game.html")
	if err != nil {
		log.Println("Erreur template:", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, ToGameView(currentGame))
}

// Revanche ou nouvelle partie
func RematchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	action := r.FormValue("type") // "revanche" ou "new"

	switch action {
	case "revanche":
		if currentGame != nil {
			currentGame.Reset()
		}
		tmpl, err := template.ParseFiles("templates/game.html")
		if err != nil {
			log.Println("Erreur template:", err)
			http.Error(w, "Erreur interne", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, ToGameView(currentGame))
		return
	case "new":
		currentGame = nil
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}