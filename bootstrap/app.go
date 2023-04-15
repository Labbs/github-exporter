package bootstrap

import (
	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/v51/github"
	"github.com/labbs/github-exporter/config"
	"github.com/labbs/github-exporter/prom"
	"github.com/rs/zerolog"
)

type Application struct {
	Logger        zerolog.Logger
	Fiber         *fiber.App
	GithubClient  *github.Client
	CronScheduler *gocron.Scheduler
}

func App(version string) Application {
	app := &Application{}
	prom.Fields = config.Metrics.WorkflowFields.Value()
	app.Logger = InitLogger(version, config.Debug)
	app.Fiber = InitFiber(app.Logger)
	app.GithubClient = NewGHClient(app.Logger)
	app.CronScheduler = InitCronScheduler()

	prom.Logger = app.Logger

	_, err := app.CronScheduler.Every(config.Github.RefreshInterval*5).Second().Do(prom.GithubFetcher, app.GithubClient)
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Error scheduling GtihubFetcher jobs")
	}

	if config.Github.EnterpriseName != "" {
		_, err = app.CronScheduler.Every(config.Github.RefreshInterval).Second().Do(prom.GetRunnersEnterpriseFromGithub, app.GithubClient)
		if err != nil {
			app.Logger.Fatal().Err(err).Msg("Error scheduling GetRunnersEnterpriseFromGithub jobs")
		}
	}

	_, err = app.CronScheduler.Every(config.Github.RefreshInterval).Second().Do(prom.GetBillableFromGithub, app.GithubClient)
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Error scheduling GetBillableFromGithub jobs")
	}
	_, err = app.CronScheduler.Every(config.Github.RefreshInterval).Second().Do(prom.GetRunnersFromGhithub, app.GithubClient)
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Error scheduling GetRunnersFromGhithub jobs")
	}
	_, err = app.CronScheduler.Every(config.Github.RefreshInterval).Second().Do(prom.GetRunnersOrganizationFromGithub, app.GithubClient)
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Error scheduling GetRunnersOrganizationFromGithub jobs")
	}
	_, err = app.CronScheduler.Every(config.Github.RefreshInterval).Second().Do(prom.GetWorkflowRunsFromGithub, app.GithubClient)
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Error scheduling GetWorkflowRunsFromGithub jobs")
	}

	return *app
}
