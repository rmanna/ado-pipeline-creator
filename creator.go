package main

import (
	"fmt"
	"net/http"

	"github.com/rmanna/ado-pipeline-creator/internal/fileutils"
	"github.com/rmanna/ado-pipeline-creator/internal/github"
	"github.com/rmanna/ado-pipeline-creator/internal/pagetmpl"
)

func creator(w http.ResponseWriter, r *http.Request) {

	pageVars := pagetmpl.PageVars{
		SelectStatus:   "false",
		ServiceName:    "ServiceName",
		BuildType:      "BuildType",
		BuildTypeValue: "select",
		Email:          "Email",
	}
	pagetmpl.Render(w, "creator.html", pageVars)
}

// On submit
func execute(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	var ServiceName = r.FormValue("ServiceName")
	var BuildType = r.FormValue("BuildType")
	var Email = r.FormValue("Email")
	var SelectStatus = r.FormValue("status")
	var BuildFields []pagetmpl.Input

	if SelectStatus == "true" {

		switch BuildType {
		case "gradle":
			BuildFields = []pagetmpl.Input{
				pagetmpl.Input{Name: "Gradle Tasks", Value: "test bootJar", ID: "inputs"},
				pagetmpl.Input{Name: "Gradle Options", Value: "-Xmx1024m", ID: "inputs"},
				pagetmpl.Input{Name: "Java Home Option", Value: "JDKVersion", ID: "inputs"},
			}
		case "maven":
			BuildFields = []pagetmpl.Input{
				pagetmpl.Input{Name: "Maven Options", Value: "-Xmx3072m", ID: "inputs"},
				pagetmpl.Input{Name: "Maven Goals", Value: "clean package", ID: "inputs"},
			}
		case "vue":
			BuildFields = []pagetmpl.Input{
				pagetmpl.Input{Name: "Vue Command", Value: "", ID: "inputs"},
			}
		case "angular":
			BuildFields = []pagetmpl.Input{
				pagetmpl.Input{Name: "Npm Command", Value: "run build", ID: "inputs"},
			}
		case "golang":
			BuildFields = []pagetmpl.Input{
				pagetmpl.Input{Name: "GO Command", Value: "", ID: "inputs"},
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
		//fmt.Println("Email is:", r.FormValue("Email"))

		//	buildType := r.FormValue("buildType")
		switch BuildType {
		case "gradle":
			fmt.Println("gradle")
			fileutils.SearchReplace("configs/azure-gradle-pipeline.yaml", "azure-pipeline.yaml", "BUILDTYPE", BuildType)
		case "maven":
			fmt.Println("maven")
			fileutils.SearchReplace("configs/azure-maven-pipeline.yaml", "azure-pipeline.yaml", "BUILDTYPE", BuildType)
		case "vue":
			fmt.Println("vue")
		case "angular":
			fmt.Println("angular")
		case "golang":
			fmt.Println("golang")
		}

		fileutils.SearchReplace("configs/sonar.properties", "sonar-project.properties", "SERVICENAME", ServiceName)
		fileutils.SearchReplace("configs/buildDefinitionTemplateRequest.json", "buildDefinitionRequest.json", "SERVICENAME", ServiceName)

		// GITHUB SETUP
		github.CreateRepository(ServiceName)
		// github.CommitBranch First Commit
		// github.InviteTeams Team invites

		// ADO SETUP
		// SEND REQUEST TO ADO USING buildDefinitionRequest.json

		// SLACK SETUP
		// CREATE OR ADD TO SLACK CHANNEL

		http.Redirect(w, r, "/results", http.StatusFound)
	}
}
