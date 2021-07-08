package get

import (
	"context"
	"log"
	"os"
	"rbac-check/pkg/client"
	"rbac-check/pkg/config"
	"testing"
)

func TestGetUserID(t *testing.T) {
	if os.Getenv("TEST_GITHUB_ACCESS") == "" {
		log.Fatalln("You must have a personal access token set in an env var called 'TEST_GITHUB_ACCESS'")
	}

	user := config.User{
		Username: "cloud-platform-moj",
	}

	opt := config.Options{
		Client: client.GitHubClient(os.Getenv("TEST_GITHUB_ACCESS")),
		Ctx:    context.Background(),
	}

	expected := int64(42068481)
	u, _ := UserID(&opt, &user)

	if int64(*u.ID) != expected {
		t.Errorf("The userID is not expected. want %v, got %v", expected, int64(*u.ID))
	}
}

func TestTeamName(t *testing.T) {
	namespace := "cloud-platform-reports-dev"

	user := config.User{
		Repo: "cloud-platform-environments",
		Org:  "ministryofjustice",
	}

	opt := config.Options{
		Client: client.GitHubClient(os.Getenv("TEST_GITHUB_ACCESS")),
		Ctx:    context.Background(),
	}

	platform := config.Platform{
		AdminTeam:        "webops",
		PrimaryCluster:   "live-1",
		SecondaryCluster: "live",
	}

	teams, _ := TeamName(namespace, &opt, &user, &platform)

	for _, team := range teams {
		if team != platform.AdminTeam {
			t.Errorf("Expecting: %s; got %s", platform.AdminTeam, team)
		}
	}
}
