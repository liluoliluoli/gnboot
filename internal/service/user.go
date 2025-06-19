package service

import (
	"context"
	"errors"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/utils/cache_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/context_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	c                    *conf.Bootstrap
	userRepo             *repo.UserRepo
	videoUserMappingRepo *repo.VideoUserMappingRepo
	cache                sdomain.Cache[*sdomain.User]
	client               redis.UniversalClient
}

func NewUserService(c *conf.Bootstrap,
	userRepo *repo.UserRepo,
	videoUserMappingRepo *repo.VideoUserMappingRepo,
	client redis.UniversalClient,
) *UserService {
	return &UserService{
		c:                    c,
		userRepo:             userRepo,
		videoUserMappingRepo: videoUserMappingRepo,
		cache:                repo.NewCache[*sdomain.User](c, userRepo.Data.Cache()),
		client:               client,
	}
}

func (s *UserService) Get(ctx context.Context, id int64) (*sdomain.User, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.User, error) {
		return s.get(ctx, id)
	})
}

func (s *UserService) get(ctx context.Context, id int64) (*sdomain.User, error) {
	item, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *UserService) QueryByUserName(ctx context.Context, userName string) (*sdomain.User, error) {
	item, err := s.userRepo.GetByUserName(ctx, userName)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *UserService) UpdateSessionToken(ctx context.Context, user *sdomain.User) error {
	err := gen.Use(s.videoUserMappingRepo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
		err := s.cache.Flush(ctx, func(ctx context.Context) error {
			return s.userRepo.Update(ctx, tx, user)
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *UserService) UpdateFavorite(ctx context.Context, userId int64, videoId int64, favorite bool) error {
	err := gen.Use(s.videoUserMappingRepo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
		err := s.cache.Flush(ctx, func(ctx context.Context) error {
			return s.videoUserMappingRepo.UpdateFavorite(ctx, tx, userId, videoId, favorite)
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *UserService) UpdatePlayStatus(ctx context.Context, userId int64, videoId int64, episodeId int64, position int64, playTimestamp int64) error {
	err := gen.Use(s.videoUserMappingRepo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
		err := s.cache.Flush(ctx, func(ctx context.Context) error {
			return s.videoUserMappingRepo.UpdatePlayStatus(ctx, tx, userId, videoId, episodeId, position, playTimestamp)
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *UserService) Create(ctx context.Context, userName string, password string) error {
	err := gen.Use(s.videoUserMappingRepo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
		err := s.cache.Flush(ctx, func(ctx context.Context) error {
			return s.userRepo.Create(ctx, tx, userName, password, constant.Trial)
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *UserService) GetCurrentUser(ctx context.Context) (*sdomain.User, error) {
	userName, err := context_util.GetGenericContext[string](ctx, constant.CTX_UserName)
	if err != nil {
		return nil, err
	}
	rs, err := s.QueryByUserName(ctx, userName)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, errors.New("用户不存在")
	}
	return rs, nil
}

func (s *UserService) UpdateNotice(ctx context.Context, title string, content string) error {
	cmd := s.client.HSet(ctx, constant.RK_Notice, constant.HK_NoticeTitle, title)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	cmd = s.client.HSet(ctx, constant.RK_Notice, constant.HK_NoticeContent, content)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
