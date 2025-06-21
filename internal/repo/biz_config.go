package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
	"github.com/liluoliluoli/gnboot/internal/common/utils/json_util"
	"github.com/redis/go-redis/v9"
)

type ConfigRepo struct {
	client redis.UniversalClient
}

func NewConfigRepo(client redis.UniversalClient) *ConfigRepo {
	return &ConfigRepo{
		client: client,
	}
}

func (s *ConfigRepo) GetConfigBySubKey(ctx context.Context, key string, subKey string) (string, error) {
	configMap, err := json_util.Unmarshal[map[string]map[string]string](gerror.HandleRedisStringNotFound(s.client.Get(ctx, constant.RK_Configs).Val()))
	if err != nil {
		return "", err
	}
	return configMap[key][subKey], nil
}
