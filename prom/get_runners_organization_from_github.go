package prom

import (
	"context"
	"strconv"
	"time"

	"github.com/google/go-github/v51/github"
	"github.com/labbs/github-exporter/config"
)

func GetRunnersOrganizationFromGithub(client *github.Client) {
	RunnersOrganizationGauge.Reset()
	opt := &github.ListOptions{PerPage: 200}

	for _, orga := range config.Github.Organizations.Value() {
		var runners []*github.Runner

		for {
			resp, rr, err := client.Actions.ListOrganizationRunners(context.Background(), orga, opt)
			if rlerr, ok := err.(*github.RateLimitError); ok {
				Logger.Info().Err(rlerr).Str("event", "get_runners_organization_from_github").Msg("Rate limit error, waiting for reset until " + rlerr.Rate.Reset.String())
				time.Sleep(time.Until(rlerr.Rate.Reset.Time))
				continue
			} else if err != nil {
				Logger.Error().Err(err).Str("event", "get_runners_organization_from_github").Str("organization", orga).Msg("Error to get runners organization from github")
				break
			}
			runners = append(runners, resp.Runners...)
			if rr.NextPage == 0 {
				break
			}
			opt.Page = rr.NextPage
		}

		for _, runner := range runners {
			if runner.GetStatus() == "online" {
				RunnersOrganizationGauge.WithLabelValues(orga, runner.GetOS(), runner.GetName(), strconv.FormatInt(runner.GetID(), 10), strconv.FormatBool(runner.GetBusy())).Set(1)
			} else {
				RunnersOrganizationGauge.WithLabelValues(orga, runner.GetOS(), runner.GetName(), strconv.FormatInt(runner.GetID(), 10), strconv.FormatBool(runner.GetBusy())).Set(0)
			}
		}
	}
}
