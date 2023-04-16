package prom

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/v51/github"
	"github.com/labbs/github-exporter/config"
)

func GithubFetcher(client *github.Client) {

	// Get repositories if not set in config
	var repos []string
	if len(config.Github.Repositories.Value()) > 0 {
		repos = config.Github.Repositories.Value()
	} else {
		for _, orga := range config.Github.Organizations.Value() {
			repos = append(repos, getReposFromOrganization(client, orga)...)
		}
	}

	// Get workflows
	non_empty_repos := make([]string, 0)
	ww := make(map[string]map[int64]github.Workflow)
	for _, repo := range repos {
		r := strings.Split(repo, "/")
		workflows_for_repo := getWorkflowsFromRepository(client, r[0], r[1])
		if len(workflows_for_repo) > 0 {
			ww[repo] = workflows_for_repo
			non_empty_repos = append(non_empty_repos, repo)
			Logger.Debug().Str("event", "get_workflows_from_repository").Str("repository", repo).Int("count", len(workflows_for_repo)).Msg("Workflows found")
		}
	}
	repositories = non_empty_repos
	workflows = ww
}

func getReposFromOrganization(client *github.Client, orga string) []string {
	var repos []string

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 200, Page: 0},
	}

	for {
		repositories, resp, err := client.Repositories.ListByOrg(context.Background(), orga, opt)
		if rlerr, ok := err.(*github.RateLimitError); ok {
			Logger.Info().Err(rlerr).Str("event", "get_repos_from_organization").Msg("Rate limit error, waiting for reset until " + rlerr.Rate.Reset.String())
			time.Sleep(time.Until(rlerr.Rate.Reset.Time))
			continue
		} else if err != nil {
			Logger.Error().Err(err).Str("event", "get_repos_from_organization").Msg("Error to get repositories from organization")
			break
		}

		for _, repo := range repositories {
			repos = append(repos, *repo.FullName)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.ListOptions.Page = resp.NextPage
	}

	return repos
}

func getWorkflowsFromRepository(client *github.Client, owner, repo string) map[int64]github.Workflow {
	workflows := make(map[int64]github.Workflow)

	opt := &github.ListOptions{PerPage: 200, Page: 0}

	for {
		workflowRuns, resp, err := client.Actions.ListWorkflows(context.Background(), owner, repo, opt)
		if rlerr, ok := err.(*github.RateLimitError); ok {
			Logger.Info().Err(rlerr).Str("event", "get_workflows_from_repository").Msg("Rate limit error, waiting for reset until " + rlerr.Rate.Reset.String())
			time.Sleep(time.Until(rlerr.Rate.Reset.Time))
			continue
		} else if err != nil {
			Logger.Error().Err(err).Str("event", "get_workflows_from_repository").Msg("Error to get workflows from repository")
			break
		}

		for _, w := range workflowRuns.Workflows {
			workflows[*w.ID] = *w
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return workflows
}
