package config

import "github.com/urfave/cli/v2"

var (
	Port       int
	Version    string
	Debug      bool
	ConfigFile string

	Github struct {
		// Github refresh interval in seconds
		RefreshInterval int
		// Github Personal Access Token
		Token string

		// Github App
		ApplicationID             int64
		ApplicationInstallationId int64
		ApplicationPrivateKey     string

		// Github Enterprise
		EnterpriseURL  string
		EnterpriseName string

		// Github Organization
		Organizations cli.StringSlice
		// Github Repository
		Repositories cli.StringSlice
	}

	Metrics struct {
		DisableGoMetrics     bool
		FetchWorkflowUsage   bool
		FetchEnterpriseStats bool
	}

	WorkflowFields cli.StringSlice
	CacheHTTPSize  int64
)
