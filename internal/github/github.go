package github

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
)

var (
	client        *github.Client
	ctx           = context.Background()
	sourceRepo    string
	sourceOwner   = flag.String("source-owner", "rmanna", "Name of the owner (user or org) of the repo to create the commit in.")
	baseBranch    = flag.String("base-branch", "master", "Name of branch to create the `commit-branch` from.")
	sourceFiles   = flag.String("files", "azure-pipeline.yaml,sonar-project.properties", "Comma-separated list of files to commit and their location.")
	authorName    = flag.String("author-name", "Self-Service Automation", "Name of the author of the commit.")
	authorEmail   = flag.String("author-email", "oldevops@openlane.com", "Email of the author of the commit.")
	commitMessage = flag.String("commit-message", "Initial Project Commit", "Content of the commit message.")
)

// CreateRepository exported
func CreateRepository(token string, sourceRepository string, repositoryPrivacy bool, autoInitialize bool) {

	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	if sourceRepository == "" {
		log.Fatal("You need to specify a non-empty value for your source repository")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	r := &github.Repository{Name: &sourceRepository, Private: &repositoryPrivacy, AutoInit: &autoInitialize}
	repo, _, err := client.Repositories.Create(ctx, "", r)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully created new repo: %v\n", repo.GetName())
}

// InviteTeams exported
func InviteTeams(teams string, permission string) {

	owner := "KAR-AUTO"
	org := "KAR-AUTO"
	repo := "ol-maven-dummy-app"

	for _, team := range strings.Split(teams, ",") {
		teamopts := &github.TeamAddTeamRepoOptions{Permission: permission}
		_, err := client.Teams.AddTeamRepoBySlug(ctx, org, team, owner, repo, teamopts)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Successfully created team invite: %v\n", team)

	}
}

// AddCollaborator exported
func AddCollaborator(token string, collaborator string, permission string) {
	useropts := &github.RepositoryAddCollaboratorOptions{Permission: "admin"}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
	cinvite, _, err := client.Repositories.AddCollaborator(ctx, "rmanna", "ol-ralph", collaborator, useropts)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully created collaborator: %v\n", cinvite.GetInvitee())
}

// ListOrgTeams exported
func ListOrgTeams(org string, permission string) {
	listopts := &github.ListOptions{Page: 1, PerPage: 100}
	teams, _, err := client.Teams.ListTeams(ctx, org, listopts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", teams)
}

// CommitBranch exported
func CommitBranch(token string, sourceRepository string) {

	flag.StringVar(&sourceRepo, "source-repo", sourceRepository, "Name of the repository.")
	flag.Parse()

	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	if *sourceOwner == "" || sourceRepo == "" || *sourceFiles == "" || *authorName == "" || *authorEmail == "" {
		log.Fatal("You need to specify a non-empty value for the flags `-source-owner`, `-source-repo`, `-commit-branch`, `-files`, `-author-name` and `-author-email`")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

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
