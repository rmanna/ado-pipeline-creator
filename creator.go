package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// Input struct
type Input struct {
	Name  string
	Value string
	ID    string
}

func creator(w http.ResponseWriter, r *http.Request) {

	pageVars := PageVars{
		SelectStatus:   "false",
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

	var ServiceName = r.FormValue("ServiceName")
	var BuildType = r.FormValue("BuildType")
	var Email = r.FormValue("Email")
	var SelectStatus = r.FormValue("status")
	var BuildFields []Input

	if SelectStatus == "true" {

		switch BuildType {
		case "gradle":
			BuildFields = []Input{
				Input{"Gradle Tasks", "test bootJar", "inputs"},
				Input{"Gradle Options", "-Xmx1024m", "inputs"},
				Input{"Java Home Option", "JDKVersion", "inputs"},
			}
		case "maven":
			BuildFields = []Input{
				Input{"Maven Options", "-Xmx3072m", "inputs"},
				Input{"Maven Goals", "clean package", "inputs"},
			}
		case "vue":
			BuildFields = []Input{
				Input{"Vue Command", "", "inputs"},
			}
		case "angular":
			BuildFields = []Input{
				Input{"Npm Command", "run build", "inputs"},
			}
		case "golang":
			BuildFields = []Input{
				Input{"GO Command", "", "inputs"},
			}
		}

		SelectStatus := "false"

		pageVars := PageVars{
			SelectStatus:   SelectStatus,
			ServiceName:    "ServiceName",
			ServiceValue:   ServiceName,
			BuildType:      "BuildType",
			BuildTypeValue: BuildType,
			Email:          "Email",
			EmailValue:     Email,
			BuildFields:    BuildFields,
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
		case "gradle":
			fmt.Println("gradle")
			updateFile("pipelineTemplates/azure-gradle-pipeline.yaml", "target/azure-pipeline.yaml", "BUILDTYPE", BuildType)
		case "maven":
			fmt.Println("maven")
			updateFile("pipelineTemplates/azure-maven-pipeline.yaml", "target/azure-pipeline.yaml", "BUILDTYPE", BuildType)
		case "vue":
			fmt.Println("vue")
		case "angular":
			fmt.Println("angular")
		case "golang":
			fmt.Println("golang")
		}

		createGithubRepo(ServiceName, "new repository")
		cloneGithubRepo(ServiceName, "rmanna")
		// cp target/* (git clone ServiceName)
		// cd (git clone ServiceName)
		// git add *
		// git commit -m"Initial Repository

		http.Redirect(w, r, "/results", http.StatusFound)
	}
}

func cloneGithubRepo(ServiceName string, OrgName string) {
	directory := "/tmp"
	url := "https://github.com/rmanna/" + ServiceName
	//token := os.Getenv("GITHUB_AUTH_TOKEN")
	token := "348a55efe1bd06a335b94f8c5788eabc81f4b558"

	r, err := git.PlainClone(directory, true, &git.CloneOptions{
		Auth: &githttp.BasicAuth{
			Username: "abc123", // yes, this can be anything except an empty string
			Password: token,
		},
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatal("Error returned from cloning repository:", err)
	}

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		log.Fatal("Error returned from cloning repository:", err)
	}

	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Fatal("Error returned from cloning repository:", err)
	}

	fmt.Println(commit)
}

func createGithubRepo(ServiceName string, Description string) {
	var (
		name        = flag.String("name", ServiceName, "Name of repo to create in authenticated user's GitHub account.")
		description = flag.String("description", Description, "Description of created repo.")
		private     = flag.Bool("private", true, "Will created repo be private.")
	)

	org := "rmanna"
	flag.Parse()
	//token := os.Getenv("GITHUB_AUTH_TOKEN")
	token := "348a55efe1bd06a335b94f8c5788eabc81f4b558"
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	if *name == "" {
		log.Fatal("No name: New repos must be given a name")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	r := &github.Repository{Name: name, Private: private, Description: description}
	repo, _, err := client.Repositories.Create(ctx, "", r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully created new repo: %v\n", repo.GetName())

	// Commit README.md as part of repository setup
	fileContent := []byte("This is a starter for the README")
	opts := &github.RepositoryContentFileOptions{
		Message:   github.String("Initial commit"),
		Content:   fileContent,
		Branch:    github.String("master"),
		Committer: &github.CommitAuthor{Name: github.String("Self-Serve Automation"), Email: github.String("oldevops@openlane.com")},
	}
	_, _, err = client.Repositories.CreateFile(ctx, org, *name, "README.md", opts)
	if err != nil {
		fmt.Println(err)
		return
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
