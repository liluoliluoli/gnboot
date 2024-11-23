package service

import (
	"context"
	"gnboot/internal/repo"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewMovieUseCase,
)

type Transaction interface {
	Tx(ctx context.Context, handler func(context.Context) error) error
}

type Cache[T any] interface {
	Cache() redis.UniversalClient
	WithPrefix(prefix string) Cache[T]
	WithRefresh() Cache[T]
	Get(ctx context.Context, action string, write func(string, context.Context) (T, error)) (T, error)
	GetPage(ctx context.Context, action string, write func(string, context.Context) (repo.PageCache[T], error)) (repo.PageCache[T], error)
	Set(ctx context.Context, action string, data any, short bool)
	Del(ctx context.Context, action string)
	SetWithExpiration(ctx context.Context, action string, data any, seconds int64)
	Flush(ctx context.Context, handler func(context.Context) error) error
	FlushByPrefix(ctx context.Context, prefix ...string) (err error)
}
