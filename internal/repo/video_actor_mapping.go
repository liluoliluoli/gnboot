package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

type VideoActorMappingRepo struct {
	Data *Data
}

func NewVideoActorMappingRepo(data *Data) *VideoActorMappingRepo {
	return &VideoActorMappingRepo{
		Data: data,
	}
}

func (r *VideoActorMappingRepo) do(ctx context.Context, tx *gen.Query) gen.IVideoActorMappingDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).VideoActorMapping.WithContext(ctx)
	} else {
		return tx.VideoActorMapping.WithContext(ctx)
	}
}

func (r *VideoActorMappingRepo) FindByVideoIdAndType(ctx context.Context, videoId []int64, videoType string) ([]*sdomain.VideoActorMapping, error) {
	finds, err := r.do(ctx, nil).Where(gen.VideoActorMapping.VideoID.In(videoId...)).Where(gen.VideoActorMapping.VideoType.Eq(videoType)).Find()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return lo.Map(finds, func(item *model.VideoActorMapping, index int) *sdomain.VideoActorMapping {
		return (&sdomain.VideoActorMapping{}).ConvertFromRepo(item)
	}), nil
}

func (r *VideoActorMappingRepo) FindByActorIdAndVideoType(ctx context.Context, actorId int64, videoType string) ([]*sdomain.VideoActorMapping, error) {
	finds, err := r.do(ctx, nil).Where(gen.VideoActorMapping.ActorID.Eq(actorId)).Where(gen.VideoActorMapping.VideoType.Eq(videoType)).Find()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return lo.Map(finds, func(item *model.VideoActorMapping, index int) *sdomain.VideoActorMapping {
		return (&sdomain.VideoActorMapping{}).ConvertFromRepo(item)
	}), nil
}
