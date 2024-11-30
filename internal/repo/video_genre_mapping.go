package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

type VideoGenreMappingRepo struct {
	Data *Data
}

func NewVideoGenreMappingRepo(data *Data) *VideoGenreMappingRepo {
	return &VideoGenreMappingRepo{
		Data: data,
	}
}

func (r *VideoGenreMappingRepo) do(ctx context.Context, tx *gen.Query) gen.IVideoGenreMappingDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).VideoGenreMapping.WithContext(ctx)
	} else {
		return tx.VideoGenreMapping.WithContext(ctx)
	}
}

func (r *VideoGenreMappingRepo) FindByVideoIdAndType(ctx context.Context, videoId []int64, videoType string) ([]*sdomain.VideoGenreMapping, error) {
	finds, err := r.do(ctx, nil).Where(gen.VideoGenreMapping.VideoID.In(videoId...)).Where(gen.VideoGenreMapping.VideoType.Eq(videoType)).Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return lo.Map(finds, func(item *model.VideoGenreMapping, index int) *sdomain.VideoGenreMapping {
		return (&sdomain.VideoGenreMapping{}).ConvertFromRepo(item)
	}), nil
}

func (r *VideoGenreMappingRepo) FindByGenreIdAndVideoType(ctx context.Context, genreId int64, videoType string) ([]*sdomain.VideoGenreMapping, error) {
	finds, err := r.do(ctx, nil).Where(gen.VideoGenreMapping.GenreID.Eq(genreId)).Where(gen.VideoGenreMapping.VideoType.Eq(videoType)).Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return lo.Map(finds, func(item *model.VideoGenreMapping, index int) *sdomain.VideoGenreMapping {
		return (&sdomain.VideoGenreMapping{}).ConvertFromRepo(item)
	}), nil
}
