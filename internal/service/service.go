package service

import (
	"github.com/go-cinch/common/worker"
	"github.com/google/wire"
	"gnboot/api/gnboot"
	"gnboot/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGnbootService)

// GnbootService is a gnboot service.
type GnbootService struct {
	gnboot.UnimplementedGnbootServer

	task  *worker.Worker
	movie *biz.MovieUseCase
}

// NewGnbootService new a service.
func NewGnbootService(task *worker.Worker, movie *biz.MovieUseCase) *GnbootService {
	return &GnbootService{task: task, movie: movie}
}