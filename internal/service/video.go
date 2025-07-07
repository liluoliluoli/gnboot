package service

import (
	"context"
	"fmt"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/utils/array_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/cache_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/context_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
	"sort"
)

type VideoService struct {
	c                        *conf.Bootstrap
	videoRepo                *repo.VideoRepo
	actorRepo                *repo.ActorRepo
	videoActorMappingRepo    *repo.VideoActorMappingRepo
	videoSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo
	userRepo                 *repo.UserRepo
	videoUserMappingRepo     *repo.VideoUserMappingRepo
	episodeRepo              *repo.EpisodeRepo
	cache                    sdomain.Cache[*sdomain.Video]
	configRepo               *repo.ConfigRepo
}

func NewVideoService(c *conf.Bootstrap,
	videoRepo *repo.VideoRepo,
	actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo,
	videoSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo,
	userRepo *repo.UserRepo, videoUserMappingRepo *repo.VideoUserMappingRepo,
	episodeRepo *repo.EpisodeRepo, configRepo *repo.ConfigRepo) *VideoService {
	return &VideoService{
		c:                        c,
		videoRepo:                videoRepo,
		actorRepo:                actorRepo,
		videoActorMappingRepo:    videoActorMappingRepo,
		videoSubtitleMappingRepo: videoSubtitleMappingRepo,
		userRepo:                 userRepo,
		videoUserMappingRepo:     videoUserMappingRepo,
		episodeRepo:              episodeRepo,
		cache:                    repo.NewCache[*sdomain.Video](c, videoRepo.Data.Cache()),
		configRepo:               configRepo,
	}
}

