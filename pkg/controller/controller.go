package controller

import (
	"context"
	"time"

	"github.com/aswa2ds/as-cron-core/pkg/cron"
	"github.com/aswa2ds/as-cron-core/pkg/executor"
	db "github.com/aswa2ds/as-cron-db"
	"github.com/aswa2ds/as-cron-db/cron_jobs"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Controller interface {
	Start(ctx context.Context)
}

type controller struct {
	executor executor.Executor
	errCh    chan error
}

func New() Controller {
	errCh := make(chan error)
	return &controller{
		errCh:    errCh,
		executor: executor.New(errCh),
	}
}

// Start implements Controller.
func (c *controller) Start(ctx context.Context) {
	go c.executor.Start(ctx)
	for {
		select {
		case <-ctx.Done():
			log.Info("controller stopped by cancel()")
			break
		case err := <-c.errCh:
			log.Error(err)
		default:
			c.syncCronJobs()
		}
		time.Sleep(time.Second)
	}
}

func (c *controller) syncCronJobs() {
	client, err := db.ClientSet()
	if err != nil {
		c.errCh <- err
	}
	cronJobs := client.CronJobs().List()
	currentTime := time.Now()
	for _, cronJob := range cronJobs {
		duration := cronJob.NextToggleTime.Sub(currentTime)
		if duration.Seconds() < 0 {
			log.Infof("Cron job: %s, should be scheduled at %s, but not\n", cronJob.JobName, cronJob.NextToggleTime.String())
			continue
		}
		if duration.Seconds() <= 5.0 {
			c.addToExecutor(cronJob)
			nextToggleTime, err := cron.NextToggleTime(cronJob.CronExpression, cronJob.NextToggleTime)
			if err != nil {
				c.errCh <- err
			}
			cronJobUpdate := cron_jobs.CronJob{
				Model: gorm.Model{
					ID: cronJob.ID,
				},
				NextToggleTime: nextToggleTime,
			}

			client.CronJobs().UpdateNextToggleTime(cronJobUpdate)
		}
	}
}

func (c *controller) addToExecutor(cronJob cron_jobs.CronJob) {
	c.executor.AddToExecutor(cronJob)
}

var _ Controller = &controller{}
