package config

// Configurations exported
type Configurations struct {
	Server ServerConfigurations
	Github GithubConfigurations
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port int
}

// GithubConfigurations exported
type GithubConfigurations struct {
	SourceOwner       string
	SourceFiles       string
	BaseBranch        string
	AuthorName        string
	AuthorEmail       string
	CommitMessage     string
	Token             string
	RepositoryPrivacy bool
	AutoInitialize    bool
}
