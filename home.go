package main

import (
	"net/http"

	logger "github.com/rmanna/ado-pipeline-creator/internal/logger"
	"github.com/rmanna/ado-pipeline-creator/internal/pagetmpl"
)

//handler for / renders the home.html
func home(w http.ResponseWriter, req *http.Request) {
	pageVars := pagetmpl.PageVars{
		Title: "Pipeline Creator",
	}
	pagetmpl.Render(w, "home.html", pageVars)
	logger.Log.RequestFields(req.Method, req.URL.Path)
}
