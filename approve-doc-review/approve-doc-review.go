package main

import (
	"flag"
	"log"
	"os"

	"approve-doc-review/check"

	ghaction "github.com/sethvargo/go-githubactions"
)

// main is intended to be used in a GitHub action. It expects a file containing all changes
// made in a PR, an environment variable containing the username of the author/owner of the pr and
// a GitHub token that can perform organisation lookups via the GitHub API. At a high level it will
// return a truthy value if the PR:
// - contains only changes containing the term "last_reviewed_on"
// - the author/owner of the PR is in the "ministryofjustice/WebOps" team.
func main() {
	var (
		team     = flag.String("team", "WebOps", "team and orgOwner are the GitHub team and organisation that we're using to validate the user.")
		prOwner  = flag.String("user", os.Getenv("PR_OWNER"), "contains the value of an environment variable that is set in the GH action container")
		orgOwner = flag.String("org", "ministryofjustice", "who owns the repository")
	)

	var (
		fileName = flag.String("file", "changes", "the file created by a GitHub action, it contains the output of a git diff")
		token    = flag.String("token", os.Getenv("GITHUB_OAUTH_TOKEN"), "Personal access token from GitHub.")
	)

	flag.Parse()

	if os.Getenv("GITHUB_OAUTH_TOKEN") == "" || os.Getenv("PR_OWNER") == "" {
		log.Fatalln("you must have the GITHUB_OAUTH_TOKEN and PR_OWNER env var set.")
	}

	// prRelevant will return true or false depending on the contents of fileName. We don't want
	// the GH action to error here so we just log the error and take no action.
	prRelevant, err := check.ParsePR(*fileName)
	if err != nil {
		log.Println("Unable to parse the PR", err)
	}

	// userAllowed will return true of false depending on the owner of the PR. We don't want
	// the GH action to error here so we just log the error and take no action.
	userAllowed, err := check.IsInGitHubTeam(*team, *orgOwner, *token, *prOwner)
	if err != nil {
		log.Println("Unable to check if the user is valid.", err)
	}

	// Conditional check to see if we should pass or fail the step. We don't want a hard fail so we set
	// the output to false and log.
	if prRelevant && userAllowed {
		log.Println("Success: The changes in this PR are only review dates and the user is valid.")
		ghaction.SetOutput("review_pr", "true")
	} else {
		log.Println("Fail: Either the PR contains more than review date changes or the user isn't a member of the webops team.")
		ghaction.SetOutput("review_pr", "false")
	}
}
