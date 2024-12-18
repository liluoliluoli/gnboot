package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

type VideoSubtitleMappingRepo struct {
	Data *Data
}

func NewVideoSubtitleMappingRepo(data *Data) *VideoSubtitleMappingRepo {
	return &VideoSubtitleMappingRepo{
		Data: data,
	}
}

func (r *VideoSubtitleMappingRepo) do(ctx context.Context, tx *gen.Query) gen.IVideoSubtitleMappingDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).VideoSubtitleMapping.WithContext(ctx)
	} else {
		return tx.VideoSubtitleMapping.WithContext(ctx)
	}
}

func (r *VideoSubtitleMappingRepo) FindByVideoIdAndType(ctx context.Context, videoId []int64, videoType string) ([]*sdomain.VideoSubtitleMapping, error) {
	finds, err := r.do(ctx, nil).Where(gen.VideoSubtitleMapping.VideoID.In(videoId...)).Where(gen.VideoSubtitleMapping.VideoType.Eq(videoType)).Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return lo.Map(finds, func(item *model.VideoSubtitleMapping, index int) *sdomain.VideoSubtitleMapping {
		return (&sdomain.VideoSubtitleMapping{}).ConvertFromRepo(item)
	}), nil
}
