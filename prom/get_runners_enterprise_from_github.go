package prom

import (
	"context"
	"strconv"
	"time"

	"github.com/google/go-github/v51/github"
	"github.com/labbs/github-exporter/config"
)

func GetRunnersEnterpriseFromGithub(client *github.Client) {
	var runners []*github.Runner
	opt := &github.ListOptions{PerPage: 200}

	for {
		resp, rr, err := client.Enterprise.ListRunners(context.Background(), config.Github.EnterpriseName, nil)
		if rlerr, ok := err.(*github.RateLimitError); ok {
			Logger.Info().Err(rlerr).Str("event", "get_runners_enterprise_from_github").Msg("Rate limit error, waiting for reset until " + rlerr.Rate.Reset.String())
			time.Sleep(time.Until(rlerr.Rate.Reset.Time))
			continue
		} else if err != nil {
			Logger.Error().Err(err).Str("event", "get_runners_enterprise_from_github").Msg("Error to get runners enterprise from github")
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
			RunnersEnterpriseGauge.WithLabelValues(*runner.OS, *runner.Name, strconv.FormatInt((runner.GetID()), 10)).Set(1)
		} else {
			RunnersEnterpriseGauge.WithLabelValues(*runner.OS, *runner.Name, strconv.FormatInt((runner.GetID()), 10)).Set(0)
		}
	}
}
