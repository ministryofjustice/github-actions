// main is expected to run in a GitHub action. It will expect to receive a
// file in its directory called `files` that contains a space seperated list of
// MoJ Cloud Platform Kubernetes namespaces changed in a PR. Each namespace contains
// an rbac file that contains a list of team names. If the PR owner (whoever raised the PR)
// is a member of that list, the package returns true. If not, it'll return false.
package main

import (
	"context"
	"log"
	"os"

	"rbac-check/pkg/client"
	"rbac-check/pkg/config"
	"rbac-check/pkg/get"
	"rbac-check/pkg/validate"

	ghaction "github.com/sethvargo/go-githubactions"
)

func main() {
	// Exit hard if the environment variables don't exist. The package requires
	// a personal access token with ORG permissions. PR_OWNER is passed upstream
	// by a GitHub Action.
	if os.Getenv("GITHUB_OAUTH_TOKEN") == "" || os.Getenv("PR_OWNER") == "" {
		log.Fatalln("you must have the GITHUB_OAUTH_TOKEN and PR_OWNER env var set.")
	}

	user := config.User{
		Username: os.Getenv("PR_OWNER"),
		Branch:   os.Getenv("BRANCH"),
		Repo:     "cloud-platform-environments",
		Org:      "ministryofjustice",
	}

	opt := config.Options{
		Client:   client.GitHubClient(os.Getenv("GITHUB_OAUTH_TOKEN")),
		Ctx:      context.Background(),
		FileName: "files",
	}

	platform := config.Platform{
		AdminTeam:        "WebOps",
		PrimaryCluster:   "live-1",
		SecondaryCluster: "live",
	}

	// On the condition where there is no fileName i.e. it isn't created upstream,
	// this package should return a pass. By doing this we're assuming the PR contains
	// changes that don't include rbac files and thus can be reviewed.
	_, err := os.Stat(opt.FileName)
	if os.IsNotExist(err) {
		log.Println("File doesn't exist. Passing.")
		ghaction.SetOutput("review_pr", "true")
		os.Exit(0)
	}

	// Parse all namespaces in the fileName variable. Fail if we can't parse because later
	// functions require this output.
	namespaces, err := get.Namespaces(opt.FileName)
	if err != nil {
		log.Fatalln("Unable to fetch namespace:", err)
	}

	// Call the GitHub API for the rbac team name of each namespace in the namespaces variable.
	// The teams are stored in a map to ensure we don't duplicate. I found maps are best for
	// deduplication in Go.
	namespaceTeams := make(map[string]int)
	for ns := range namespaces {
		teams, err := get.TeamName(ns, &opt, &user, &platform)
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

	// Add the WebOps team as admins over all namespaces
	namespaceTeams[platform.AdminTeam] = 1

	// Convert the PR_OWNER string into a GitHub user ID. This is used later, to compare the
	// list of users in a team.
	userID, err := get.UserID(&opt, &user)
	if err != nil {
		log.Fatalln("Unable to fetch userID", err)
	}

	// Call the GitHub API to confirm if the user exists in the GitHub team name.
	valid, team, err := validate.UserPermissions(namespaceTeams, userID, &opt, &user)
	if err != nil {
		log.Fatalln("Unable to check if the user is valid:", err)
	}

	// Send result back to GitHub Action.
	if valid {
		log.Println("\n The user:", userID.GetName(), "\n is in team:", team)
		ghaction.SetOutput("review_pr", "true")
	} else {
		log.Println("\n The user:", userID.GetName(), "\n can't be found in teams:", namespaceTeams)
		ghaction.SetOutput("review_pr", "false")
	}
}
