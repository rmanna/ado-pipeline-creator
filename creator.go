package main

import (
	"fmt"
	"net/http"

	"github.com/rmanna/ado-pipeline-creator/internal/fileutils"
	logger "github.com/rmanna/ado-pipeline-creator/internal/logger"
	"github.com/rmanna/ado-pipeline-creator/internal/pagetmpl"
)

func creator(w http.ResponseWriter, req *http.Request) {

	pageVars := pagetmpl.PageVars{
		SelectStatus:   "false",
		ServiceName:    "ServiceName",
		BuildType:      "BuildType",
		BuildTypeValue: "select",
		Email:          "Email",
	}
	pagetmpl.Render(w, "creator.html", pageVars)
	logger.Log.RequestFields(req.Method, req.URL.Path)
}

// On submit
func execute(w http.ResponseWriter, r *http.Request) {

	config := fileutils.ReadYamlConfig("./internal", "config")
	r.ParseForm()

	var ServiceName = r.FormValue("ServiceName")
	var BuildType = r.FormValue("BuildType")
	var Email = r.FormValue("Email")
	var SelectStatus = r.FormValue("status")
	var BuildFields []pagetmpl.Input

	if SelectStatus == "true" {

		switch BuildType {
		case "gradle":
			logger.Log.InfoArg("Selected gradle build type")
			BuildFields = []pagetmpl.Input{
				pagetmpl.Input{Name: "Gradle Tasks", Value: config.BuildType.Gradle.Tasks, ID: "inputs"},
				pagetmpl.Input{Name: "Gradle Options", Value: config.BuildType.Gradle.Options, ID: "inputs"},
				pagetmpl.Input{Name: "Gradle Java Home", Value: config.BuildType.Gradle.JavaHomeOptions, ID: "inputs"},
			}
		case "maven":
			logger.Log.InfoArg("Selected maven build type")
			BuildFields = []pagetmpl.Input{
				pagetmpl.Input{Name: "Maven Options", Value: config.BuildType.Maven.Options, ID: "inputs"},
				pagetmpl.Input{Name: "Maven Goals", Value: config.BuildType.Maven.Goals, ID: "inputs"},
			}
		case "vue":
			logger.Log.InfoArg("Selected vue build type")
			BuildFields = []pagetmpl.Input{
				pagetmpl.Input{Name: "Vue Command", Value: config.BuildType.Vue.Command, ID: "inputs"},
			}
		case "angular":
			logger.Log.InfoArg("Selected angular build type")
			BuildFields = []pagetmpl.Input{
				pagetmpl.Input{Name: "Angular Command", Value: config.BuildType.Angular.Command, ID: "inputs"},
			}
		case "golang":
			logger.Log.InfoArg("Selected golang build type")
			BuildFields = []pagetmpl.Input{
				pagetmpl.Input{Name: "Golang Command", Value: config.BuildType.Golang.Command, ID: "inputs"},
			}
		}

		SelectStatus := "false"

		pageVars := pagetmpl.PageVars{
			SelectStatus:   SelectStatus,
			ServiceName:    "ServiceName",
			ServiceValue:   ServiceName,
			BuildType:      "BuildType",
			BuildTypeValue: BuildType,
			Email:          "Email",
			EmailValue:     Email,
			BuildFields:    BuildFields,
		}
		pagetmpl.Render(w, "creator.html", pageVars)

	} else {
		fmt.Println("Service Name is:", r.FormValue("ServiceName"))
		fmt.Println("Build Type is:", r.FormValue("BuildType"))
		fmt.Println("Email is:", r.FormValue("Email"))

		var buildTypeTemplate string
		switch BuildType {
		case "gradle":
			logger.Log.InfoArg("Selected gradle build type")
			buildTypeTemplate = "configs/azure-gradle-pipeline.yaml"
		case "maven":
			logger.Log.InfoArg("Selected maven build type")
			buildTypeTemplate = "configs/azure-maven-pipeline.yaml"
		case "vue":
			logger.Log.InfoArg("Selected vue build type")
			buildTypeTemplate = "configs/azure-vue-pipeline.yaml"
		case "angular":
			logger.Log.InfoArg("Selected angular build type")
			buildTypeTemplate = "configs/azure-angular-pipeline.yaml"
		case "golang":
			logger.Log.InfoArg("Selected golang build type")
			buildTypeTemplate = "configs/azure-golang-pipeline.yaml"
		}

		fileutils.SearchReplace(buildTypeTemplate, "azure-pipeline.yaml", "BUILDTYPE", BuildType)
		logger.Log.InfoArg("Successfully generated pipeline definition")

		fileutils.SearchReplace("configs/sonar.properties", "sonar-project.properties", "SERVICENAME", ServiceName)
		logger.Log.InfoArg("Successfully generated sonar properties")

		fileutils.SearchReplace("configs/buildDefinitionTemplateRequest.json", "buildDefinitionRequest.json", "SERVICENAME", ServiceName)
		logger.Log.InfoArg("Successfully create pipeline build in azure devops")

		// GITHUB SETUP
		//github.CreateRepository("ca4b5735341d33b0fd6fdf214244f0f5909c901d", ServiceName, true, true)
		//github.CommitBranch("ca4b5735341d33b0fd6fdf214244f0f5909c901d", ServiceName)
		//github.AddCollaborator("ca4b5735341d33b0fd6fdf214244f0f5909c901d", "ralphmanna", "admin")

		//adminTeams := "ol-devops"
		//writeTeams := "ol-offshore-devops,ol-development,ol-qa"
		//github.InviteTeams(adminTeams, "admin")
		//github.InviteTeams(writeTeams, "write")

		// ADO SETUP
		// SEND REQUEST TO ADO USING buildDefinitionRequest.json

		// SLACK SETUP
		// CREATE OR ADD TO SLACK CHANNEL

		http.Redirect(w, r, "/results", http.StatusFound)
	}
}
