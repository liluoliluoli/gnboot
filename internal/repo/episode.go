package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

type EpisodeRepo struct {
	Data *Data
}

func NewEpisodeRepo(data *Data) *EpisodeRepo {
	return &EpisodeRepo{
		Data: data,
	}
}

func (r *EpisodeRepo) do(ctx context.Context, tx *gen.Query) gen.IEpisodeDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).Episode.WithContext(ctx)
	} else {
		return tx.Episode.WithContext(ctx)
	}
}

func (r *EpisodeRepo) Get(ctx context.Context, id int64) (*sdomain.Episode, error) {
	find, err := r.do(ctx, nil).Where(gen.Episode.ID.Eq(id)).First()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return (&sdomain.Episode{}).ConvertFromRepo(find), nil
}

func (r *EpisodeRepo) QueryBySeasonId(ctx context.Context, seasonId int64) ([]*sdomain.Episode, error) {
	finds, err := r.do(ctx, nil).Order(gen.Episode.Episode).Find()
	if err != nil {
		return nil, handleQueryError(err)
	}
	return lo.Map(finds, func(item *model.Episode, index int) *sdomain.Episode {
		return (&sdomain.Episode{}).ConvertFromRepo(item)
	}), nil
}
