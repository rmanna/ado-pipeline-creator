package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func creator(w http.ResponseWriter, req *http.Request) {

	pageVars := PageVars{
		ServiceName: "serviceName",
		AgentType:   "agentType",
		Email:       "email",
	}
	render(w, "creator.html", pageVars)
}

func execute(w http.ResponseWriter, req *http.Request) {

	setupTargetFiles()
	req.ParseForm()

	fmt.Println("Service Name is:", req.FormValue("serviceName"))
	fmt.Println("Agent Type is:", req.FormValue("agentType"))
	fmt.Println("Email is:", req.FormValue("email"))

	serviceName := req.FormValue("serviceName")
	updateFile("pipelineTemplates/buildDefinitionTemplateRequest.json", "target/buildDefinitionRequest.json", "SERVICENAME", serviceName)
	updateFile("pipelineTemplates/sonar.properties", "target/sonar-project.properties", "SERVICENAME", serviceName)

	agentType := req.FormValue("agentType")
	switch agentType {
	case "javaGradle":
		fmt.Println("javaGradle")
		updateFile("pipelineTemplates/azure-gradle-pipeline.yaml", "target/azure-pipeline.yaml", "AGENTTYPE", agentType)
	case "javaMvn":
		fmt.Println("javaMvn")
		updateFile("pipelineTemplates/azure-maven-pipeline.yaml", "target/azure-pipeline.yaml", "AGENTTYPE", agentType)
	case "vueNpm":
		fmt.Println("vueNpm")
	case "angularNpm":
		fmt.Println("angularNpm")
	case "golang":
		fmt.Println("golang")
	}

	http.Redirect(w, req, "/creator", http.StatusFound)
}

func setupTargetFiles() {
	if _, err := os.Stat("./target"); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("target", 0755)
		} else {
			os.Remove("target/*")
		}
	}
}

func updateFile(sourceFileName string, targetFileName string, sourceString string, targetString string) {
	input, err := ioutil.ReadFile(sourceFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output := strings.Replace(string(input), sourceString, targetString, -1)

	if err = ioutil.WriteFile(targetFileName, []byte(output), 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
