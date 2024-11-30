package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
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
