package validate

import (
	"rbac-check/pkg/config"

	"github.com/google/go-github/v35/github"
)

func UserPermissions(namespaceTeams map[string]int, githubUser *github.User, opt *config.Options, user *config.User) (bool, string, error) {
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
