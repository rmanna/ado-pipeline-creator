package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type PageVars struct {
	SelectStatus   string
	AgentType      string
	Title          string
	ServiceName    string
	ServiceValue   string
	BuildType      string
	BuildTypeValue string
	Email          string
	EmailValue     string
	BuildFields    []Input
}

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
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

func render(w http.ResponseWriter, tmpl string, pageVars PageVars) {

	// prefix the name passed in with templates/
	tmpl = fmt.Sprintf("templates/%s", tmpl)
	//parse the template file held in the templates folder
	t, err := template.ParseFiles(tmpl)

	if err != nil { // if there is an error
		log.Print("template parsing error: ", err)
	}

	//execute the template and pass in the variables to fill the gaps
	err = t.Execute(w, pageVars)

	if err != nil {
		log.Print("template executing error: ", err)
	}
}
