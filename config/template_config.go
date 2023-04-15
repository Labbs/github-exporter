package config

import (
	"os"
)

var templateConfigFile string = `{
  "port": 8080,
  "debug": false,
  "cache_http_size_bytes": 104857600,

  "github": {
    "token": "######",
    "refresh_interval": 30,
    "organizations": ["mycompany"],
    "repositories": ["mycompany/test"],
    "application_id": 123456,
    "application_installation_id": 123456,
    "application_private_key": "#######",
    "enterprise_url": "https://github.mycompany.com",
    "enterprise_name": "mycompany"
  },

  "metrics": {
    "disable_go_metrics": true,
    "fetch_workflow_usage": true,
    "fetch_enterprise_stats": true,
    "workflow_fields": ["repo","id","..."]
  }
}`

func GenerateTemplateConfigFile(path string) {
	if path == "" {
		path = "config.json"
	}
	err := os.WriteFile(path, []byte(templateConfigFile), 0644)
	if err != nil {
		panic(err)
	}
}
