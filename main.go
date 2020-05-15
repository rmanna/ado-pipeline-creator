package main

import (
	"net/http"
	"os"
)

func main() {
	http.Handle("/web/css/", http.StripPrefix("/web/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/", home)
	http.HandleFunc("/creator", creator)
	http.HandleFunc("/execute", execute)
	http.HandleFunc("/results", results)
	http.ListenAndServe(getPort(), nil)
}

// Detect $PORT and if present uses it for listen and serve else defaults to :8080
func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}
