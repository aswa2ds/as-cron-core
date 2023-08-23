package executor

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/aswa2ds/as-cron-db/cron_jobs"
	log "github.com/sirupsen/logrus"
)

type Executor interface {
	Start(ctx context.Context)
	AddToExecutor(cronJob cron_jobs.CronJob)
}

type executor struct {
	jobQueue []cron_jobs.CronJob
	lock     sync.Mutex
	errCh    chan<- error
}

func New(errCh chan<- error) Executor {
	return &executor{
		errCh: errCh,
	}
}

// Start implements Executor.
func (e *executor) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Info("executor stopped by stopCh")
			break
		default:
			if len(e.jobQueue) != 0 {
				e.lock.Lock()
				currJob := e.jobQueue[0]
				e.jobQueue = e.jobQueue[1:]
				e.lock.Unlock()
				err := e.do(currJob)
				if err != nil {
					e.errCh <- err
				}
			}
		}
	}
}

// AddToExecutor implements Executor.
func (e *executor) AddToExecutor(cronJob cron_jobs.CronJob) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.jobQueue = append(e.jobQueue, cronJob)
}

func (e *executor) do(cronJob cron_jobs.CronJob) error {
	switch cronJob.Protocal {
	case 0:
		_, err := http.Get(fmt.Sprintf("http://%s:%s%s", cronJob.Address, cronJob.Port, cronJob.Path))
		if err != nil {
			return err
		}
	default:
		log.Error("unsupported protocal\n")
	}
	log.Infof("Cron Job: %s has been called\n", cronJob.JobName)
	return nil
}

var _ Executor = &executor{}