func (s *VideoService) Create(ctx context.Context, item *sdomain.Video) error {
	err := gen.Use(s.videoRepo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
		err := s.cache.Flush(ctx, func(ctx context.Context) error {
			return s.videoRepo.Create(ctx, tx, item.ConvertToRepo())
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *VideoService) Get(ctx context.Context, id int64, userId int64) (*sdomain.Video, error) {
	return s.cache.Get(ctx, cache_util.GetCacheActionName(id), func(action string, ctx context.Context) (*sdomain.Video, error) {
		return s.get(ctx, id, userId)
	})
}

func (s *VideoService) get(ctx context.Context, id int64, userId int64) (*sdomain.Video, error) {
	item, err := s.videoRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	actorsMap, err := s.buildVideoActorsMap(ctx, []*sdomain.Video{item})
	if err != nil {
		return nil, err
	}
	item.Actors = actorsMap[item.ID]

	episodes, err := s.episodeRepo.QueryByVideoId(ctx, id)
	if err != nil {
		return nil, err
	}
	item.Episodes = episodes

	if userId != 0 {
		videoPlayedMap, err := s.buildVideoLastPlayedInfo(ctx, userId, []*sdomain.Video{item})
		if err != nil {
			return nil, err
		}
		if videoPlayedMap[item.ID] != nil {
			item.LastPlayedTime = videoPlayedMap[item.ID].LastPlayedTime
			item.LastPlayedEpisodeId = videoPlayedMap[item.ID].LastPlayedEpisodeId
			item.LastPlayedPosition = videoPlayedMap[item.ID].LastPlayedPosition
			item.IsFavorite = videoPlayedMap[item.ID].IsFavorite
		}
	}
	return item, nil
}

func (s *VideoService) Page(ctx context.Context, condition *sdomain.VideoSearch, userId int64) (*sdomain.PageResult[*sdomain.Video], error) {
	rp, err := s.cache.Page(ctx, cache_util.GetCacheActionName(condition), func(action string, ctx context.Context) (*sdomain.PageResult[*sdomain.Video], error) {
		return s.page(ctx, condition, userId)
	})
	if err != nil {
		return nil, err
	}
	return rp, nil
}

func (s *VideoService) page(ctx context.Context, condition *sdomain.VideoSearch, userId int64) (*sdomain.PageResult[*sdomain.Video], error) {
	var pageResult *sdomain.PageResult[*sdomain.Video]
	var err error
	if condition.IsHistory {
		pageResult, err = s.pageUserVideo(ctx, userId, nil, condition.Page)
		if err != nil {
			return nil, err
		}
	} else {
		pageResult, err = s.videoRepo.Page(ctx, condition)
		if err != nil {
			return nil, err
		}
	}

	if pageResult != nil && len(pageResult.List) != 0 {
		idToIndex := make(map[string]int)
		for i, id := range condition.Ids {
			idToIndex[fmt.Sprintf("%d", id)] = i
		}
		sort.Slice(pageResult.List, func(i, j int) bool {
			return idToIndex[fmt.Sprintf("%d", pageResult.List[i].ID)] < idToIndex[fmt.Sprintf("%d", pageResult.List[j].ID)]
		})
		jellyfinBoxIpStr, _ := s.configRepo.GetConfigBySubKey(ctx, constant.Key_BoxIpMapping, constant.SubKey_JellyfinBoxIp)
		clientIp, err := context_util.GetGenericContext[string](ctx, constant.CTX_ClientIp)
		if err != nil {
			clientIp = "127.0.0.1"
		}
		boxIp := array_util.GetHashElement(jellyfinBoxIpStr, clientIp)
		actorsMap, err := s.buildVideoActorsMap(ctx, pageResult.List)
		if err != nil {
			return nil, err
		}
		var videoPlayedMap map[int64]*sdomain.VideoUserMapping
		if userId != 0 {
			videoPlayedMap, err = s.buildVideoLastPlayedInfo(ctx, userId, pageResult.List)
			if err != nil {
				return nil, err
			}
		}
		for _, item := range pageResult.List {
			item.Actors = actorsMap[item.ID]
			if userId != 0 {
				if videoPlayedMap[item.ID] != nil {
					item.LastPlayedTime = videoPlayedMap[item.ID].LastPlayedTime
					item.LastPlayedEpisodeId = videoPlayedMap[item.ID].LastPlayedEpisodeId
					item.LastPlayedPosition = videoPlayedMap[item.ID].LastPlayedPosition
				}
			}
			if item.Thumbnail != "" && !condition.IsHistory {
				item.Thumbnail = boxIp + item.Thumbnail
			}
		}
	}
	return pageResult, nil
}

func (s *VideoService) PageFavorites(ctx context.Context, userId int64, page *sdomain.Page) (*sdomain.PageResult[*sdomain.Video], error) {
	rp, err := s.cache.Page(ctx, cache_util.GetCacheActionName(userId, page), func(action string, ctx context.Context) (*sdomain.PageResult[*sdomain.Video], error) {
		return s.pageUserVideo(ctx, userId, lo.ToPtr(true), page)
	})
	if err != nil {
		return nil, err
	}
	return rp, nil
}

func (s *VideoService) pageUserVideo(ctx context.Context, userId int64, isFavorite *bool, page *sdomain.Page) (*sdomain.PageResult[*sdomain.Video], error) {
	pageResult, err := s.videoUserMappingRepo.Page(ctx, userId, isFavorite, page)
	if err != nil {
		return nil, err
	}
	if pageResult != nil && len(pageResult.List) != 0 {
		return s.Page(ctx, &sdomain.VideoSearch{
			Ids: lo.Map(pageResult.List, func(item *sdomain.VideoUserMapping, index int) int64 {
				return item.VideoId
			}),
			Page:      page,
			IsHistory: false,
		}, userId)
	}
	return &sdomain.PageResult[*sdomain.Video]{
		Page: &sdomain.Page{
			CurrentPage: page.CurrentPage,
			PageSize:    page.PageSize,
			Count:       0,
		},
	}, nil
}

func (s *VideoService) Update(ctx context.Context, item *sdomain.UpdateVideo) error {
	err := gen.Use(s.videoRepo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
		err := s.cache.Flush(ctx, func(ctx context.Context) error {
			return s.videoRepo.Update(ctx, tx, item.ConvertToRepo())
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *VideoService) Delete(ctx context.Context, ids ...int64) error {
	err := gen.Use(s.videoRepo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
		err := s.cache.Flush(ctx, func(ctx context.Context) error {
			return s.videoRepo.Delete(ctx, tx, ids...)
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// 补齐每部影片的VideoActorMapping数组
func (s *VideoService) buildVideoActorsMap(ctx context.Context, videos []*sdomain.Video) (map[int64][]*sdomain.VideoActorMapping, error) {
	videoIds := lo.Map(videos, func(item *sdomain.Video, index int) int64 {
		return item.ID
	})
	actorMappings, err := s.videoActorMappingRepo.FindByVideoIds(ctx, videoIds)
	if err != nil {
		return nil, err
	}
	actors, err := s.actorRepo.FindByIds(ctx, lo.Map(actorMappings, func(item *sdomain.VideoActorMapping, index int) int64 {
		return item.ActorId
	}))
	if err != nil {
		return nil, err
	}
	actorsMap := lo.SliceToMap(actors, func(item *sdomain.Actor) (int64, *sdomain.Actor) {
		return item.ID, item
	})

	for _, actorMapping := range actorMappings {
		actorMapping.Actor = actorsMap[actorMapping.ActorId]
	}
	return lo.GroupBy(actorMappings, func(item *sdomain.VideoActorMapping) int64 {
		return item.VideoId
	}), nil
}

// 补齐最后播放信息
func (s *VideoService) buildVideoLastPlayedInfo(ctx context.Context, userId int64, videos []*sdomain.Video) (map[int64]*sdomain.VideoUserMapping, error) {
	videoIds := lo.Map(videos, func(item *sdomain.Video, index int) int64 {
		return item.ID
	})
	userMappings, err := s.videoUserMappingRepo.FindByUserIdAndVideoIds(ctx, userId, videoIds)
	if err != nil {
		return nil, err
	}
	return lo.SliceToMap(userMappings, func(item *sdomain.VideoUserMapping) (int64, *sdomain.VideoUserMapping) {
		return item.VideoId, item
	}), nil
}
