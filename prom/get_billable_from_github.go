package prom

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v51/github"
)

func GetBillableFromGithub(client *github.Client) {
	for _, repo := range repositories {
		for k, v := range workflows[repo] {
			r := strings.Split(repo, "/")

			for {
				resp, _, err := client.Actions.GetWorkflowUsageByID(context.Background(), r[0], r[1], k)
				if rlerr, ok := err.(*github.RateLimitError); ok {
					Logger.Info().Err(rlerr).Str("event", "get_billable_from_github").Msg("Rate limit error, waiting for reset until " + rlerr.Rate.Reset.String())
					time.Sleep(time.Until(rlerr.Rate.Reset.Time))
					continue
				} else if err != nil {
					Logger.Error().Err(err).Str("event", "get_billable_from_github").Str("repo", repo).Str("workflow", *v.Name).Str("node_id", *v.NodeID).Msg("Error to get billable from github")
					break
				}
				for os, value := range *resp.Billable {
					WorkflowBillGauge.WithLabelValues(repo, strconv.FormatInt(*v.ID, 10), *v.NodeID, *v.Name, *v.State, os).Set(float64(*value.TotalMS) / 1000)
				}
				break
			}
		}
	}
}
