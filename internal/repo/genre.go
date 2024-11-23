package repo

import (
	"context"
	"strings"

	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/utils"
	"gnboot/internal/repo/model"
	"gnboot/internal/repo/query"
	"gnboot/internal/service"
	"gorm.io/gen"
)

type genreRepo struct {
	data *Data
}

func NewGenreRepo(data *Data) service.GenreRepo {
	return &genreRepo{
		data: data,
	}
}

func (ro genreRepo) Create(ctx context.Context, item *service.CreateGenre) (err error) {
	err = ro.NameExists(ctx, item.Name)
	if err == nil {
		err = service.ErrDuplicateField(ctx, "name", item.Name)
		return
	}
	var m model.Genre
	copierx.Copy(&m, item)
	p := query.Use(ro.data.DB(ctx)).Genre
	db := p.WithContext(ctx)
	m.ID = ro.data.ID(ctx)
	err = db.Create(&m)
	return
}

func (ro genreRepo) Get(ctx context.Context, id uint64) (item *service.Genre, err error) {
	item = &service.Genre{}
	p := query.Use(ro.data.DB(ctx)).Genre
	db := p.WithContext(ctx)
	m := db.GetByID(id)
	if m.ID == constant.UI0 {
		err = service.ErrRecordNotFound(ctx)
		return
	}
	copierx.Copy(&item, m)
	return
}

func (ro genreRepo) Find(ctx context.Context, condition *service.FindGenre) (rp []service.Genre) {
	p := query.Use(ro.data.DB(ctx)).Genre
	db := p.WithContext(ctx)
	rp = make([]service.Genre, 0)
	list := make([]model.Genre, 0)
	conditions := make([]gen.Condition, 0, 2)
	if condition.Name != nil {
		conditions = append(conditions, p.Name.Like(strings.Join([]string{"%", *condition.Name, "%"}, "")))
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

func (ro genreRepo) Update(ctx context.Context, item *service.UpdateGenre) (err error) {
	p := query.Use(ro.data.DB(ctx)).Genre
	db := p.WithContext(ctx)
	m := db.GetByID(item.ID)
	if m.ID == constant.UI0 {
		err = service.ErrRecordNotFound(ctx)
		return
	}
	change := make(map[string]interface{})
	utils.CompareDiff(m, item, &change)
	if len(change) == 0 {
		err = service.ErrDataNotChange(ctx)
		return
	}
	if item.Name != nil && *item.Name != m.Name {
		err = ro.NameExists(ctx, *item.Name)
		if err == nil {
			err = service.ErrDuplicateField(ctx, "name", *item.Name)
			return
		}
	}
	_, err = db.
		Where(p.ID.Eq(item.ID)).
		Updates(&change)
	return
}

func (ro genreRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	p := query.Use(ro.data.DB(ctx)).Genre
	db := p.WithContext(ctx)
	_, err = db.
		Where(p.ID.In(ids...)).
		Delete()
	return
}

func (ro genreRepo) NameExists(ctx context.Context, name string) (err error) {
	p := query.Use(ro.data.DB(ctx)).Genre
	db := p.WithContext(ctx)
	arr := strings.Split(name, ",")
	for _, item := range arr {
		res := db.GetByCol("name", item)
		if res.ID == constant.UI0 {
			err = service.ErrRecordNotFound(ctx)
			log.
				WithError(err).
				Warn("invalid `name`: %s", name)
			return
		}
	}
	return
}
