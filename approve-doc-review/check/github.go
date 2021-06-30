// Check is used to perform bool checks on whether a PR and its user
// are valid to be auto approved.
package check

import (
	"context"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

// IsUserInWebOps takes a GitHub token and a string containing the owner of the
// PR. It'll perform a lookup for the following:
// - The PR owner's GitHub user ID, denoted by the user var.
// - The org user ID, denoted by org. We could've hard coded this.
// - The team ID, denoted by teamId.
// It then performs a lookup and confirms if the user's ID appears in the team ID memebers.
// If so, this function returns true.
func IsUserInWebOps(team, orgOwner, token, prOwner string) (bool, error) {
	// Setting up GitHub client.
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Fetch the user's GitHub user ID.
	user, _, err := client.Users.Get(ctx, prOwner)
	if err != nil {
		return false, err
	}

	// Fetch the orgOwner user ID.
	org, _, err := client.Organizations.Get(ctx, orgOwner)
	if err != nil {
		return false, err
	}

	// Fetch the team's GitHub ID.
	opts := &github.ListOptions{}
	var teamId int64
	for {
		teams, resp, err := client.Teams.ListTeams(ctx, orgOwner, opts)
		if err != nil {
			return false, err
		}
		for _, t := range teams {
			if t.GetName() == team {
				teamId = t.GetID()
				break
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	// Create a collection of members of the team variable.
	teamOpts := &github.TeamListTeamMembersOptions{}
	teamMembers, _, err := client.Teams.ListTeamMembersByID(ctx, *org.ID, teamId, teamOpts)
	if err != nil {
		return false, err
	}

	// Loop the collection of team members and check to see if the prOwner is a member.
	for _, v := range teamMembers {
		if int64(*v.ID) == int64(*user.ID) {
			return true, nil
		}
	}

	return false, nil
}
