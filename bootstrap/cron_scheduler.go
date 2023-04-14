package bootstrap

import (
	"time"

	"github.com/go-co-op/gocron"
)

func InitCronScheduler() *gocron.Scheduler {
	cronScheduler := gocron.NewScheduler(time.UTC)
	cronScheduler.TagsUnique()
	cronScheduler.StartAsync()
	return cronScheduler
}
