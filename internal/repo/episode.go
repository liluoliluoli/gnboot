package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
	"gorm.io/gorm"
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
		return nil, handleQueryError(ctx, err)
	}
	return (&sdomain.Episode{}).ConvertFromRepo(find), nil
}

func (r *EpisodeRepo) QueryByVideoId(ctx context.Context, videoId int64) ([]*sdomain.Episode, error) {
	finds, err := r.do(ctx, nil).Where(gen.Episode.VideoID.Eq(videoId)).Order(gen.Episode.Episode).Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return lo.Map(finds, func(item *model.Episode, index int) *sdomain.Episode {
		return (&sdomain.Episode{}).ConvertFromRepo(item)
	}), nil
}

func (r *EpisodeRepo) QueryByIds(ctx context.Context, ids []int64) ([]*sdomain.Episode, error) {
	finds, err := r.do(ctx, nil).Where(gen.Episode.ID.In(ids...)).Order(gen.Episode.Episode).Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return lo.Map(finds, func(item *model.Episode, index int) *sdomain.Episode {
		return (&sdomain.Episode{}).ConvertFromRepo(item)
	}), nil
}

func (r *EpisodeRepo) Create(ctx context.Context, tx *gen.Query, episode *model.Episode) error {
	err := r.do(ctx, tx).Save(episode)
	if err != nil {
		return err
	}
	return nil
}

func (r *EpisodeRepo) Updates(ctx context.Context, tx *gen.Query, episode *sdomain.Episode) error {
	updates, err := r.do(ctx, tx).Updates(episode.ConvertToRepo())
	if err != nil {
		return err
	}
	if updates.RowsAffected != 1 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

// QueryByPathAndName 实现查询 episode 表，判断传入参数 path 是否在 xiaoya_Path 中存在，
// 传入参数 name 是否在 episode_title 中存在，如果两个条件都不存在则返回空
func (r *EpisodeRepo) QueryByPathAndName(ctx context.Context, path, name string) ([]*sdomain.Episode, error) {
	query := r.do(ctx, nil)
	if path != "" {
		query = query.Where(gen.Episode.XiaoyaPath.Eq(path))
	}
	if name != "" {
		query = query.Where(gen.Episode.EpisodeTitle.Eq(name))
	}
	if path == "" && name == "" {
		return []*sdomain.Episode{}, nil
	}
	finds, err := query.Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return lo.Map(finds, func(item *model.Episode, index int) *sdomain.Episode {
		return (&sdomain.Episode{}).ConvertFromRepo(item)
	}), nil
}
