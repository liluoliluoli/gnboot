package task

import (
	"gnboot/internal/service"
	"gnboot/internal/service/sdomain"
)

type I4kSyncTask struct {
	movieService *service.MovieService
}

func NewI4kSyncTask(movieService *service.MovieService) *I4kSyncTask {
	return &I4kSyncTask{
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
