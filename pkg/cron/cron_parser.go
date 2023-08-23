package cron

import (
	"time"

	cron "github.com/robfig/cron/v3"
)

func NextToggleTime(cronExpression string, currToggleTime time.Time) (time.Time, error) {
	standardParser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)
	schedule, err := standardParser.Parse(cronExpression)
	if err != nil {
		return time.Time{}, err
	}

	nextToggleTime := schedule.Next(currToggleTime)

	return nextToggleTime, nil
}
