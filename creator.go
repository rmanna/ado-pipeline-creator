package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	sourceOwner   = flag.String("source-owner", "rmanna", "Name of the owner (user or org) of the repo to create the commit in.")
	baseBranch    = flag.String("base-branch", "master", "Name of branch to create the `commit-branch` from.")
	sourceFiles   = flag.String("files", "README.md,azure-pipeline.yaml,sonar-project.properties", "Comma-separated list of files to commit and their location.")
	authorName    = flag.String("author-name", "Self-Service Automation", "Name of the author of the commit.")
	authorEmail   = flag.String("author-email", "oldevops@openlane.com", "Email of the author of the commit.")
	commitMessage = flag.String("commit-message", "Initial Project Commit", "Content of the commit message.")
	repoPrivacy   = flag.Bool("private", true, "Whether the repo will be private or not.")
	autoInit      = flag.Bool("initialize", true, "Whether to auto initialize the repo or not.")
)

var sourceRepo string
var client *github.Client
var ctx = context.Background()

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
		fmt.Println("Service Name is:", r.FormValue("ServiceName"))
		fmt.Println("Build Type is:", r.FormValue("BuildType"))
		//fmt.Println("Email is:", r.FormValue("Email"))

		//	buildType := r.FormValue("buildType")
		switch BuildType {
		case "gradle":
			fmt.Println("gradle")
			updateFile("pipelineTemplates/azure-gradle-pipeline.yaml", "azure-pipeline.yaml", "BUILDTYPE", BuildType)
		case "maven":
			fmt.Println("maven")
			updateFile("pipelineTemplates/azure-maven-pipeline.yaml", "azure-pipeline.yaml", "BUILDTYPE", BuildType)
		case "vue":
			fmt.Println("vue")
		case "angular":
			fmt.Println("angular")
		case "golang":
			fmt.Println("golang")
		}

		updateFile("pipelineTemplates/sonar.properties", "sonar-project.properties", "SERVICENAME", ServiceName)
		updateFile("pipelineTemplates/buildDefinitionTemplateRequest.json", "buildDefinitionRequest.json", "SERVICENAME", ServiceName)
		flag.StringVar(&sourceRepo, "source-repo", ServiceName, "Name of the repository.")

		createGithubRepo()
		// SEND REQUEST TO ADO USING buildDefinitionRequest.json
		// CREATE OR ADD TO SLACK CHANNEL

		http.Redirect(w, r, "/results", http.StatusFound)
	}
}

func createGithubRepo() {

	flag.Parse()

	token := os.Getenv("GITHUB_AUTH_TOKEN")
	//token := "b4efb01ef973314a0ed8c8ff1453e3186210bf8c"
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	if *sourceOwner == "" || sourceRepo == "" || *sourceFiles == "" || *authorName == "" || *authorEmail == "" {
		log.Fatal("You need to specify a non-empty value for the flags `-source-owner`, `-source-repo`, `-commit-branch`, `-files`, `-author-name` and `-author-email`")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	r := &github.Repository{Name: &sourceRepo, Private: repoPrivacy, AutoInit: autoInit}

	repo, _, err := client.Repositories.Create(ctx, "", r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully created new repo: %v\n", repo.GetName())

	ref, err := getRef()
	if err != nil {
		log.Fatalf("Unable to get/create the commit reference: %s\n", err)
	}

	if ref == nil {
		log.Fatalf("No error where returned but the reference is nil")
	}

	tree, err := getTree(ref)
	if err != nil {
		log.Fatalf("Unable to create the tree based on the provided files: %s\n", err)
	}

	if err := pushCommit(ref, tree); err != nil {
		log.Fatalf("Unable to create the commit: %s\n", err)
	}

	fmt.Printf("Successfully committed pipeline files\n")
}

// getRef returns the base branch reference object
func getRef() (ref *github.Reference, err error) {
	if ref, _, err = client.Git.GetRef(ctx, *sourceOwner, sourceRepo, "refs/heads/"+*baseBranch); err == nil {
		return ref, nil
	}

	return ref, err
}

// getTree generates the tree to commit based on the given files and the commit of the ref you got in getRef.
func getTree(ref *github.Reference) (tree *github.Tree, err error) {
	// Create a tree with what to commit.
	entries := []*github.TreeEntry{}

	// Load each file into the tree.
	for _, fileArg := range strings.Split(*sourceFiles, ",") {
		file, content, err := getFileContent(fileArg)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &github.TreeEntry{Path: github.String(file), Type: github.String("blob"), Content: github.String(string(content)), Mode: github.String("100644")})
	}

	tree, _, err = client.Git.CreateTree(ctx, *sourceOwner, sourceRepo, *ref.Object.SHA, entries)
	return tree, err
}

// getFileContent loads the local content of a file and return the target name
// of the file in the target repository and its contents.
func getFileContent(fileArg string) (targetName string, b []byte, err error) {
	var localFile string
	files := strings.Split(fileArg, ":")
	switch {
	case len(files) < 1:
		return "", nil, errors.New("empty `-files` parameter")
	case len(files) == 1:
		localFile = files[0]
		targetName = files[0]
	default:
		localFile = files[0]
		targetName = files[1]
	}

	b, err = ioutil.ReadFile(localFile)
	return targetName, b, err
}

// createCommit creates the commit in the given reference using the given tree.
func pushCommit(ref *github.Reference, tree *github.Tree) (err error) {
	// Get the parent commit to attach the commit to.
	parent, _, err := client.Repositories.GetCommit(ctx, *sourceOwner, sourceRepo, *ref.Object.SHA)
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()
	author := &github.CommitAuthor{Date: &date, Name: authorName, Email: authorEmail}
	commit := &github.Commit{Author: author, Message: commitMessage, Tree: tree, Parents: []*github.Commit{parent.Commit}}
	newCommit, _, err := client.Git.CreateCommit(ctx, *sourceOwner, sourceRepo, commit)
	if err != nil {
		return err
	}

	// Attach the commit to the master branch.
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, *sourceOwner, sourceRepo, ref, false)
	return err
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
