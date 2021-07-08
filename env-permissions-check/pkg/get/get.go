// Package get is a getter package for all GitHub API calls
package get

import (
	"bufio"
	"log"
	"os"
	"rbac-check/pkg/config"
	"strings"

	"github.com/google/go-github/v35/github"
	"gopkg.in/yaml.v2"
)

// UserID takes an Options and User data type from the caller and returns
// the users GitHub user object.
func UserID(opt *config.Options, user *config.User) (*github.User, error) {
	userID, _, err := opt.Client.Users.Get(opt.Ctx, user.Username)
	if err != nil {
		return nil, err
	}

	return userID, nil
}

func origin(namespace string, opt *config.Options, user *config.User, platform *config.Platform, repoOpts *github.RepositoryContentGetOptions) (string, error) {
	secondaryCluster := platform.SecondaryCluster
	primaryCluster := platform.PrimaryCluster

	cluster := primaryCluster
	user.Path = "namespaces/" + cluster + ".cloud-platform.service.justice.gov.uk/" + namespace + "/01-rbac.yaml"

	_, _, resp, err := opt.Client.Repositories.GetContents(opt.Ctx, user.Org, user.Repo, user.Path, repoOpts)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 200 {
		return cluster, nil
	} else {
		cluster = secondaryCluster
		_, _, resp, err := opt.Client.Repositories.GetContents(opt.Ctx, user.Org, user.Repo, user.Path, repoOpts)
		if err != nil {
			return "", err
		}
		if resp.StatusCode == 200 {
			return cluster, nil
		}
	}

	return "none", nil
}

func TeamName(namespace string, opt *config.Options, user *config.User, platform *config.Platform) ([]string, error) {
	repoOpts := &github.RepositoryContentGetOptions{}
	// 3 cases here. Live-1, live if a namespace doesn't yet exist

	// is it live, live-1 or doesn't it exist
	// Find out if it's live-1, live or in the PR.
	ori, err := origin(namespace, opt, user, platform, repoOpts)
	if err != nil {
		log.Println(err)
	}

	// if the namespace doesn't exist yet, check the pr.
	if ori == "none" {
		repoOpts = &github.RepositoryContentGetOptions{
			Ref: user.Branch,
		}
		ori, err = origin(namespace, opt, user, platform, repoOpts)
		if err != nil {
			return nil, err
		}
	}

	user.Path = "namespaces/" + ori + ".cloud-platform.service.justice.gov.uk/" + namespace + "/01-rbac.yaml"
	file, _, _, err := opt.Client.Repositories.GetContents(opt.Ctx, user.Org, user.Repo, user.Path, repoOpts)
	if err != nil {
		return nil, err
	}

	cont, err := file.GetContent()
	if err != nil {
		return nil, err
	}

	fullName := config.Rbac{}

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

func Namespaces(fileName string) (map[string]int, error) {
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
