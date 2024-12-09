package i4k

import (
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/service"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/liluoliluoli/gnboot/internal/task/i4k/movie"
	"strconv"
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
	_, err := t.movieService.Get(task.Ctx, 1, 1)
	//httpclient分页循环查询i4k获取电影列表

	sum := 0
	for sum <= 10 {
		movie.TaskList(strconv.Itoa(sum))
		sum += sum
	}

	//遍历电影列表，插入
	//err = t.movieService.Create(task.Ctx, &sdomain.CreateMovie{
	//	ExternalID:    "123",
	//	OriginalTitle: "sss",
	//})
	if err != nil {
		return err
	}
	return nil
}
