package prom

import (
	"github.com/labbs/github-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	WorkflowBillGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_workflow_usage_seconds",
			Help: "Number of billable seconds used by a specific workflow during the current billing cycle. Any job re-runs are also included in the usage. Only apply to workflows in private repositories that use GitHub-hosted runners.",
		},
		[]string{"repo", "id", "node_id", "name", "state", "os"},
	)

	RunnersEnterpriseGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_enterprise_status",
			Help: "runner status",
		},
		[]string{"os", "name", "id"},
	)

	RunnersGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_status",
			Help: "runner status",
		},
		[]string{"repo", "os", "name", "id", "busy"},
	)

	RunnersOrganizationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_organization_status",
			Help: "runner status",
		},
		[]string{"organization", "os", "name", "id", "busy"},
	)

	WorkflowRunStatusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_workflow_run_status",
			Help: "Workflow run status of all workflow runs created in the last 12hr",
		},
		config.WorkflowFields.Value(),
	)

	WorkflowRunDurationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_workflow_run_duration_ms",
			Help: "Workflow run duration (in milliseconds) of all workflow runs created in the last 12hr",
		},
		config.WorkflowFields.Value(),
	)

	AdminStatsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_admin_stats",
			Help: "Admin stats for a GitHub Enterprise instance",
		},
		[]string{"name", "type"},
	)
)
