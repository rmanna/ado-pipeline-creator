package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Input struct {
	Name  string
	Value string
	Id    string
}

func creator(w http.ResponseWriter, r *http.Request) {

	pageVars := PageVars{
		ServiceName: "serviceName",
		BuildType:   "buildType",
		Email:       "email",
	}
	render(w, "creator.html", pageVars)
}

func userSelected(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	//buildType := r.FormValue("buildType")

	Inputs := []Input{
		Input{"gradleTasks", "test bootJar", "inputs"},
		Input{"gradleOptions", "-Xmx1024m", "inputs"},
		Input{"javaHomeOption", "JDKVersion", "inputs"},
	}
	pageVars := PageVars{
		BuildFields: Inputs,
	}
	render(w, "creator.html", pageVars)
}

// On submit
func execute(w http.ResponseWriter, r *http.Request) {

	setupTargetFiles()
	r.ParseForm()

	fmt.Println("Service Name is:", r.FormValue("serviceName"))
	fmt.Println("Build Type is:", r.FormValue("buildType"))
	fmt.Println("Email is:", r.FormValue("email"))

	serviceName := r.FormValue("serviceName")
	updateFile("pipelineTemplates/buildDefinitionTemplateRequest.json", "target/buildDefinitionRequest.json", "SERVICENAME", serviceName)
	updateFile("pipelineTemplates/sonar.properties", "target/sonar-project.properties", "SERVICENAME", serviceName)

	buildType := r.FormValue("buildType")
	switch buildType {
	case "javaGradle":
		fmt.Println("javaGradle")
		updateFile("pipelineTemplates/azure-gradle-pipeline.yaml", "target/azure-pipeline.yaml", "BUILDTYPE", buildType)
	case "javaMvn":
		fmt.Println("javaMvn")
		updateFile("pipelineTemplates/azure-maven-pipeline.yaml", "target/azure-pipeline.yaml", "BUILDTYPE", buildType)
	case "vueNpm":
		fmt.Println("vueNpm")
	case "angularNpm":
		fmt.Println("angularNpm")
	case "golang":
		fmt.Println("golang")
	}

	//http.Redirect(w, r, "/creator", http.StatusFound)
}

// Create Dir or Cleanup Dir
func setupTargetFiles() {
	if _, err := os.Stat("./target"); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("target", 0755)
		} else {
			os.Remove("target/*")
		}
	}
}

// Find and Replace
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
