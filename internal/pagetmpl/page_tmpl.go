package pagetmpl

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Input struct
type Input struct {
	Name  string
	Value string
	ID    string
}

// PageVars Struct
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

// Render website template page
func Render(w http.ResponseWriter, tmpl string, pageVars PageVars) {

	// prefix the name passed in with templates/
	tmpl = fmt.Sprintf("web/template/%s", tmpl)
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
