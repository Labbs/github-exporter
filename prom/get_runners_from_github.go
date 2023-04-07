package prom

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v51/github"
	"github.com/labbs/github-exporter/config"
	"github.com/rs/zerolog"
)

func GetRunnersFromGhithub(logger zerolog.Logger, client *github.Client) {
	opt := &github.ListOptions{PerPage: 200}

	for _, repo := range config.Github.Repositories.Value() {
		var runners []*github.Runner
		r := strings.Split(repo, "/")

		for {
			resp, rr, err := client.Actions.ListRunners(context.Background(), r[0], r[1], opt)
			if rlerr, ok := err.(*github.RateLimitError); ok {
				logger.Info().Err(rlerr).Str("event", "get_runners_from_github").Msg("Rate limit error, waiting for reset until " + rlerr.Rate.Reset.String())
				time.Sleep(time.Until(rlerr.Rate.Reset.Time))
				continue
			} else if err != nil {
				logger.Error().Err(err).Str("event", "get_runners_from_github").Str("repo", repo).Msg("Error to get runners from github")
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
				RunnersGauge.WithLabelValues(repo, runner.GetOS(), runner.GetName(), strconv.FormatInt(runner.GetID(), 10), strconv.FormatBool(runner.GetBusy())).Set(1)
			} else {
				RunnersGauge.WithLabelValues(repo, runner.GetOS(), runner.GetName(), strconv.FormatInt(runner.GetID(), 10), strconv.FormatBool(runner.GetBusy())).Set(0)
			}
		}
	}
}
