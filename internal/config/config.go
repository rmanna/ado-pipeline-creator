package config

// Configurations exported
type Configurations struct {
	Server    ServerConfiguration
	Github    GithubConfiguration
	BuildType BuildTypeConfiguration
}

// ServerConfiguration exported
type ServerConfiguration struct {
	Port string
}

// GithubConfiguration exported
type GithubConfiguration struct {
	Owner         string
	Organization  string
	PipelineFiles string
	BaseBranch    string
	AuthorName    string
	AuthorEmail   string
	CommitMessage string
	Token         string
	Permissions   struct {
		Admin string
		Write string
	}
	RepositoryPrivacy bool
	AutoInitialize    bool
}

// BuildTypeConfiguration exported
type BuildTypeConfiguration struct {
	Gradle struct {
		Tasks           string
		Options         string
		JavaHomeOptions string
	}
	Maven struct {
		Options string
		Goals   string
	}
	Vue struct {
		Command string
	}
	Angular struct {
		Command string
	}
	Golang struct {
		Command string
	}
}
