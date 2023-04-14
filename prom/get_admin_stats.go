package prom

import (
	"context"

	"github.com/google/go-github/v51/github"
)

// GetAdminStats get admin stats from github enterprise
func GetAdminStats(client *github.Client) {
	resp, _, err := client.Admin.GetAdminStats(context.Background())
	if err != nil {
		Logger.Error().Err(err).Str("event", "get_repos").Msg("Error to get repos from github")
		return
	}

	AdminStatsGauge.WithLabelValues("issues", "open").Set(float64(resp.Issues.GetOpenIssues()))
	AdminStatsGauge.WithLabelValues("issues", "closed").Set(float64(resp.Issues.GetClosedIssues()))
	AdminStatsGauge.WithLabelValues("issues", "total").Set(float64(resp.Issues.GetTotalIssues()))
	AdminStatsGauge.WithLabelValues("users", "admin").Set(float64(resp.Users.GetAdminUsers()))
	AdminStatsGauge.WithLabelValues("users", "suspended").Set(float64(resp.Users.GetSuspendedUsers()))
	AdminStatsGauge.WithLabelValues("users", "total").Set(float64(resp.Users.GetTotalUsers()))
	AdminStatsGauge.WithLabelValues("repos", "total").Set(float64(resp.Repos.GetTotalRepos()))
	AdminStatsGauge.WithLabelValues("repos", "fork").Set(float64(resp.Repos.GetForkRepos()))
	AdminStatsGauge.WithLabelValues("repos", "org").Set(float64(resp.Repos.GetOrgRepos()))
	AdminStatsGauge.WithLabelValues("repos", "root").Set(float64(resp.Repos.GetRootRepos()))
	AdminStatsGauge.WithLabelValues("repos", "total_pushes").Set(float64(resp.Repos.GetTotalPushes()))
	AdminStatsGauge.WithLabelValues("repos", "total_wikie").Set(float64(resp.Repos.GetTotalWikis()))
	AdminStatsGauge.WithLabelValues("orgs", "disabled").Set(float64(resp.Orgs.GetDisabledOrgs()))
	AdminStatsGauge.WithLabelValues("orgs", "total").Set(float64(resp.Orgs.GetTotalOrgs()))
	AdminStatsGauge.WithLabelValues("orgs", "teams").Set(float64(resp.Orgs.GetTotalTeams()))
	AdminStatsGauge.WithLabelValues("orgs", "team_members").Set(float64(resp.Orgs.GetTotalTeamMembers()))
}
