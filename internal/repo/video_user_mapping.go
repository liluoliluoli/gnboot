package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
	"gorm.io/gorm"
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

func (r *VideoUserMappingRepo) FindByUserIdAndVideoIdAndType(ctx context.Context, userId int64, videoIds []int64, videoType string) ([]*sdomain.VideoUserMapping, error) {
	do := r.do(ctx, nil)
	if userId != 0 {
		do = do.Where(gen.VideoUserMapping.UserID.Eq(userId))
	}
	if len(videoIds) > 0 {
		do.Where(gen.VideoUserMapping.VideoID.In(videoIds...))
	}
	if videoType != "" {
		do = do.Where(gen.VideoUserMapping.VideoType.Eq(videoType))
	}
	finds, err := do.Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return lo.Map(finds, func(item *model.VideoUserMapping, index int) *sdomain.VideoUserMapping {
		return (&sdomain.VideoUserMapping{}).ConvertFromRepo(item)
	}), nil
}

func (r *VideoUserMappingRepo) UpdateFavorite(ctx context.Context, tx *gen.Query, userId int64, videoId int64, videoType string, favorite bool) error {
	first, err := r.do(ctx, tx).Where(gen.VideoUserMapping.UserID.Eq(userId)).Where(gen.VideoUserMapping.VideoID.Eq(videoId)).Where(gen.VideoUserMapping.VideoType.Eq(videoType)).First()
	if err != nil {
		return err
	}
	if first != nil {
		first.Favorited = lo.ToPtr(favorite)
		updates, err := r.do(ctx, tx).Select(gen.VideoUserMapping.Favorited).Updates(first)
		if err != nil {
			return err
		}
		if updates.RowsAffected != 1 {
			return gorm.ErrDuplicatedKey
		}
	} else {
		err = r.do(ctx, tx).Create(&model.VideoUserMapping{
			VideoID:   videoId,
			UserID:    userId,
			VideoType: videoType,
			Favorited: lo.ToPtr(favorite),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *VideoUserMappingRepo) UpdatePlayStatus(ctx context.Context, tx *gen.Query, userId int64, videoId int64, videoType string, position int32) error {
	first, err := r.do(ctx, tx).Where(gen.VideoUserMapping.UserID.Eq(userId)).Where(gen.VideoUserMapping.VideoID.Eq(videoId)).Where(gen.VideoUserMapping.VideoType.Eq(videoType)).First()
	if err != nil {
		return err
	}
	if first != nil {
		first.LastPlayedPosition = lo.ToPtr(position)
		first.LastPlayedTime = lo.ToPtr(time.Now())
		updates, err := r.do(ctx, tx).Updates(first)
		if err != nil {
			return err
		}
		if updates.RowsAffected != 1 {
			return gorm.ErrDuplicatedKey
		}
	} else {
		err = r.do(ctx, tx).Create(&model.VideoUserMapping{
			VideoID:            videoId,
			UserID:             userId,
			VideoType:          videoType,
			Favorited:          lo.ToPtr(false),
			LastPlayedPosition: lo.ToPtr(position),
			LastPlayedTime:     lo.ToPtr(time.Now()),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
