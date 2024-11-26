package task

import (
	"gnboot/internal/conf"
	"gnboot/internal/service"
	"gnboot/internal/service/sdomain"
)

type I4kSyncTask struct {
	c            *conf.Bootstrap
	movieService *service.MovieService
}

func NewI4kSyncTask(c *conf.Bootstrap, movieService *service.MovieService) *I4kSyncTask {
	return &I4kSyncTask{
		c:            c,
		movieService: movieService,
	}
}

func (t *I4kSyncTask) ProcessTest(task *sdomain.Task) error {
	_, err := t.movieService.Get(task.Ctx, 1)
	if err != nil {
		return err
	}
	return nil
}
