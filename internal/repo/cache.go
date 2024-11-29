package repo

import (
	"context"
	"errors"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/liluoliluoli/gnboot/internal/utils/json_util"
	"math/rand"
	"strings"
	"time"

	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/plugins/gorm/tenant"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

// Cache .
type Cache[T any] struct {
	redis   redis.UniversalClient
	disable bool
	prefix  string
	lock    string
	val     T
	refresh bool
}

// NewCache .
func NewCache[T any](c *conf.Bootstrap, client redis.UniversalClient) sdomain.Cache[T] {
	return &Cache[T]{
		redis:   client,
		disable: c.Server.Nocache,
		lock:    "lock",
	}
}

func (c *Cache[T]) Cache() redis.UniversalClient {
	return c.redis
}

func (c *Cache[T]) WithPrefix(prefix string) sdomain.Cache[T] {
	return &Cache[T]{
		redis:   c.redis,
		disable: c.disable,
		prefix:  prefix,
		lock:    c.lock,
		val:     c.val,
	}
}

func (c *Cache[T]) WithRefresh() sdomain.Cache[T] {
	return &Cache[T]{
		redis:   c.redis,
		disable: c.disable,
		prefix:  c.prefix,
		lock:    c.lock,
		val:     c.val,
		refresh: true,
	}
}

func (c *Cache[T]) Get(
	ctx context.Context,
	action string,
	write func(string, context.Context) (T, error),
) (T, error) {
	var _res T
	if c.disable {
		return write(action, ctx)
	}
	key := c.getValKey(ctx, action)
	if !c.refresh { //无需回源
		// 1. first get cache
		rs, err := c.redis.Get(ctx, key).Result()
		if err == nil {
			// cache exists
			res, err := json_util.Unmarshal[T](rs)
			return res, err
		}
	}
	// 2. get lock before read db
	ok := c.Lock(ctx, action)
	if !ok {
		return _res, sdomain.ErrTooManyRequests(ctx)
	}
	defer c.Unlock(ctx, action)
	// 3. load repo from db and write to cache
	if write != nil {
		res, err := write(action, ctx)
		if err != nil {
			return res, err
		}
		c.Set(ctx, action, res, false)
		return res, nil
	}
	return _res, nil
}

func (c *Cache[T]) List(
	ctx context.Context,
	action string,
	write func(string, context.Context) ([]T, error),
) ([]T, error) {
	var _res []T
	if c.disable {
		return write(action, ctx)
	}
	key := c.getValKey(ctx, action)
	if !c.refresh { //无需回源
		// 1. first get cache
		rs, err := c.redis.Get(ctx, key).Result()
		if err == nil {
			// cache exists
			res, err := json_util.Unmarshal[[]T](rs)
			return res, err
		}
	}
	// 2. get lock before read db
	ok := c.Lock(ctx, action)
	if !ok {
		return _res, sdomain.ErrTooManyRequests(ctx)
	}
	defer c.Unlock(ctx, action)
	// 3. load repo from db and write to cache
	if write != nil {
		res, err := write(action, ctx)
		if err != nil {
			return res, err
		}
		c.Set(ctx, action, res, false)
		return res, nil
	}
	return _res, nil
}

func (c *Cache[T]) Page(
	ctx context.Context,
	action string,
	write func(string, context.Context) (*sdomain.PageResult[T], error),
) (*sdomain.PageResult[T], error) {
	var _res *sdomain.PageResult[T]
	if c.disable {
		return write(action, ctx)
	}
	key := c.getValKey(ctx, action)
	if !c.refresh { //无需回源
		// 1. first get cache
		rs, err := c.redis.Get(ctx, key).Result()
		if err == nil {
			// cache exists
			res, err := json_util.Unmarshal[*sdomain.PageResult[T]](rs)
			return res, err
		}
	}
	// 2. get lock before read db
	ok := c.Lock(ctx, action)
	if !ok {
		return _res, sdomain.ErrTooManyRequests(ctx)
	}
	defer c.Unlock(ctx, action)
	// 3. load repo from db and write to cache
	if write != nil {
		res, err := write(action, ctx)
		if err != nil {
			return res, err
		}
		c.Set(ctx, action, res, false)
		return res, nil
	}
	return _res, nil
}

func (c *Cache[T]) Set(ctx context.Context, action string, data any, short bool) {
	seconds := rand.New(rand.NewSource(time.Now().Unix())).Int63n(300) + 300
	if short {
		// if record not found, set a short expiration
		seconds = 60
	}
	c.SetWithExpiration(ctx, action, data, seconds)
}

func (c *Cache[T]) SetWithExpiration(ctx context.Context, action string, data any, seconds int64) {
	if c.disable {
		return
	}
	// set random expiration avoid a large number of keys expire at the same time
	err := c.redis.Set(ctx, c.getValKey(ctx, action), data, time.Duration(seconds)*time.Second).Err()
	if err != nil {
		log.
			WithContext(ctx).
			WithError(err).
			WithFields(log.Fields{
				"action":  action,
				"seconds": seconds,
			}).
			Warn("set cache failed")
		return
	}
}

func (c *Cache[T]) Del(ctx context.Context, action string) {
	if c.disable {
		return
	}
	key := c.getValKey(ctx, action)
	err := c.redis.Del(ctx, key).Err()
	if err != nil {
		log.
			WithContext(ctx).
			WithError(err).
			WithFields(log.Fields{
				"action": action,
				"key":    key,
			}).
			Warn("del cache failed")
	}
}

func (c *Cache[T]) Flush(ctx context.Context, handler func(ctx context.Context) error) (err error) {
	err = handler(ctx)
	if err != nil {
		return
	}
	if c.disable {
		return
	}
	action := c.getPrefixKey(ctx)
	arr := c.redis.Keys(ctx, action).Val()
	p := c.redis.Pipeline()
	for _, item := range arr {
		if item == c.lock {
			continue
		}
		p.Del(ctx, item)
	}
	_, pErr := p.Exec(ctx)
	if pErr != nil {
		log.
			WithContext(ctx).
			WithError(pErr).
			WithFields(log.Fields{
				"action": action,
			}).
			Warn("flush cache failed")
	}
	return
}

func (c *Cache[T]) FlushByPrefix(ctx context.Context, prefix ...string) (err error) {
	action := c.getPrefixKey(ctx, prefix...)
	arr := c.redis.Keys(ctx, action).Val()
	p := c.redis.Pipeline()
	for _, item := range arr {
		if item == c.lock {
			continue
		}
		p.Del(ctx, item)
	}
	_, pErr := p.Exec(ctx)
	if pErr != nil {
		log.
			WithContext(ctx).
			WithError(pErr).
			WithFields(log.Fields{
				"action": action,
			}).
			Warn("flush cache by prefix failed")
	}
	return
}

func (c *Cache[T]) Lock(ctx context.Context, action string) (ok bool) {
	if c.disable {
		ok = true
		return
	}
	retry := 0
	var e error
	for retry < 600 && !ok {
		ok, e = c.redis.SetNX(ctx, c.getLockKey(ctx, action), 1, time.Minute).Result()
		if errors.Is(e, context.DeadlineExceeded) ||
			errors.Is(e, context.Canceled) ||
			(e != nil && e.Error() == "gredis: connection pool timeout") {
			log.
				WithContext(ctx).
				WithError(e).
				WithFields(log.Fields{
					"action": action,
				}).
				Warn("lock failed")
			return
		}
		time.Sleep(25 * time.Millisecond)
		retry++
	}
	return
}

func (c *Cache[T]) Unlock(ctx context.Context, action string) {
	if c.disable {
		return
	}
	// get span and create new ctx since current ctx maybe timeout, unlock must be execution
	span := trace.SpanFromContext(ctx)
	ctx = trace.ContextWithSpan(context.Background(), span)
	err := c.redis.Del(ctx, c.getLockKey(ctx, action)).Err()
	if err != nil {
		log.
			WithContext(ctx).
			WithError(err).
			WithFields(log.Fields{
				"action": action,
			}).
			Warn("unlock cache failed")
	}
}

func (c *Cache[T]) getPrefixKey(ctx context.Context, arr ...string) string {
	id := tenant.FromContext(ctx)
	prefix := c.prefix
	if len(arr) > 0 {
		// append params prefix need add val
		prefix = strings.Join(append([]string{prefix}, arr...), "_")
	}
	if strings.TrimSpace(prefix) == "" {
		// avoid flush all key
		log.
			WithContext(ctx).
			Warn("invalid prefix")
		prefix = "prefix"
	}
	if id == "" {
		return strings.Join([]string{prefix, "*"}, "")
	}
	return strings.Join([]string{id, "_", prefix, "*"}, "")
}

func (c *Cache[T]) getValKey(ctx context.Context, action string) string {
	id := tenant.FromContext(ctx)
	if id == "" {
		return strings.Join([]string{c.prefix, action}, "_")
	}
	return strings.Join([]string{id, c.prefix, action}, "_")
}

func (c *Cache[T]) getLockKey(ctx context.Context, action string) string {
	id := tenant.FromContext(ctx)
	if id == "" {
		return strings.Join([]string{c.prefix, c.lock, action}, "_")
	}
	return strings.Join([]string{id, c.prefix, c.lock, action}, "_")
}
