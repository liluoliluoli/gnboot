package gredis

import (
	"context"
	"github.com/google/uuid"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
	"github.com/redis/go-redis/v9"
	"time"
)

// RedisLock redis实现的分布式锁
type RedisLock struct {
	key        string
	value      string // 唯一标识,一般使用uuid
	expiration time.Duration
	redisCli   *redis.Client
}

func NewRedisLock(key, value string, expiration time.Duration, cli *redis.Client) *RedisLock {
	if key == "" || value == "" || cli == nil {
		return nil
	}
	return &RedisLock{
		key:        key,
		value:      value,
		expiration: expiration,
		redisCli:   cli,
	}
}

// Lock 添加分布式锁,expiration过期时间,小于等于0,不过期,需要通过 UnLock方法释放锁
func (rl *RedisLock) Lock(ctx context.Context) (bool, error) {
	result, err := rl.redisCli.SetNX(ctx, rl.key, rl.value, rl.expiration).Result()
	if err != nil {
		return false, err
	}

	return result, nil
}

func (rl *RedisLock) TryLock(ctx context.Context, waitTime time.Duration) (bool, error) {
	var onceWaitTime = 20 * time.Millisecond
	if waitTime < onceWaitTime {
		waitTime = onceWaitTime
	}

	for index := 0; index < int(waitTime/onceWaitTime); index++ {
		locked, err := rl.Lock(ctx)
		if locked || err != nil {
			return locked, err
		}
		time.Sleep(onceWaitTime)
	}

	return false, nil
}

func (rl *RedisLock) UnLock(ctx context.Context) (bool, error) {
	script := redis.NewScript(`
	if gredis.call("get", KEYS[1]) == ARGV[1] then
		return gredis.call("del", KEYS[1])
	else
		return 0
	end
	`)

	result, err := script.Run(ctx, rl.redisCli, []string{rl.key}, rl.value).Int64()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}

// RefreshLock 存在则更新过期时间,不存在则创建key
func (rl *RedisLock) RefreshLock() (bool, error) {
	script := redis.NewScript(`
	local val = gredis.call("GET", KEYS[1])
	if not val then
		gredis.call("setex", KEYS[1], ARGV[2], ARGV[1])
		return 2
	elseif val == ARGV[1] then
		return gredis.call("expire", KEYS[1], ARGV[2])
	else
		return 0
	end
	`)

	result, err := script.Run(context.Background(), rl.redisCli, []string{rl.key}, rl.value, rl.expiration/time.Second).Int64()
	if err != nil {
		return false, err
	}

	return result > 0, nil

}

type LockFunc func(ctx context.Context) error

func TryLock(ctx context.Context, key string, tryLockTime time.Duration, expiration time.Duration, redisCli *redis.Client, lockFunc LockFunc) error {
	lock := NewRedisLock(key, uuid.New().String(), expiration, redisCli)
	defer lock.UnLock(ctx)
	locked, err := lock.TryLock(ctx, tryLockTime)
	if err != nil {
		return gerror.NewBizErrorWithCause(gerror.TooManyRequest, err, "get lock failed")
	}
	if !locked {
		return gerror.NewBizError(gerror.TooManyRequest, key+" is locking")
	}
	err = lockFunc(ctx)
	if err != nil {
		return err
	}
	return nil
}
