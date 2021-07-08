// Package config implements internal configuration.
package config

import (
	"context"

	"github.com/google/go-github/v35/github"
)

// Rbac type is used to parse a yaml file and a list of subjects in a RoleBinding.
type Rbac struct {
	Subjects []Subjects `yaml:"subjects"`
}

// Subjects type is a child of Rbac and contains the name of the GitHub team.
type Subjects struct {
	Name string `yaml:"name"`
}

// Options defines the contents of a github client, context and other variables.
type Options struct {
	Client *github.Client
	Ctx    context.Context
	// FileName defines the file name passed by the GitHub Action.
	FileName string
}

// User defines the structure of the user of the PR.
type User struct {
	// Username is the creator of the PR, passed by a GitHub Action.
	Username string
	// Branch is the git branch the users PR lives on.
	Branch string
	// Repo contains the string value of the repository this is executed against.
	Repo string
	// Org is the Organisation the repository/Repo lives in.
	Org string
	// Path is the file path of the rbac file in a namespace.
	Path string
}

type Platform struct {
	// AdminTeam is the admin of the users repository, i.e. can PR in all namespaces.
	AdminTeam string
	// PrimaryCluster is the cluster name where most MoJ Kubernetes namespaces are created.
	PrimaryCluster string
	// SecondaryCluster is the cluster name of the backup or future cluster for an MoJ Kubernetes namespace.
	SecondaryCluster string
}
