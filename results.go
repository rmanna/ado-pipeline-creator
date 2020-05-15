package main

import (
	"net/http"
	"github.com/rmanna/ado-pipeline-creator/internal/pagetmpl"
)

//handler for / renders the home.html
func results(w http.ResponseWriter, req *http.Request) {
	pageVars := pagetmpl.PageVars{
		Title: "Pipeline Creator",
	}
	pagetmpl.Render(w, "results.html", pageVars)
}
