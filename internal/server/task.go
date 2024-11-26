package server

import (
	"context"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/worker"
	"gnboot/internal/conf"
	"gnboot/internal/service/sdomain"
	"gnboot/internal/task"
)

type Job struct {
	I4kSyncTask *task.I4kSyncTask
	Worker      *worker.Worker
}

func (j Job) Start(ctx context.Context) error {
	return nil
}

func (j Job) Stop(ctx context.Context) error {
	return nil
}

func NewJob(c *conf.Bootstrap, i4kSyncTask *task.I4kSyncTask) *Job {
	job := &Job{
		I4kSyncTask: i4kSyncTask,
	}
	job.Worker = NewWorker(c, job)
	return job
}

// New is initialize task worker from config
func NewWorker(c *conf.Bootstrap, job *Job) *worker.Worker {
	w := worker.New(
		worker.WithRedisURI(c.Data.Redis.Dsn),
		worker.WithGroup(c.Name),
		worker.WithHandler(func(ctx context.Context, p worker.Payload) error {
			switch p.UID {
			case "task1":
				job.I4kSyncTask.ProcessTest(&sdomain.Task{
					Ctx:     ctx,
					Payload: p,
				})
			}
			return nil
		}),
	)
	if w.Error != nil {
		log.Error(w.Error)
		panic(w.Error)
	}

	for id, item := range c.Task {
		err := w.Cron(
			worker.WithRunUUID(id),
			worker.WithRunGroup(item.Name),
			worker.WithRunExpr(item.Expr),
			worker.WithRunTimeout(int(item.Timeout)),
			worker.WithRunMaxRetry(int(item.Retry)),
		)
		if err != nil {
			log.Error(err)
			panic(err)
		}
	}
	log.Info("initialize worker success")
	return w
}
