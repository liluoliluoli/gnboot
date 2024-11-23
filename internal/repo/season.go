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

type seasonRepo struct {
	data *Data
}

func NewSeasonRepo(data *Data) service.SeasonRepo {
	return &seasonRepo{
		data: data,
	}
}

func (ro seasonRepo) Create(ctx context.Context, item *service.CreateSeason) (err error) {
	err = ro.NameExists(ctx, item.Name)
	if err == nil {
		err = service.ErrDuplicateField(ctx, "name", item.Name)
		return
	}
	var m model.Season
	copierx.Copy(&m, item)
	p := query.Use(ro.data.DB(ctx)).Season
	db := p.WithContext(ctx)
	m.ID = ro.data.ID(ctx)
	err = db.Create(&m)
	return
}

func (ro seasonRepo) Get(ctx context.Context, id uint64) (item *service.Season, err error) {
	item = &service.Season{}
	p := query.Use(ro.data.DB(ctx)).Season
	db := p.WithContext(ctx)
	m := db.GetByID(id)
	if m.ID == constant.UI0 {
		err = service.ErrRecordNotFound(ctx)
		return
	}
	copierx.Copy(&item, m)
	return
}

func (ro seasonRepo) Find(ctx context.Context, condition *service.FindSeason) (rp []service.Season) {
	p := query.Use(ro.data.DB(ctx)).Season
	db := p.WithContext(ctx)
	rp = make([]service.Season, 0)
	list := make([]model.Season, 0)
	conditions := make([]gen.Condition, 0, 2)
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

func (ro seasonRepo) Update(ctx context.Context, item *service.UpdateSeason) (err error) {
	p := query.Use(ro.data.DB(ctx)).Season
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
	_, err = db.
		Where(p.ID.Eq(item.ID)).
		Updates(&change)
	return
}

func (ro seasonRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	p := query.Use(ro.data.DB(ctx)).Season
	db := p.WithContext(ctx)
	_, err = db.
		Where(p.ID.In(ids...)).
		Delete()
	return
}

func (ro seasonRepo) NameExists(ctx context.Context, name string) (err error) {
	p := query.Use(ro.data.DB(ctx)).Season
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
