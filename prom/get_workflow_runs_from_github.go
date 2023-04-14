package prom

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v51/github"
	"github.com/labbs/github-exporter/config"
)

func GetWorkflowRunsFromGithub(client *github.Client) {
	window_start := time.Now().Add(time.Duration(-12) * time.Hour).Format(time.RFC3339)

	for _, repo := range repositories {
		r := strings.Split(repo, "/")
		var runs []*github.WorkflowRun
		opt := &github.ListWorkflowRunsOptions{
			ListOptions: github.ListOptions{PerPage: 200},
			Created:     ">=" + window_start,
		}

		for {
			workflowRuns, resp, err := client.Actions.ListRepositoryWorkflowRuns(context.Background(), r[0], r[1], opt)
			if rlerr, ok := err.(*github.RateLimitError); ok {
				Logger.Info().Err(rlerr).Str("event", "get_workflow_runs_from_github").Msg("Rate limit error, waiting for reset until " + rlerr.Rate.Reset.String())
				time.Sleep(time.Until(rlerr.Rate.Reset.Time))
				continue
			} else if err != nil {
				Logger.Error().Err(err).Str("event", "get_workflow_runs_from_github").Str("repository", repo).Msg("Error to get workflow runs from github")
				break
			}
			runs = append(runs, workflowRuns.WorkflowRuns...)
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}

		for _, run := range runs {
			var status float64 = 0
			if run.GetStatus() == "success" {
				status = 1
			} else if run.GetStatus() == "skipped" {
				status = 2
			} else if run.GetStatus() == "in_progress" {
				status = 3
			} else if run.GetStatus() == "queued" {
				status = 4
			}

			fields := getRelevantFields(repo, run)
			WorkflowRunStatusGauge.WithLabelValues(fields...).Set(status)

			var run_usage *github.WorkflowRunUsage = nil
			if config.Metrics.FetchWorkflowUsage {
				run_usage = getRunUsage(client, r[0], r[1], *run.ID)
			}

			if run_usage != nil {
				WorkflowRunDurationGauge.WithLabelValues(fields...).Set(float64(run_usage.GetRunDurationMS()))
			} else {
				created := run.CreatedAt.Time.Unix()
				updated := run.UpdatedAt.Time.Unix()
				elapsed := updated - created
				WorkflowRunDurationGauge.WithLabelValues(fields...).Set(float64(elapsed))
			}
		}
	}
}

func getFieldValue(run github.WorkflowRun, repo, field string) string {
	switch field {
	case "repo":
		return repo
	case "id":
		return strconv.FormatInt(*run.ID, 10)
	case "node_id":
		return *run.NodeID
	case "head_branch":
		return *run.HeadBranch
	case "head_sha":
		return *run.HeadSHA
	case "run_number":
		return strconv.Itoa(*run.RunNumber)
	case "workflow_id":
		return strconv.FormatInt(*run.WorkflowID, 10)
	case "workflow":
		r, exist := workflows[repo]
		if !exist {
			Logger.Error().Str("event", "get_workflow_runs_from_github").Str("repository", repo).Msg("Couldn't fetch repo from workflow cache.")
			return "unknown"
		}
		w, exist := r[*run.WorkflowID]
		if !exist {
			Logger.Error().Str("event", "get_workflow_runs_from_github").Str("repository", repo).Int64("workflow", *run.WorkflowID).Msg("Couldn't fetch repo from workflow cache.")
			return "unknown"
		}
		return *w.Name
	case "event":
		return *run.Event
	case "status":
		return *run.Status
	}
	Logger.Error().Str("event", "get_workflow_runs_from_github").Str("repository", repo).Str("field", field).Msg("Tried to fetch invalid field.")
	return ""
}

func getRelevantFields(repo string, run *github.WorkflowRun) []string {
	result := make([]string, len(config.Metrics.WorkflowFields.Value()))
	for i, field := range config.Metrics.WorkflowFields.Value() {
		result[i] = getFieldValue(*run, repo, field)
	}
	return result
}

func getRunUsage(client *github.Client, owner string, repo string, runId int64) *github.WorkflowRunUsage {
	for {
		resp, _, err := client.Actions.GetWorkflowRunUsageByID(context.Background(), owner, repo, runId)
		if rl_err, ok := err.(*github.RateLimitError); ok {
			Logger.Info().Err(rl_err).Str("event", "get_run_usage").Msg("Rate limit error, waiting for reset until " + rl_err.Rate.Reset.String())
			time.Sleep(time.Until(rl_err.Rate.Reset.Time))
			continue
		} else if err != nil {
			Logger.Error().Err(err).Str("event", "get_run_usage").Str("repository", repo).Int64("run_id", runId).Msg("Error to get workflow run usage from github")
			return nil
		}
		return resp
	}
}
