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

type appVersionRepo struct {
	data *Data
}

func NewAppVersionRepo(data *Data) biz.AppVersionRepo {
	return &appVersionRepo{
		data: data,
	}
}

func (ro appVersionRepo) Create(ctx context.Context, item *biz.CreateAppVersion) (err error) {
	err = ro.NameExists(ctx, item.Name)
	if err == nil {
		err = biz.ErrDuplicateField(ctx, "name", item.Name)
		return
	}
	var m model.AppVersion
	copierx.Copy(&m, item)
	p := query.Use(ro.data.DB(ctx)).AppVersion
	db := p.WithContext(ctx)
	m.ID = ro.data.ID(ctx)
	err = db.Create(&m)
	return
}

func (ro appVersionRepo) Get(ctx context.Context, id uint64) (item *biz.AppVersion, err error) {
	item = &biz.AppVersion{}
	p := query.Use(ro.data.DB(ctx)).AppVersion
	db := p.WithContext(ctx)
	m := db.GetByID(id)
	if m.ID == constant.UI0 {
		err = biz.ErrRecordNotFound(ctx)
		return
	}
	copierx.Copy(&item, m)
	return
}

func (ro appVersionRepo) Find(ctx context.Context, condition *biz.FindAppVersion) (rp []biz.AppVersion) {
	p := query.Use(ro.data.DB(ctx)).AppVersion
	db := p.WithContext(ctx)
	rp = make([]biz.AppVersion, 0)
	list := make([]model.AppVersion, 0)
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

func (ro appVersionRepo) Update(ctx context.Context, item *biz.UpdateAppVersion) (err error) {
	p := query.Use(ro.data.DB(ctx)).AppVersion
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
	if item.Name != nil && *item.Name != m.Name {
		err = ro.NameExists(ctx, *item.Name)
		if err == nil {
			err = biz.ErrDuplicateField(ctx, "name", *item.Name)
			return
		}
	}
	_, err = db.
		Where(p.ID.Eq(item.ID)).
		Updates(&change)
	return
}

func (ro appVersionRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	p := query.Use(ro.data.DB(ctx)).AppVersion
	db := p.WithContext(ctx)
	_, err = db.
		Where(p.ID.In(ids...)).
		Delete()
	return
}

func (ro appVersionRepo) NameExists(ctx context.Context, name string) (err error) {
	p := query.Use(ro.data.DB(ctx)).AppVersion
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
