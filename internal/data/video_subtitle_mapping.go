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

type videoSubtitleMappingRepo struct {
	data *Data
}

func NewVideoSubtitleMappingRepo(data *Data) biz.VideoSubtitleMappingRepo {
	return &videoSubtitleMappingRepo{
		data: data,
	}
}

func (ro videoSubtitleMappingRepo) Create(ctx context.Context, item *biz.CreateVideoSubtitleMapping) (err error) {
	err = ro.NameExists(ctx, item.Name)
	if err == nil {
		err = biz.ErrDuplicateField(ctx, "name", item.Name)
		return
	}
	var m model.VideoSubtitleMapping
	copierx.Copy(&m, item)
	p := query.Use(ro.data.DB(ctx)).VideoSubtitleMapping
	db := p.WithContext(ctx)
	m.ID = ro.data.ID(ctx)
	err = db.Create(&m)
	return
}

func (ro videoSubtitleMappingRepo) Get(ctx context.Context, id uint64) (item *biz.VideoSubtitleMapping, err error) {
	item = &biz.VideoSubtitleMapping{}
	p := query.Use(ro.data.DB(ctx)).VideoSubtitleMapping
	db := p.WithContext(ctx)
	m := db.GetByID(id)
	if m.ID == constant.UI0 {
		err = biz.ErrRecordNotFound(ctx)
		return
	}
	copierx.Copy(&item, m)
	return
}

func (ro videoSubtitleMappingRepo) Find(ctx context.Context, condition *biz.FindVideoSubtitleMapping) (rp []biz.VideoSubtitleMapping) {
	p := query.Use(ro.data.DB(ctx)).VideoSubtitleMapping
	db := p.WithContext(ctx)
	rp = make([]biz.VideoSubtitleMapping, 0)
	list := make([]model.VideoSubtitleMapping, 0)
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

func (ro videoSubtitleMappingRepo) Update(ctx context.Context, item *biz.UpdateVideoSubtitleMapping) (err error) {
	p := query.Use(ro.data.DB(ctx)).VideoSubtitleMapping
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
	_, err = db.
		Where(p.ID.Eq(item.ID)).
		Updates(&change)
	return
}

func (ro videoSubtitleMappingRepo) Delete(ctx context.Context, ids ...uint64) (err error) {
	p := query.Use(ro.data.DB(ctx)).VideoSubtitleMapping
	db := p.WithContext(ctx)
	_, err = db.
		Where(p.ID.In(ids...)).
		Delete()
	return
}

func (ro videoSubtitleMappingRepo) NameExists(ctx context.Context, name string) (err error) {
	p := query.Use(ro.data.DB(ctx)).VideoSubtitleMapping
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
