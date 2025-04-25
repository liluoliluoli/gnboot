package i4k

import (
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
)

type I4kSyncTask struct {
	c *conf.Bootstrap
}

func NewI4kSyncTask(c *conf.Bootstrap) *I4kSyncTask {
	return &I4kSyncTask{
		c: c,
	}
}

func (t *I4kSyncTask) ProcessTest(task *sdomain.Task) error {
	return nil
}
