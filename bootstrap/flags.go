package bootstrap

import (
	"github.com/labbs/github-exporter/config"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func GenericFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			EnvVars:     []string{"CONFIG"},
			Usage:       "Config file path",
			Value:       "config.json",
			Destination: &config.ConfigFile,
		},
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"d"},
			EnvVars:     []string{"DEBUG"},
			Value:       false,
			Usage:       "Enable debug mode",
			Destination: &config.Debug,
		}),
	}
}

func GithubFlags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "github.token",
			Usage:       "Github Token",
			Aliases:     []string{"gt"},
			Destination: &config.Github.Token,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "github.refresh_interval",
			Usage:       "Refresh interval in seconds",
			Aliases:     []string{"gri"},
			Destination: &config.Github.RefreshInterval,
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:        "github.organizations",
			Usage:       "Github Organizations",
			Aliases:     []string{"go"},
			Destination: &config.Github.Organizations,
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:        "github.repositories",
			Usage:       "Github Repositories",
			Aliases:     []string{"gr"},
			Destination: &config.Github.Repositories,
		}),
		altsrc.NewInt64Flag(&cli.Int64Flag{
			Name:        "github.application_id",
			Usage:       "Github Application ID",
			Aliases:     []string{"gai"},
			Destination: &config.Github.ApplicationID,
		}),
		altsrc.NewInt64Flag(&cli.Int64Flag{
			Name:        "github.application_installation_id",
			Usage:       "Github Application Installation ID",
			Aliases:     []string{"gaii"},
			Destination: &config.Github.ApplicationInstallationId,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "github.application_private_key",
			Usage:       "Github Application Private Key",
			Aliases:     []string{"gapk"},
			Destination: &config.Github.ApplicationPrivateKey,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "github.enterprise_url",
			Usage:       "Github Enterprise URL",
			Aliases:     []string{"geu"},
			Destination: &config.Github.EnterpriseURL,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "github.enterprise_name",
			Usage:       "Github Enterprise Name",
			Aliases:     []string{"gen"},
			Destination: &config.Github.EnterpriseName,
		}),
	}
}

func MetricsFlags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "metrics.disable_go_metrics",
			Usage:       "Disable Go Metrics",
			Aliases:     []string{"mdgm"},
			Value:       true,
			Destination: &config.Metrics.DisableGoMetrics,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "metrics.fetch_workflow_usage",
			Usage:       "Fetch Workflow Usage",
			Aliases:     []string{"mfwu"},
			Value:       true,
			Destination: &config.Metrics.FetchWorkflowUsage,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "metrics.fetch_enterprise_stats",
			Usage:       "Fetch Enterprise Stats",
			Aliases:     []string{"mfes"},
			Value:       true,
			Destination: &config.Metrics.FetchEnterpriseStats,
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:        "workflow_fields",
			Usage:       "Workflow fields to export",
			Aliases:     []string{"wf"},
			Value:       cli.NewStringSlice("repo", "id", "node_id", "head_branch", "head_sha", "run_number", "workflow_id", "workflow", "event", "status"),
			Destination: &config.Metrics.WorkflowFields,
		}),
	}
}

func ServerFlags() []cli.Flag {
	flags := []cli.Flag{
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "port",
			Aliases:     []string{"p"},
			EnvVars:     []string{"PORT"},
			Usage:       "Server Port",
			Destination: &config.Port,
		}),
		altsrc.NewInt64Flag(&cli.Int64Flag{
			Name:        "cache_http_size_bytes",
			Usage:       "Cache HTTP Size in bytes",
			Aliases:     []string{"chs"},
			Value:       100 * 1024 * 1024,
			Destination: &config.CacheHTTPSize,
		}),
	}
	flags = append(flags, GithubFlags()...)
	flags = append(flags, MetricsFlags()...)
	flags = append(flags, GenericFlags()...)
	return flags
}
