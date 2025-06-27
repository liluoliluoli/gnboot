package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

type EpisodeAudioMappingRepo struct {
	Data *Data
}

func NewEpisodeAudioMappingRepo(data *Data) *EpisodeAudioMappingRepo {
	return &EpisodeAudioMappingRepo{
		Data: data,
	}
}

func (r *EpisodeAudioMappingRepo) do(ctx context.Context, tx *gen.Query) gen.IEpisodeAudioMappingDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).EpisodeAudioMapping.WithContext(ctx)
	} else {
		return tx.EpisodeAudioMapping.WithContext(ctx)
	}
}

func (r *EpisodeAudioMappingRepo) FindByEpisodeId(ctx context.Context, episodeId int64) ([]*sdomain.EpisodeAudioMapping, error) {
	finds, err := r.do(ctx, nil).Where(gen.EpisodeAudioMapping.EpisodeID.Eq(episodeId)).Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return lo.Map(finds, func(item *model.EpisodeAudioMapping, index int) *sdomain.EpisodeAudioMapping {
		return (&sdomain.EpisodeAudioMapping{}).ConvertFromRepo(item)
	}), nil
}

func (r *EpisodeAudioMappingRepo) Delete(ctx context.Context, tx *gen.Query, episodeId int64) error {
	_, err := r.do(ctx, tx).Where(gen.EpisodeAudioMapping.EpisodeID.Eq(episodeId)).Delete()
	if err != nil {
		return err
	}
	return nil
}

func (r *EpisodeAudioMappingRepo) Create(ctx context.Context, tx *gen.Query, mapping *model.EpisodeAudioMapping) error {
	err := r.do(ctx, tx).Save(mapping)
	if err != nil {
		return err
	}
	return nil
}
