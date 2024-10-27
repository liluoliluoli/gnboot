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

	task   *worker.Worker
	gnboot *biz.GnbootUseCase
}

// NewGnbootService new a service.
func NewGnbootService(task *worker.Worker, gnboot *biz.GnbootUseCase) *GnbootService {
	return &GnbootService{task: task, gnboot: gnboot}
}
