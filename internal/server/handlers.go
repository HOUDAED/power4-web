package web

import (
	"html/template"
	"log"
	"net/http"
)

type Game struct {
    Player1       string
    Player2       string
    Difficulty    string
    Rows          int
    Cols          int
    Grid          [][]int
    CurrentPlayer int
    Winner        int
    Draw          bool
    ColRange      []int  // ajout pour générer les boutons colonnes
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
        return
	}

    player1 := r.FormValue("player1")
	player2 := r.FormValue("player2")
	difficulty := r.FormValue("difficulty")
    
	if player1 == "" || player2 == "" || difficulty == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

func NewGame(player1, player2, difficulty string) *Game {
    var rows, cols int
    switch difficulty {
    case "easy":
        rows, cols = 6, 7
    case "normal":
        rows, cols = 6, 9
    case "hard":
        rows, cols = 7, 8
    }

    grid := make([][]int, rows)
    for i := range grid {
        grid[i] = make([]int, cols)
    }

    colRange := make([]int, cols)
    for i := 0; i < cols; i++ {
        colRange[i] = i
    }

    return &Game{
        Player1:       player1,
        Player2:       player2,
        Difficulty:    difficulty,
        Rows:          rows,
        Cols:          cols,
        Grid:          grid,
        CurrentPlayer: 1,
        ColRange:      colRange,
    }
}

func GameHandler(w http.ResponseWriter, r *http.Request) {
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
    game := NewGame(player1, player2, difficulty)

    tmpl, err := template.ParseFiles("templates/game.html")
    if err != nil {
        log.Fatal(err)
    }
    tmpl.Execute(w, game)
}
