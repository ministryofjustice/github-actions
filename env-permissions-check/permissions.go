// Package main is expected to run in a GitHub action. It will expect to receive a
// file in its directory called `files` that contains a space seperated list of
// MoJ Cloud Platform Kubernetes namespaces changed in a PR. Each namespace contains
// an rbac file that contains a list of team names. If the PR owner (whoever raised the PR)
// is a member of that list, the package returns true. If not, it'll return false.
package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"

	"rbac-check/pkg/client"

	"github.com/google/go-github/v35/github"
	ghaction "github.com/sethvargo/go-githubactions"
	"gopkg.in/yaml.v2"
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
	// AdminTeam is the admin of the users repository, i.e. can PR in all namespaces.
	AdminTeam string
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
}

func main() {
	// Exit hard if the environment variables don't exist. The package requires
	// a personal access token with ORG permissions. PR_OWNER is passed upstream
	// by a GitHub Action.
	if os.Getenv("GITHUB_OAUTH_TOKEN") == "" || os.Getenv("PR_OWNER") == "" {
		log.Fatalln("you must have the GITHUB_OAUTH_TOKEN and PR_OWNER env var set.")
	}

	user := User{
		Username: os.Getenv("PR_OWNER"),
		Branch:   os.Getenv("BRANCH"),
		Repo:     "cloud-platform-environments",
		Org:      "ministryofjustice",
	}

	opt := Options{
		Client:    client.GitHubClient(os.Getenv("GITHUB_OAUTH_TOKEN")),
		Ctx:       context.Background(),
		FileName:  "files",
		AdminTeam: "WebOps",
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
	namespaces, err := getNamespaces(opt.FileName)
	if err != nil {
		log.Fatalln("Unable to fetch namespace:", err)
	}

	// Call the GitHub API for the rbac team name of each namespace in the namespaces variable.
	// The teams are stored in a map to ensure we don't duplicate. I found maps are best for
	// deduplication in Go.
	namespaceTeams := make(map[string]int)
	for ns := range namespaces {
		teams, err := getTeamName(ns, &opt, &user)
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
	namespaceTeams[opt.AdminTeam] = 1

	// Convert the PR_OWNER string into a GitHub user ID. This is used later, to compare the
	// list of users in a team.
	userID, err := getUserID(&opt, &user)
	if err != nil {
		log.Fatalln("Unable to fetch userID", err)
	}

	// Call the GitHub API to confirm if the user exists in the GitHub team name.
	valid, team, err := isUserValid(namespaceTeams, userID, &opt, &user)
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

func getUserID(opt *Options, user *User) (*github.User, error) {
	userID, _, err := opt.Client.Users.Get(opt.Ctx, user.Username)
	if err != nil {
		return nil, err
	}

	return userID, nil
}

func getOrigin(namespace string, opt *Options, user *User, repoOpts *github.RepositoryContentGetOptions) (string, error) {
	secondaryCluster := "live"
	primaryCluster := "live-1"

	cluster := primaryCluster
	path := "namespaces/" + cluster + ".cloud-platform.service.justice.gov.uk/" + namespace + "/01-rbac.yaml"

	_, _, resp, err := opt.Client.Repositories.GetContents(opt.Ctx, user.Org, user.Repo, path, repoOpts)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 200 {
		return cluster, nil
	} else {
		cluster = secondaryCluster
		_, _, resp, err := opt.Client.Repositories.GetContents(opt.Ctx, user.Org, user.Repo, path, repoOpts)
		if err != nil {
			return "", err
		}
		if resp.StatusCode == 200 {
			return cluster, nil
		}
	}

	return "none", nil
}

func getTeamName(namespace string, opt *Options, user *User) ([]string, error) {
	repoOpts := &github.RepositoryContentGetOptions{}
	// 3 cases here. Live-1, live if a namespace doesn't yet exist

	// is it live, live-1 or doesn't it exist
	// Find out if it's live-1, live or in the PR.
	origin, err := getOrigin(namespace, opt, user, repoOpts)
	if err != nil {
		log.Println(err)
	}

	// if the namespace doesn't exist yet, check the pr.
	if origin == "none" {
		repoOpts = &github.RepositoryContentGetOptions{
			Ref: user.Branch,
		}
		origin, err = getOrigin(namespace, opt, user, repoOpts)
		if err != nil {
			return nil, err
		}
	}

	path := "namespaces/" + origin + ".cloud-platform.service.justice.gov.uk/" + namespace + "/01-rbac.yaml"
	file, _, _, err := opt.Client.Repositories.GetContents(opt.Ctx, user.Org, user.Repo, path, repoOpts)
	if err != nil {
		return nil, err
	}

	cont, err := file.GetContent()
	if err != nil {
		return nil, err
	}

	fullName := Rbac{}

	err = yaml.Unmarshal([]byte(cont), &fullName)
	if err != nil {
		return nil, err
	}

	var namespaceTeams []string
	for _, name := range fullName.Subjects {
		str := strings.SplitAfter(string(name.Name), ":")
		namespaceTeams = append(namespaceTeams, str[1])
	}

	return namespaceTeams, nil
}

func getNamespaces(fileName string) (map[string]int, error) {
	namespaces := make(map[string]int)

	file, err := os.Open(fileName)
	if err != nil {
		return namespaces, err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		if namespaces[scanner.Text()] == 0 {
			namespaces[scanner.Text()] = 1
		} else {
			namespaces[scanner.Text()]++
		}
	}

	return namespaces, nil
}

func isUserValid(namespaceTeams map[string]int, githubUser *github.User, opt *Options, user *User) (bool, string, error) {
	teamOpts := &github.TeamListTeamMembersOptions{}
	for team := range namespaceTeams {
		members, _, err := opt.Client.Teams.ListTeamMembersBySlug(opt.Ctx, user.Org, team, teamOpts)
		if err != nil {
			return false, "", nil
		}
		for _, member := range members {
			if member.GetID() == githubUser.GetID() {
				return true, team, nil
			}
		}
	}

	return false, "", nil
}
