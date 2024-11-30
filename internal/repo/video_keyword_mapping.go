package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

type VideoKeywordMappingRepo struct {
	Data *Data
}

func NewVideoKeywordMappingRepo(data *Data) *VideoKeywordMappingRepo {
	return &VideoKeywordMappingRepo{
		Data: data,
	}
}

func (r *VideoKeywordMappingRepo) do(ctx context.Context, tx *gen.Query) gen.IVideoKeywordMappingDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).VideoKeywordMapping.WithContext(ctx)
	} else {
		return tx.VideoKeywordMapping.WithContext(ctx)
	}
}

func (r *VideoKeywordMappingRepo) FindByVideoIdAndType(ctx context.Context, videoId []int64, videoType string) ([]*sdomain.VideoKeywordMapping, error) {
	finds, err := r.do(ctx, nil).Where(gen.VideoKeywordMapping.VideoID.In(videoId...)).Where(gen.VideoKeywordMapping.VideoType.Eq(videoType)).Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return lo.Map(finds, func(item *model.VideoKeywordMapping, index int) *sdomain.VideoKeywordMapping {
		return (&sdomain.VideoKeywordMapping{}).ConvertFromRepo(item)
	}), nil
}

func (r *VideoKeywordMappingRepo) FindByKeywordIdAndVideoType(ctx context.Context, keywordId int64, videoType string) ([]*sdomain.VideoKeywordMapping, error) {
	finds, err := r.do(ctx, nil).Where(gen.VideoKeywordMapping.KeywordID.Eq(keywordId)).Where(gen.VideoKeywordMapping.VideoType.Eq(videoType)).Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return lo.Map(finds, func(item *model.VideoKeywordMapping, index int) *sdomain.VideoKeywordMapping {
		return (&sdomain.VideoKeywordMapping{}).ConvertFromRepo(item)
	}), nil
}
