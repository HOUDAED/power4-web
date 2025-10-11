package main

import (
	"net/http"
	"power4/internal/server"
)


func main() {
	http.HandleFunc("/", web.IndexHandler)
    http.HandleFunc("/welcome", web.WelcomeHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe(":8080", nil)
}