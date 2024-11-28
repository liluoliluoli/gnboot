package sdomain

import (
	"context"
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/worker"
	"github.com/redis/go-redis/v9"
)

type Cache[T any] interface {
	Cache() redis.UniversalClient
	WithPrefix(prefix string) Cache[T]
	WithRefresh() Cache[T]
	Get(ctx context.Context, action string, write func(string, context.Context) (T, error)) (T, error)
	GetPage(ctx context.Context, action string, write func(string, context.Context) (*PageResult[T], error)) (*PageResult[T], error)
	Set(ctx context.Context, action string, data any, short bool)
	Del(ctx context.Context, action string)
	SetWithExpiration(ctx context.Context, action string, data any, seconds int64)
	Flush(ctx context.Context, handler func(context.Context) error) error
	FlushByPrefix(ctx context.Context, prefix ...string) (err error)
}

type Task struct {
	Ctx     context.Context
	Payload worker.Payload
}

type PageResult[T any] struct {
	Page *page.Page `json:"page"`
	List []T        `json:"list"`
}

type Sort struct {
	Filter    string `json:"filter"`
	Type      string `json:"type"`
	Direction string `json:"direction"`
}
