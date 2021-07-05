package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type Rbac struct {
	Subjects []Subjects `yaml:"subjects"`
}

// Subjects
type Subjects struct {
	Name string `yaml:"name"`
}

func main() {
	var (
		fileName = "files"
		token    = os.Getenv("GITHUB_OAUTH_TOKEN")
		prOwner  = os.Getenv("PR_OWNER")
		branch   = os.Getenv("BRANCH")
		// valid    = false
	)

	if os.Getenv("GITHUB_OAUTH_TOKEN") == "" || os.Getenv("PR_OWNER") == "" {
		log.Fatalln("you must have the GITHUB_OAUTH_TOKEN and PR_OWNER env var set.")
	}

	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		log.Println("File doesn't exist. Passing.")
	}

	namespaces, err := getNamespaces(fileName)
	if err != nil {
		log.Fatalln("Unable to fetch namespace:", err)
	}

	namespaceTeams := make(map[string]int)
	for ns := range namespaces {
		teams, err := getTeamName(token, ns, branch)
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

	userID, err := getUserID(prOwner, token)
	if err != nil {
		log.Println("Unable to fetch userID", err)
	}

	orgID, err := getOrgID(token)
	if err != nil {
		log.Println("Unable to fetch userID", err)
	}

	userTeams := getUserTeams(token, prOwner)

	fmt.Println(userTeams, namespaces)

	// Get the namespace team name i.e. the name of the team in the rbac file
	// Get a collection of teams the users in
	// For each item in the collection, if it matches with the above rbac team name; pass
	// else; fail
}

func getOrgID(token string) (*github.Organization, error) {
	orgOwner := "ministryofjustice"
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Fetch the orgOwner user ID.
	org, _, err := client.Organizations.Get(ctx, orgOwner)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func getUserID(prOwner, token string) (*github.User, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Fetch the user's GitHub user ID.
	user, _, err := client.Users.Get(ctx, prOwner)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func getOrigin(namespace string, ctx context.Context, client *github.Client, opts *github.RepositoryContentGetOptions) (string, error) {
	secondaryCluster := "live"
	primaryCluster := "live-1"

	cluster := primaryCluster
	path := "namespaces/" + cluster + ".cloud-platform.service.justice.gov.uk/" + namespace + "/01-rbac.yaml"

	_, _, resp, err := client.Repositories.GetContents(ctx, "ministryofjustice", "cloud-platform-environments", path, opts)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 200 {
		return cluster, nil
	} else {
		cluster = secondaryCluster
		_, _, resp, err := client.Repositories.GetContents(ctx, "ministryofjustice", "cloud-platform-environments", path, opts)
		if err != nil {
			return "", err
		}
		if resp.StatusCode == 200 {
			return cluster, nil
		}
	}

	return "none", nil
}

func getTeamName(token, namespace, branch string) ([]string, error) {
	// call the github api for the namespace passed to get yaml
	// parse the yaml and get subject name
	// strip the name so it appears as webops note: this is all lowercase
	// return it as a string

	// Setting up GitHub client.
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	opts := &github.RepositoryContentGetOptions{}
	// 3 cases here. Live-1, live if a namespace doesn't yet exist

	// is it live, live-1 or doesn't it exist
	// Find out if it's live-1, live or in the PR.
	origin, err := getOrigin(namespace, ctx, client, opts)
	if err != nil {
		log.Println(err)
	}

	// if the namespace doesn't exist yet, check the pr.
	if origin == "none" {
		opts := &github.RepositoryContentGetOptions{
			Ref: branch,
		}
		origin, err = getOrigin(namespace, ctx, client, opts)
		if err != nil {
			return nil, err
		}
	}

	path := "namespaces/" + origin + ".cloud-platform.service.justice.gov.uk/" + namespace + "/01-rbac.yaml"
	file, _, _, err := client.Repositories.GetContents(ctx, "ministryofjustice", "cloud-platform-environments", path, opts)
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

	// get all pr file changes
	// grep the namespace name
	// grep the rbac file within the namespace name
	// return the team name

}

func getUserTeams(token, prOwner string) []string {
	// Setting up GitHub client.
	// ctx := context.Background()
	// ts := oauth2.StaticTokenSource(
	// 	&oauth2.Token{AccessToken: token},
	// )

	// tc := oauth2.NewClient(ctx, ts)

	// _ = github.NewClient(tc)

	// We have the username
	// we need all teams the user's in
	// grab the user id
	// perform a lookup for all teams the users in. Add to collection.
	// return collection.
	return nil

}
