package adaptor

import (
	"github.com/go-cinch/common/worker"
	"github.com/google/wire"
	"gnboot/api/movie"
	"gnboot/internal/service"
)

// ProviderSet is adaptor providers.
var ProviderSet = wire.NewSet(NewGnbootService)

// GnbootService is a gnboot adaptor.
type GnbootService struct {
	movie.UnimplementedMovieRemoteServiceServer

	task  *worker.Worker
	movie *service.MovieUseCase
}

// NewGnbootService new a adaptor.
func NewGnbootService(task *worker.Worker, movie *service.MovieUseCase) *GnbootService {
	return &GnbootService{task: task, movie: movie}
}
