package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
	yaml "gopkg.in/yaml.v2"
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

	teamNames := make(map[string]int)
	for ns := range namespaces {
		team, _ := getTeamName(token, ns)
		fmt.Println(team)
		if teamNames[team] == 0 {
			teamNames[team] = 1
		} else {
			teamNames[team]++
		}
	}

	userTeams := getUserTeams(token, prOwner)

	fmt.Println(userTeams, namespaces)

	// Get the namespace team name i.e. the name of the team in the rbac file
	// Get a collection of teams the users in
	// For each item in the collection, if it matches with the above rbac team name; pass
	// else; fail
}

func getTeamName(token, namespace string) (string, error) {
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
	path := "namespaces/live-1.cloud-platform.service.justice.gov.uk/abundant-namespace-dev/01-rbac.yaml"

	opts := &github.RepositoryContentGetOptions{}
	// reader, _, err := client.Repositories.DownloadContents(ctx, "ministryofjustice", "cloud-platform-environments", path, opts)
	fileContent, _, _, err := client.Repositories.GetContents(ctx, "ministryofjustice", "cloud-platform-environments", path, opts)
	if err != nil {
		log.Println(err)
	}

	cont, err := fileContent.GetContent()

	fullName := Rbac{}

	err = yaml.Unmarshal([]byte(cont), &fullName)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("--- t:\n%v\n\n", fullName)

	return "", nil
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
