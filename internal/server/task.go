package server

import (
	"context"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/worker"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/liluoliluoli/gnboot/internal/task"
)

type Job struct {
	c           *conf.Bootstrap
	i4kSyncTask *task.I4kSyncTask
	worker      *worker.Worker
}

func (j *Job) Start(ctx context.Context) error {
	j.worker = NewWorker(j.c, j)
	return nil
}

func (j *Job) Stop(ctx context.Context) error {
	return nil
}

func NewJob(c *conf.Bootstrap, i4kSyncTask *task.I4kSyncTask) *Job {
	return &Job{
		c:           c,
		i4kSyncTask: i4kSyncTask,
	}
}

func NewWorker(c *conf.Bootstrap, job *Job) *worker.Worker {
	w := worker.New(
		worker.WithRedisURI(c.Data.Redis.Dsn),
		worker.WithGroup(c.Name),
		worker.WithHandler(func(ctx context.Context, p worker.Payload) error {
			switch p.UID {
			case "task1":
				job.i4kSyncTask.ProcessTest(&sdomain.Task{
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
