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

type videoActorMappingRepo struct {
	data *Data
}

func NewVideoActorMappingRepo(data *Data) biz.VideoActorMappingRepo {
	return &videoActorMappingRepo{
		data: data,
	}
}

func (ro videoActorMappingRepo) Create(ctx context.Context, item *biz.CreateVideoActorMapping) (err error) {
	err = ro.NameExists(ctx, item.Name)
	if err == nil {
		err = biz.ErrDuplicateField(ctx, "name", item.Name)
		return
	}
	var m model.VideoActorMapping
	copierx.Copy(&m, item)
	p := query.Use(ro.data.DB(ctx)).VideoActorMapping
	db := p.WithContext(ctx)
	m.ID = ro.data.ID(ctx)
	err = db.Create(&m)
	return
}

func (ro videoActorMappingRepo) Get(ctx context.Context, id uint64) (item *biz.VideoActorMapping, err error) {
	item = &biz.VideoActorMapping{}
	p := query.Use(ro.data.DB(ctx)).VideoActorMapping
	db := p.WithContext(ctx)
	m := db.GetByID(id)
	if m.ID == constant.UI0 {
		err = biz.ErrRecordNotFound(ctx)
		return
	}
	copierx.Copy(&item, m)
	return
}

func (ro videoActorMappingRepo) Find(ctx context.Context, condition *biz.FindVideoActorMapping) (rp []biz.VideoActorMapping) {
	p := query.Use(ro.data.DB(ctx)).VideoActorMapping
	db := p.WithContext(ctx)
	rp = make([]biz.VideoActorMapping, 0)
	list := make([]model.VideoActorMapping, 0)
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

func (ro videoActorMappingRepo) Update(ctx context.Context, item *biz.UpdateVideoActorMapping) (err error) {
	p := query.Use(ro.data.DB(ctx)).VideoActorMapping
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

func (ro videoActorMappingRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	p := query.Use(ro.data.DB(ctx)).VideoActorMapping
	db := p.WithContext(ctx)
	_, err = db.
		Where(p.ID.In(ids...)).
		Delete()
	return
}

func (ro videoActorMappingRepo) NameExists(ctx context.Context, name string) (err error) {
	p := query.Use(ro.data.DB(ctx)).VideoActorMapping
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
