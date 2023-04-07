package prom

import "github.com/google/go-github/v51/github"

var (
	repositories []string
	workflows    map[string]map[int64]github.Workflow
)
