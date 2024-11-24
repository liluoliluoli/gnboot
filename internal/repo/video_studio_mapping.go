package repo

import (
	"context"
	"github.com/samber/lo"
	"gnboot/internal/repo/gen"
	"gnboot/internal/repo/model"
	"gnboot/internal/service/sdomain"
)

type VideoStudioMappingRepo struct {
	Data *Data
}

func NewVideoStudioMappingRepo(data *Data) *VideoStudioMappingRepo {
	return &VideoStudioMappingRepo{
		Data: data,
	}
}

func (r *VideoStudioMappingRepo) do(ctx context.Context, tx *gen.Query) gen.IVideoStudioMappingDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).VideoStudioMapping.WithContext(ctx)
	} else {
		return tx.VideoStudioMapping.WithContext(ctx)
	}
}

func (r *VideoStudioMappingRepo) FindByVideoIdAndType(ctx context.Context, videoId []int64, videoType string) ([]*sdomain.VideoStudioMapping, error) {
	finds, err := r.do(ctx, nil).Where(gen.VideoStudioMapping.VideoID.In(videoId...)).Where(gen.VideoStudioMapping.VideoType.Eq(videoType)).Find()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return lo.Map(finds, func(item *model.VideoStudioMapping, index int) *sdomain.VideoStudioMapping {
		return (&sdomain.VideoStudioMapping{}).ConvertFromRepo(item)
	}), nil
}
