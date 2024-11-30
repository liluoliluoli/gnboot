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

func (r *VideoUserMappingRepo) FindByUserIdAndVideoIdAndType(ctx context.Context, userId int64, videoId []int64, videoType string) ([]*sdomain.VideoUserMapping, error) {
	finds, err := r.do(ctx, nil).Where(gen.VideoUserMapping.UserID.Eq(userId)).Where(gen.VideoUserMapping.VideoID.In(videoId...)).Where(gen.VideoUserMapping.VideoType.Eq(videoType)).Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return lo.Map(finds, func(item *model.VideoUserMapping, index int) *sdomain.VideoUserMapping {
		return (&sdomain.VideoUserMapping{}).ConvertFromRepo(item)
	}), nil
}
