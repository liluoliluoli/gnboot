package data

import (
	"context"
	"strings"

	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"
	"gnboot/internal/biz"
	"gnboot/internal/data/model"
	"gnboot/internal/data/query"
	"gorm.io/gen"
)

type movieRepo struct {
	data *Data
}

func NewMovieRepo(data *Data) biz.MovieRepo {
	return &movieRepo{
		data: data,
	}
}

func (ro movieRepo) Create(ctx context.Context, item *biz.CreateMovie) (err error) {
	err = ro.NameExists(ctx, item.Name)
	if err == nil {
		err = biz.ErrDuplicateField(ctx, "name", item.Name)
		return
	}
	var m model.Movie
	copierx.Copy(&m, item)
	p := query.Use(ro.data.DB(ctx)).Movie
	db := p.WithContext(ctx)
	m.ID = ro.data.ID(ctx)
	err = db.Create(&m)
	return
}

func (ro movieRepo) Get(ctx context.Context, id uint64) (item *biz.Movie, err error) {
	item = &biz.Movie{}
	p := query.Use(ro.data.DB(ctx)).Movie
	db := p.WithContext(ctx)
	m := db.GetByID(id)
	if m.ID == constant.UI0 {
		err = biz.ErrRecordNotFound(ctx)
		return
	}
	copierx.Copy(&item, m)
	return
}

func (ro movieRepo) Find(ctx context.Context, condition *biz.FindMovie) (rp []biz.Movie) {
	p := query.Use(ro.data.DB(ctx)).Movie
	db := p.WithContext(ctx)
	rp = make([]biz.Movie, 0)
	list := make([]model.Movie, 0)
	conditions := make([]gen.Condition, 0, 2)
	if condition.Search != nil {
		conditions = append(conditions, p.OriginalTitle.Like(strings.Join([]string{"%", *condition.Search, "%"}, "")))
	}

	condition.Page.Primary = "id"
	condition.Page.
		WithContext(ctx).
		Query(
			db.
				Order(p.ID.Desc()).
				Where(conditions...).
				UnderlyingDB(),
		).
		Find(&list)
	copierx.Copy(&rp, list)
	return
}

func (ro movieRepo) Update(ctx context.Context, item *biz.UpdateMovie) (err error) {
	p := query.Use(ro.data.DB(ctx)).Movie
	db := p.WithContext(ctx)
	m := db.GetByID(item.ID)
	if m.ID == constant.UI0 {
		err = biz.ErrRecordNotFound(ctx)
		return
	}
	change := make(map[string]interface{})
	utils.CompareDiff(m, item, &change)
	if len(change) == 0 {
		err = biz.ErrDataNotChange(ctx)
		return
	}
	if item.Title != nil && *item.Title != m.OriginalTitle {
		err = ro.NameExists(ctx, *item.Title)
		if err == nil {
			err = biz.ErrDuplicateField(ctx, "name", *item.Title)
			return
		}
	}
	_, err = db.
		Where(p.ID.Eq(item.ID)).
		Updates(&change)
	return
}

func (ro movieRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	p := query.Use(ro.data.DB(ctx)).Movie
	db := p.WithContext(ctx)
	_, err = db.
		Where(p.ID.In(ids...)).
		Delete()
	return
}

func (ro movieRepo) NameExists(ctx context.Context, name string) (err error) {
	p := query.Use(ro.data.DB(ctx)).Movie
	db := p.WithContext(ctx)
	arr := strings.Split(name, ",")
	for _, item := range arr {
		res := db.GetByCol("name", item)
		if res.ID == constant.UI0 {
			err = biz.ErrRecordNotFound(ctx)
			log.
				WithError(err).
				Warn("invalid `name`: %s", name)
			return
		}
	}
	return
}
