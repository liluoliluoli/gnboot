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

type videoUserMappingRepo struct {
	data *Data
}

func NewVideoUserMappingRepo(data *Data) service.VideoUserMappingRepo {
	return &videoUserMappingRepo{
		data: data,
	}
}

func (ro videoUserMappingRepo) Create(ctx context.Context, item *service.CreateVideoUserMapping) (err error) {
	err = ro.NameExists(ctx, item.Name)
	if err == nil {
		err = service.ErrDuplicateField(ctx, "name", item.Name)
		return
	}
	var m model.VideoUserMapping
	copierx.Copy(&m, item)
	p := query.Use(ro.data.DB(ctx)).VideoUserMapping
	db := p.WithContext(ctx)
	m.ID = ro.data.ID(ctx)
	err = db.Create(&m)
	return
}

func (ro videoUserMappingRepo) Get(ctx context.Context, id uint64) (item *service.VideoUserMapping, err error) {
	item = &service.VideoUserMapping{}
	p := query.Use(ro.data.DB(ctx)).VideoUserMapping
	db := p.WithContext(ctx)
	m := db.GetByID(id)
	if m.ID == constant.UI0 {
		err = service.ErrRecordNotFound(ctx)
		return
	}
	copierx.Copy(&item, m)
	return
}

func (ro videoUserMappingRepo) Find(ctx context.Context, condition *service.FindVideoUserMapping) (rp []service.VideoUserMapping) {
	p := query.Use(ro.data.DB(ctx)).VideoUserMapping
	db := p.WithContext(ctx)
	rp = make([]service.VideoUserMapping, 0)
	list := make([]model.VideoUserMapping, 0)
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

func (ro videoUserMappingRepo) Update(ctx context.Context, item *service.UpdateVideoUserMapping) (err error) {
	p := query.Use(ro.data.DB(ctx)).VideoUserMapping
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

func (ro videoUserMappingRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	p := query.Use(ro.data.DB(ctx)).VideoUserMapping
	db := p.WithContext(ctx)
	_, err = db.
		Where(p.ID.In(ids...)).
		Delete()
	return
}

func (ro videoUserMappingRepo) NameExists(ctx context.Context, name string) (err error) {
	p := query.Use(ro.data.DB(ctx)).VideoUserMapping
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