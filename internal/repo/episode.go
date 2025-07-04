package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"time"
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

func (r *EpisodeRepo) Next(ctx context.Context, videoId int64, id int64) (*sdomain.Episode, error) {
	finds, err := r.do(ctx, nil).Where(gen.Episode.VideoID.Eq(videoId)).Order(gen.Episode.EpisodeTitle).Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	_, index, exist := lo.FindIndexOf(finds, func(item *model.Episode) bool {
		return item.ID == id
	})
	var nextEpisode *model.Episode
	if exist {
		if index < len(finds)-1 {
			nextEpisode = finds[index+1]
		}
	} else {
		nextEpisode = finds[0]
	}
	return (&sdomain.Episode{}).ConvertFromRepo(nextEpisode), nil
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

func (r *EpisodeRepo) QueryByPathAndName(ctx context.Context, path, name string) (*model.Episode, error) {
	if path == "" || name == "" {
		return nil, nil
	}
	query := r.do(ctx, nil).Where(gen.Episode.XiaoyaPath.Eq(path)).Where(gen.Episode.EpisodeTitle.Eq(name))
	find, err := query.First()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return find, nil
}

func (r *EpisodeRepo) QueryLastJfCreateTimeByJfId(ctx context.Context, jfRootId string) (*time.Time, error) {
	finds, err := r.do(ctx, nil).Select(gen.Episode.JfCreateTime.Max().As("jf_create_time")).Where(gen.Episode.JfRootPathID.Eq(jfRootId)).Group(gen.Episode.JfRootPathID).Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	first, _ := lo.First(finds)
	if first == nil {
		return nil, nil
	}
	return first.JfCreateTime, nil
}

func (r *EpisodeRepo) QueryLastPublishTimeByJfId(ctx context.Context, jfRootId string) (*time.Time, error) {
	finds, err := r.do(ctx, nil).Select(gen.Episode.JfPublishTime.Max().As("jf_publish_time")).Where(gen.Episode.JfRootPathID.Eq(jfRootId)).Group(gen.Episode.JfRootPathID).Find()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	first, _ := lo.First(finds)
	if first == nil {
		return nil, nil
	}
	return first.JfCreateTime, nil
}
