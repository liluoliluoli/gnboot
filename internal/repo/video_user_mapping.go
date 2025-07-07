package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
	"time"
)

type VideoUserMappingRepo struct {
	Data *Data
}

func NewVideoUserMappingRepo(data *Data) *VideoUserMappingRepo {
	return &VideoUserMappingRepo{
		Data: data,
	}
}

func (r *VideoUserMappingRepo) do(ctx context.Context, tx *gen.Query) gen.IVideoUserMappingDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).VideoUserMapping.WithContext(ctx)
	} else {
		return tx.VideoUserMapping.WithContext(ctx)
	}
}

func (r *VideoUserMappingRepo) FindByUserIdAndVideoIds(ctx context.Context, userId int64, videoIds []int64) ([]*sdomain.VideoUserMapping, error) {
	do := r.do(ctx, nil)
	if userId != 0 {
		do = do.Where(gen.VideoUserMapping.UserID.Eq(userId))
	}
	if len(videoIds) > 0 {
		do = do.Where(gen.VideoUserMapping.VideoID.In(videoIds...))
	}
	finds, err := do.Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return lo.Map(finds, func(item *model.VideoUserMapping, index int) *sdomain.VideoUserMapping {
		return (&sdomain.VideoUserMapping{}).ConvertFromRepo(item)
	}), nil
}

func (r *VideoUserMappingRepo) UpdateFavorite(ctx context.Context, tx *gen.Query, userId int64, videoId int64, favorite bool) error {
	first, err := r.do(ctx, tx).Where(gen.VideoUserMapping.UserID.Eq(userId)).Where(gen.VideoUserMapping.VideoID.Eq(videoId)).First()
	if handleQueryError(ctx, err) != nil {
		return err
	}

	if first != nil {
		first.IsFavorite = favorite
		first.UpdateTime = time.Now()
		_, err := r.do(ctx, tx).Select(gen.VideoUserMapping.IsFavorite, gen.VideoUserMapping.UpdateTime).Updates(first)
		if err != nil {
			return err
		}
	} else {
		err = r.do(ctx, tx).Create(&model.VideoUserMapping{
			VideoID:    videoId,
			UserID:     userId,
			IsFavorite: favorite,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *VideoUserMappingRepo) UpdatePlayStatus(ctx context.Context, tx *gen.Query, userId int64, videoId int64, episodeId int64, position int64, playTimestamp int64) error {
	first, err := r.do(ctx, tx).Where(gen.VideoUserMapping.UserID.Eq(userId)).Where(gen.VideoUserMapping.VideoID.Eq(videoId)).First()
	if handleQueryError(ctx, err) != nil {
		return err
	}
	if first != nil {
		first.LastPlayedPosition = lo.ToPtr(position)
		first.LastPlayedTime = lo.ToPtr(time.Unix(playTimestamp, 0))
		first.LastPlayedEpisodeID = lo.ToPtr(episodeId)
		first.UpdateTime = time.Now()
		_, err := r.do(ctx, tx).Updates(first)
		if err != nil {
			return err
		}
	} else {
		err = r.do(ctx, tx).Create(&model.VideoUserMapping{
			VideoID:             videoId,
			UserID:              userId,
			IsFavorite:          false,
			LastPlayedEpisodeID: lo.ToPtr(episodeId),
			LastPlayedPosition:  lo.ToPtr(position),
			LastPlayedTime:      lo.ToPtr(time.Unix(playTimestamp, 0)),
			CreateTime:          time.Now(),
			UpdateTime:          time.Now(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *VideoUserMappingRepo) Page(ctx context.Context, userId int64, isFavorite *bool, page *sdomain.Page) (*sdomain.PageResult[*sdomain.VideoUserMapping], error) {
	do := r.do(ctx, nil).Where(gen.VideoUserMapping.UserID.Eq(userId))
	if isFavorite != nil {
		do = do.Where(gen.VideoUserMapping.IsFavorite.Value(lo.FromPtr(isFavorite)))
	}
	do = do.Order(gen.VideoUserMapping.UpdateTime.Desc()).Order(gen.VideoUserMapping.LastPlayedTime.Desc())
	list, total, err := do.FindByPage(int((page.CurrentPage-1)*page.PageSize), int(page.PageSize))
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return &sdomain.PageResult[*sdomain.VideoUserMapping]{
		Page: &sdomain.Page{
			CurrentPage: page.CurrentPage,
			PageSize:    page.PageSize,
			Count:       total,
		},
		List: lo.Map(list, func(item *model.VideoUserMapping, index int) *sdomain.VideoUserMapping {
			return (&sdomain.VideoUserMapping{}).ConvertFromRepo(item)
		}),
	}, nil
}
