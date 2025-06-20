package repo

import (
	"context"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type VideoRepo struct {
	Data *Data
}

func NewVideoRepo(data *Data) *VideoRepo {
	return &VideoRepo{
		Data: data,
	}
}

func (r *VideoRepo) do(ctx context.Context, tx *gen.Query) gen.IVideoDo {
	if tx == nil {
		return gen.Use(r.Data.DB(ctx)).Video.WithContext(ctx)
	} else {
		return tx.Video.WithContext(ctx)
	}
}

func (r *VideoRepo) Get(ctx context.Context, id int64) (*sdomain.Video, error) {
	find, err := r.do(ctx, nil).Where(gen.Video.ID.Eq(id)).First()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return (&sdomain.Video{}).ConvertFromRepo(find), nil
}

func (r *VideoRepo) GetByJellyfinId(ctx context.Context, jellyfinId string) (*sdomain.Video, error) {
	find, err := r.do(ctx, nil).Where(gen.Video.JellyfinID.Eq(jellyfinId)).First()
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return (&sdomain.Video{}).ConvertFromRepo(find), nil
}

func (r *VideoRepo) Page(ctx context.Context, condition *sdomain.VideoSearch) (*sdomain.PageResult[*sdomain.Video], error) {
	do := r.do(ctx, nil)
	if condition.Search != "" {
		do = do.Where(gen.Video.Title.Like("%" + condition.Search + "%"))
	}
	if len(condition.Ids) > 0 {
		do = do.Where(gen.Video.ID.In(condition.Ids...))
	}
	if condition.Type != "" {
		do = do.Where(gen.Video.VideoType.Eq(condition.Type))
	}
	if condition.Genre != "" {
		do = do.Where(gen.Video.Genres.Like("%" + condition.Genre + "%"))
	}
	if condition.Region != "" {
		do = do.Where(gen.Video.Region.Eq(condition.Region))
	}
	if condition.Year != "" {
		do = do.Where(gen.Video.PublishDay.Between(condition.Year+"0101", condition.Year+"1231"))
	}
	if condition.Sort == constant.SortByHot {
		do = do.Order(gen.Video.VoteCount.Desc())
	} else if condition.Sort == constant.SortByRate {
		do = do.Order(gen.Video.VoteRate.Desc())
	} else if condition.Sort == constant.SortByPublish {
		do = do.Order(gen.Video.PublishDay.Desc())
	} else {
		do = do.Order(gen.Video.UpdateTime.Desc())
	}
	list, total, err := do.FindByPage(int((condition.Page.CurrentPage-1)*condition.Page.PageSize), int(condition.Page.PageSize))
	if err != nil {
		return nil, handleQueryError(ctx, err)
	}
	return &sdomain.PageResult[*sdomain.Video]{
		Page: &sdomain.Page{
			CurrentPage: condition.Page.CurrentPage,
			PageSize:    condition.Page.PageSize,
			Count:       total,
		},
		List: lo.Map(list, func(item *model.Video, index int) *sdomain.Video {
			return (&sdomain.Video{}).ConvertFromRepo(item)
		}),
	}, nil
}

func (r *VideoRepo) Create(ctx context.Context, tx *gen.Query, movie *sdomain.Video) error {
	err := r.do(ctx, tx).Save(movie.ConvertToRepo())
	if err != nil {
		return err
	}
	return nil
}

func (r *VideoRepo) Update(ctx context.Context, tx *gen.Query, movie *sdomain.UpdateVideo) error {
	updates, err := r.do(ctx, tx).Updates(movie.ConvertToRepo())
	if err != nil {
		return err
	}
	if updates.RowsAffected != 1 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

func (r *VideoRepo) Delete(ctx context.Context, tx *gen.Query, ids ...int64) error {
	_, err := r.do(ctx, tx).Where(gen.Video.ID.In(ids...)).Delete()
	if err != nil {
		return err
	}
	return nil
}
