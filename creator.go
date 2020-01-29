package main

import (
	"net/http"
)

func Creator(w http.ResponseWriter, req *http.Request) {

	pageVars := PageVars{
		Id:            "12345",
		Title:         "Pipeline Creator",
		Email:         "email",
                AgentType:     "agentType",
	}
	render(w, "creator.html", pageVars)
}
