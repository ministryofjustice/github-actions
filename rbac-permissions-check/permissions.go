// main is expected to run in a GitHub action. It will expect to receive a
// file in its directory called `files` that contains a space seperated list of
// MoJ Cloud Platform Kubernetes namespaces changed in a PR. Each namespace contains
// an rbac file that contains a list of team names. If the PR owner (whoever raised the PR)
// is a member of that list, the package returns true. If not, it'll return false.
package main

import (
	"context"
	"flag"
	"log"
	"os"

	"rbac-check/pkg/client"
	"rbac-check/pkg/config"
	"rbac-check/pkg/get"
	"rbac-check/pkg/validate"

	ghaction "github.com/sethvargo/go-githubactions"
)

// Default values have been set to ministryofjustice Cloud-Platform specifics.
var (
	token            = flag.String("token", os.Getenv("GITHUB_OAUTH_TOKEN"), "Personal access token from GitHub.")
	branch           = flag.String("branch", os.Getenv("BRANCH"), "Branch of changes in GitHub.")
	username         = flag.String("user", os.Getenv("PR_OWNER"), "Branch of changes in GitHub.")
	repo             = flag.String("repository", "cloud-platform-environments", "The repository of the Cloud Platform repository.")
	org              = flag.String("org", "ministryofjustice", "Name of the orgnanisation i.e. ministryofjustice.")
	file             = flag.String("file", "files", "File name containing namespaces with changes.")
	adminTeam        = flag.String("admin", "WebOps", "Admin team looking after repository.")
	primaryCluster   = flag.String("primary", "live-1", "Name of the primary cluster in use.")
	secondaryCluster = flag.String("secondary", "live", "Name of the secondary cluster in use.")
)

func main() {
	flag.Parse()

	// Fail if the relevant flags or environment variables haven't set.
	// token = a personal access token from GitHub
	// branch = the branch name of your PR
	// username = the github username used to create the PR.
	// All of these values will be passed upstream by a GitHub action.
	if *token == "" || *branch == "" || *username == "" {
		log.Fatalln("You need to specify a non-empty value for token, branch and username.")
	}

	user := config.User{
		PrimaryCluster:   *primaryCluster,
		SecondaryCluster: *secondaryCluster,
		Username:         *username,
	}

	opt := config.Options{
		Client: client.GitHubClient(*token),
		Ctx:    context.Background(),
	}

	repo := config.Repository{
		AdminTeam:    *adminTeam,
		Branch:       *branch,
		ChangedFiles: *file,
		Name:         *repo,
		Org:          *org,
	}

	// On the condition where there is no file ChangedFiles i.e. it hasn't been created upstream,
	// this package should return a pass. By doing this we're assuming the PR contains
	// changes that don't include rbac files and thus can be reviewed.
	_, err := os.Stat(repo.ChangedFiles)
	if os.IsNotExist(err) {
		log.Println("File doesn't exist. Passing.")
		ghaction.SetOutput("review_pr", "true")
		os.Exit(0)
	}

	// Parse all namespaces in the ChangedFiles variable. Fail if we can't parse because later
	// functions require this output.
	user.Namespaces, err = get.Namespaces(repo.ChangedFiles)
	if err != nil {
		log.Fatalln("Unable to fetch namespace:", err)
	}

	// Call the GitHub API for the rbac team name of each namespace in the namespaces variable.
	// The teams are stored in a map to ensure we don't duplicate. Maps are best for
	// deduplication in Go.
	namespaceTeams := make(map[string]int)
	for ns := range user.Namespaces {
		teams, err := get.TeamName(ns, &opt, &user, &repo)
		if err != nil {
			log.Fatalln("Unable to get team names:", err)
		}
		for _, team := range teams {
			if namespaceTeams[team] == 0 {
				namespaceTeams[team] = 1
			} else {
				namespaceTeams[team]++
			}
		}
	}

	// Add the admin team so they can make changes to any namespace.
	namespaceTeams[repo.AdminTeam] = 1

	// Convert the username string into a GitHub user ID. This is used later, to compare the
	// list of users in a team.
	user.Id, err = get.UserID(&opt, &user)
	if err != nil {
		log.Fatalln("Unable to fetch userID", err)
	}

	// Call the GitHub API to confirm if the user exists in the GitHub team name.
	valid, team, err := validate.UserPermissions(namespaceTeams, &opt, &user, &repo)
	if err != nil {
		log.Fatalln("Unable to check if the user is valid:", err)
	}

	// Send result back to GitHub Action.
	if valid {
		log.Println("\n The user:", user.Id.GetName(), "\n is in team:", team)
		ghaction.SetOutput("review_pr", "true")
	} else {
		log.Println("\n The user:", user.Id.GetName(), "\n can't be found in teams:", namespaceTeams)
		ghaction.SetOutput("review_pr", "false")
	}
}
