package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

type EpisodeSubtitleMappingRepo struct {
	Data *Data
}

func NewEpisodeSubtitleMappingRepo(data *Data) *EpisodeSubtitleMappingRepo {
	return &EpisodeSubtitleMappingRepo{
		Data: data,
	}
}

func (r *EpisodeSubtitleMappingRepo) do(ctx context.Context, tx *gen.Query) gen.IEpisodeSubtitleMappingDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).EpisodeSubtitleMapping.WithContext(ctx)
	} else {
		return tx.EpisodeSubtitleMapping.WithContext(ctx)
	}
}

func (r *EpisodeSubtitleMappingRepo) FindByEpisodeId(ctx context.Context, episodeId int64) ([]*sdomain.EpisodeSubtitleMapping, error) {
	finds, err := r.do(ctx, nil).Where(gen.EpisodeSubtitleMapping.EpisodeID.Eq(episodeId)).Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return lo.Map(finds, func(item *model.EpisodeSubtitleMapping, index int) *sdomain.EpisodeSubtitleMapping {
		return (&sdomain.EpisodeSubtitleMapping{}).ConvertFromRepo(item)
	}), nil
}

func (r *EpisodeSubtitleMappingRepo) Delete(ctx context.Context, tx *gen.Query, episodeId int64) error {
	_, err := r.do(ctx, tx).Where(gen.EpisodeSubtitleMapping.EpisodeID.Eq(episodeId)).Delete()
	if err != nil {
		return err
	}
	return nil
}

func (r *EpisodeSubtitleMappingRepo) Create(ctx context.Context, tx *gen.Query, mapping *model.EpisodeSubtitleMapping) error {
	err := r.do(ctx, tx).Save(mapping)
	if err != nil {
		return err
	}
	return nil
}
