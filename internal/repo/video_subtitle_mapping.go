package repo

import (
	"context"
	"github.com/samber/lo"
	"gnboot/internal/repo/gen"
	"gnboot/internal/repo/model"
	"gnboot/internal/service/sdomain"
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
		return nil, handleQueryError(err)
	}
	return lo.Map(finds, func(item *model.VideoSubtitleMapping, index int) *sdomain.VideoSubtitleMapping {
		return (&sdomain.VideoSubtitleMapping{}).ConvertFromRepo(item)
	}), nil
}
