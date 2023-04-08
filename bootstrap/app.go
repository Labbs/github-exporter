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
	app.Logger = InitLogger(version, config.Debug)
	app.Fiber = InitFiber(app.Logger)
	app.GithubClient = NewGHClient(app.Logger)
	app.CronScheduler = InitCronScheduler()

	prom.Logger = app.Logger

	app.CronScheduler.Every(config.Github.RefreshInterval*5).Second().Do(prom.GithubFetcher, app.GithubClient)

	if config.Github.EnterpriseName != "" {
		app.CronScheduler.Every(config.Github.RefreshInterval).Second().Do(prom.GetRunnersEnterpriseFromGithub, app.GithubClient)
	}

	app.CronScheduler.Every(config.Github.RefreshInterval).Second().Do(prom.GetBillableFromGithub, app.GithubClient)
	app.CronScheduler.Every(config.Github.RefreshInterval).Second().Do(prom.GetRunnersFromGhithub, app.GithubClient)
	app.CronScheduler.Every(config.Github.RefreshInterval).Second().Do(prom.GetRunnersOrganizationFromGithub, app.GithubClient)
	app.CronScheduler.Every(config.Github.RefreshInterval).Second().Do(prom.GetWorkflowRunsFromGithub, app.GithubClient)

	return *app
}
