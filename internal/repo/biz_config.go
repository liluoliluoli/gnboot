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
	config := &ConfigRepo{
		client: client,
	}
	config.InitConfig(context.Background())
	return config
}

func (s *ConfigRepo) GetConfigBySubKey(ctx context.Context, key string, subKey string) (string, error) {
	return constant.ConfigMap[key][subKey], nil
}

func (s *ConfigRepo) GetConfigMapByKey(ctx context.Context, key string) (map[string]string, error) {
	return constant.ConfigMap[key], nil
}

func (s *ConfigRepo) InitConfig(ctx context.Context) error {
	configMap, err := json_util.Unmarshal[map[string]map[string]string](gerror.HandleRedisStringNotFound(s.client.Get(ctx, constant.RK_Configs).Val()))
	if err != nil {
		return err
	}
	constant.ConfigMap = configMap
	return nil
}
