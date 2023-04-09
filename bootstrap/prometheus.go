package bootstrap

import (
	"github.com/labbs/github-exporter/config"
	"github.com/labbs/github-exporter/prom"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func InitPrometheusMetrics() *prometheus.Registry {
	registry := prometheus.NewRegistry()

	registry.MustRegister(
		prom.WorkflowBillGauge,
		prom.RunnersGauge,
		prom.RunnersOrganizationGauge,
		prom.WorkflowRunStatusGauge,
		prom.WorkflowRunDurationGauge,
	)

	if !config.Metrics.DisableGoMetrics {
		registry.MustRegister(
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		)
	}

	if config.Github.EnterpriseName != "" {
		registry.MustRegister(prom.RunnersEnterpriseGauge)
		if config.Metrics.FetchEnterpriseStats {
			registry.MustRegister(prom.AdminStatsGauge)
		}
	}

	return registry
}
