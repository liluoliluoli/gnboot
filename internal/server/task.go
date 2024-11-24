package server

import (
	"context"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/worker"
	"github.com/pkg/errors"
	"gnboot/internal/adaptor/task"
	"gnboot/internal/conf"
	"gnboot/internal/service/sdomain"
)

// New is initialize task worker from config
func NewWorker(c *conf.Bootstrap) (w *worker.Worker, err error) {
	w = worker.New(
		worker.WithRedisURI(c.Data.Redis.Dsn),
		worker.WithGroup(c.Name),
		worker.WithHandler(func(ctx context.Context, p worker.Payload) error {
			return process(&sdomain.Task{
				Ctx:     ctx,
				Payload: p,
			})
		}),
	)
	if w.Error != nil {
		log.Error(w.Error)
		err = errors.New("initialize worker failed")
		return
	}

	for id, item := range c.Task {
		err = w.Cron(
			worker.WithRunUUID(id),
			worker.WithRunGroup(item.Name),
			worker.WithRunExpr(item.Expr),
			worker.WithRunTimeout(int(item.Timeout)),
			worker.WithRunMaxRetry(int(item.Retry)),
		)
		if err != nil {
			log.Error(err)
			err = errors.New("initialize worker failed")
			return
		}
	}

	log.Info("initialize worker success")
	return
}

func process(t *sdomain.Task) (err error) {
	switch t.Payload.UID {
	case "task1":
		task.ProcessTest(t)
	case "task2":
		log.WithContext(t.Ctx).Info("task2: %s", t.Payload.Payload)
	}
	return
}
