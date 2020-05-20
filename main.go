package main

import (
	"net/http"

	logger "github.com/rmanna/ado-pipeline-creator/internal/logger"
	"github.com/rmanna/ado-pipeline-creator/internal/network"
)

func main() {

	logger.NewLogger()

	http.Handle("/web/css/", http.StripPrefix("/web/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/", home)
	http.HandleFunc("/creator", creator)
	http.HandleFunc("/execute", execute)
	http.HandleFunc("/results", results)
	http.ListenAndServe(network.SetPort(), nil)
}
