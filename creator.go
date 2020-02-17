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
		Status:         "notchanged",
		ServiceName:    "ServiceName",
		BuildType:      "BuildType",
		BuildTypeValue: "select",
		Email:          "Email",
	}
	render(w, "creator.html", pageVars)
}

// On submit
func execute(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	ServiceName := r.FormValue("ServiceName")
	BuildType := r.FormValue("BuildType")
	Email := r.FormValue("Email")
	Status := r.FormValue("status")

	if Status == "changed" {
		Inputs := []Input{
			Input{"Gradle Tasks", "test bootJar", "inputs"},
			Input{"Gradle Options", "-Xmx1024m", "inputs"},
			Input{"Java Home Option", "JDKVersion", "inputs"},
		}

		Status := "notchanged"
		pageVars := PageVars{
			Status:         Status,
			ServiceName:    "ServiceName",
			ServiceValue:   ServiceName,
			BuildType:      "BuildType",
			BuildTypeValue: BuildType,
			Email:          "Email",
			EmailValue:     Email,
			BuildFields:    Inputs,
		}
		render(w, "creator.html", pageVars)
	} else {
		setupTargetFiles()

		fmt.Println("Service Name is:", r.FormValue("ServiceName"))
		fmt.Println("Build Type is:", r.FormValue("BuildType"))
		//fmt.Println("Email is:", r.FormValue("Email"))

		updateFile("pipelineTemplates/buildDefinitionTemplateRequest.json", "target/buildDefinitionRequest.json", "SERVICENAME", ServiceName)
		updateFile("pipelineTemplates/sonar.properties", "target/sonar-project.properties", "SERVICENAME", ServiceName)

		//	buildType := r.FormValue("buildType")
		switch BuildType {
		case "javaGradle":
			fmt.Println("javaGradle")
			updateFile("pipelineTemplates/azure-gradle-pipeline.yaml", "target/azure-pipeline.yaml", "BUILDTYPE", BuildType)
		case "javaMvn":
			fmt.Println("javaMvn")
			updateFile("pipelineTemplates/azure-maven-pipeline.yaml", "target/azure-pipeline.yaml", "BUILDTYPE", BuildType)
		case "vueNpm":
			fmt.Println("vueNpm")
		case "angularNpm":
			fmt.Println("angularNpm")
		case "golang":
			fmt.Println("golang")
		}
	}
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
