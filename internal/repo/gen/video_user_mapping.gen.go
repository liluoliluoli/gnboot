// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package gen

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"gnboot/internal/repo/model"
)

func newVideoUserMapping(db *gorm.DB, opts ...gen.DOOption) videoUserMapping {
	_videoUserMapping := videoUserMapping{}

	_videoUserMapping.videoUserMappingDo.UseDB(db, opts...)
	_videoUserMapping.videoUserMappingDo.UseModel(&model.VideoUserMapping{})

	tableName := _videoUserMapping.videoUserMappingDo.TableName()
	_videoUserMapping.ALL = field.NewAsterisk(tableName)
	_videoUserMapping.ID = field.NewInt64(tableName, "id")
	_videoUserMapping.VideoType = field.NewString(tableName, "video_type")
	_videoUserMapping.VideoID = field.NewInt64(tableName, "video_id")
	_videoUserMapping.LastPlayedPosition = field.NewInt64(tableName, "last_played_position")
	_videoUserMapping.LastPlayedTime = field.NewInt64(tableName, "last_played_time")
	_videoUserMapping.Favorited = field.NewBool(tableName, "favorited")

	_videoUserMapping.fillFieldMap()

	return _videoUserMapping
}

type videoUserMapping struct {
	videoUserMappingDo videoUserMappingDo

	ALL                field.Asterisk
	ID                 field.Int64  // 主键
	VideoType          field.String // 影片类型，movie,series,season,episode
	VideoID            field.Int64  // 影片id，根据video_type类型分别来自movie,series,season,episode表
	LastPlayedPosition field.Int64  // 影片上次播放位置，第n秒
	LastPlayedTime     field.Int64  // 上次播放时候
	Favorited          field.Bool   // 是否收藏喜欢

	fieldMap map[string]field.Expr
}

func (v videoUserMapping) Table(newTableName string) *videoUserMapping {
	v.videoUserMappingDo.UseTable(newTableName)
	return v.updateTableName(newTableName)
}

func (v videoUserMapping) As(alias string) *videoUserMapping {
	v.videoUserMappingDo.DO = *(v.videoUserMappingDo.As(alias).(*gen.DO))
	return v.updateTableName(alias)
}

func (v *videoUserMapping) updateTableName(table string) *videoUserMapping {
	v.ALL = field.NewAsterisk(table)
	v.ID = field.NewInt64(table, "id")
	v.VideoType = field.NewString(table, "video_type")
	v.VideoID = field.NewInt64(table, "video_id")
	v.LastPlayedPosition = field.NewInt64(table, "last_played_position")
	v.LastPlayedTime = field.NewInt64(table, "last_played_time")
	v.Favorited = field.NewBool(table, "favorited")

	v.fillFieldMap()

	return v
}

func (v *videoUserMapping) WithContext(ctx context.Context) IVideoUserMappingDo {
	return v.videoUserMappingDo.WithContext(ctx)
}

func (v videoUserMapping) TableName() string { return v.videoUserMappingDo.TableName() }

func (v videoUserMapping) Alias() string { return v.videoUserMappingDo.Alias() }

func (v *videoUserMapping) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := v.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (v *videoUserMapping) fillFieldMap() {
	v.fieldMap = make(map[string]field.Expr, 6)
	v.fieldMap["id"] = v.ID
	v.fieldMap["video_type"] = v.VideoType
	v.fieldMap["video_id"] = v.VideoID
	v.fieldMap["last_played_position"] = v.LastPlayedPosition
	v.fieldMap["last_played_time"] = v.LastPlayedTime
	v.fieldMap["favorited"] = v.Favorited
}

func (v videoUserMapping) clone(db *gorm.DB) videoUserMapping {
	v.videoUserMappingDo.ReplaceConnPool(db.Statement.ConnPool)
	return v
}

func (v videoUserMapping) replaceDB(db *gorm.DB) videoUserMapping {
	v.videoUserMappingDo.ReplaceDB(db)
	return v
}

type videoUserMappingDo struct{ gen.DO }

type IVideoUserMappingDo interface {
	gen.SubQuery
	Debug() IVideoUserMappingDo
	WithContext(ctx context.Context) IVideoUserMappingDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IVideoUserMappingDo
	WriteDB() IVideoUserMappingDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IVideoUserMappingDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IVideoUserMappingDo
	Not(conds ...gen.Condition) IVideoUserMappingDo
	Or(conds ...gen.Condition) IVideoUserMappingDo
	Select(conds ...field.Expr) IVideoUserMappingDo
	Where(conds ...gen.Condition) IVideoUserMappingDo
	Order(conds ...field.Expr) IVideoUserMappingDo
	Distinct(cols ...field.Expr) IVideoUserMappingDo
	Omit(cols ...field.Expr) IVideoUserMappingDo
	Join(table schema.Tabler, on ...field.Expr) IVideoUserMappingDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IVideoUserMappingDo
	RightJoin(table schema.Tabler, on ...field.Expr) IVideoUserMappingDo
	Group(cols ...field.Expr) IVideoUserMappingDo
	Having(conds ...gen.Condition) IVideoUserMappingDo
	Limit(limit int) IVideoUserMappingDo
	Offset(offset int) IVideoUserMappingDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IVideoUserMappingDo
	Unscoped() IVideoUserMappingDo
	Create(values ...*model.VideoUserMapping) error
	CreateInBatches(values []*model.VideoUserMapping, batchSize int) error
	Save(values ...*model.VideoUserMapping) error
	First() (*model.VideoUserMapping, error)
	Take() (*model.VideoUserMapping, error)
	Last() (*model.VideoUserMapping, error)
	Find() ([]*model.VideoUserMapping, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.VideoUserMapping, err error)
	FindInBatches(result *[]*model.VideoUserMapping, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.VideoUserMapping) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IVideoUserMappingDo
	Assign(attrs ...field.AssignExpr) IVideoUserMappingDo
	Joins(fields ...field.RelationField) IVideoUserMappingDo
	Preload(fields ...field.RelationField) IVideoUserMappingDo
	FirstOrInit() (*model.VideoUserMapping, error)
	FirstOrCreate() (*model.VideoUserMapping, error)
	FindByPage(offset int, limit int) (result []*model.VideoUserMapping, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IVideoUserMappingDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	SelectByID(id int64) (result *model.VideoUserMapping, err error)
}

// SELECT * FROM @@table WHERE id=@id
func (v videoUserMappingDo) SelectByID(id int64) (result *model.VideoUserMapping, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("SELECT * FROM video_user_mapping WHERE id=? ")

	var executeSQL *gorm.DB
	executeSQL = v.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (v videoUserMappingDo) Debug() IVideoUserMappingDo {
	return v.withDO(v.DO.Debug())
}

func (v videoUserMappingDo) WithContext(ctx context.Context) IVideoUserMappingDo {
	return v.withDO(v.DO.WithContext(ctx))
}

func (v videoUserMappingDo) ReadDB() IVideoUserMappingDo {
	return v.Clauses(dbresolver.Read)
}

func (v videoUserMappingDo) WriteDB() IVideoUserMappingDo {
	return v.Clauses(dbresolver.Write)
}

func (v videoUserMappingDo) Session(config *gorm.Session) IVideoUserMappingDo {
	return v.withDO(v.DO.Session(config))
}

func (v videoUserMappingDo) Clauses(conds ...clause.Expression) IVideoUserMappingDo {
	return v.withDO(v.DO.Clauses(conds...))
}

func (v videoUserMappingDo) Returning(value interface{}, columns ...string) IVideoUserMappingDo {
	return v.withDO(v.DO.Returning(value, columns...))
}

func (v videoUserMappingDo) Not(conds ...gen.Condition) IVideoUserMappingDo {
	return v.withDO(v.DO.Not(conds...))
}

func (v videoUserMappingDo) Or(conds ...gen.Condition) IVideoUserMappingDo {
	return v.withDO(v.DO.Or(conds...))
}

func (v videoUserMappingDo) Select(conds ...field.Expr) IVideoUserMappingDo {
	return v.withDO(v.DO.Select(conds...))
}

func (v videoUserMappingDo) Where(conds ...gen.Condition) IVideoUserMappingDo {
	return v.withDO(v.DO.Where(conds...))
}

func (v videoUserMappingDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IVideoUserMappingDo {
	return v.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (v videoUserMappingDo) Order(conds ...field.Expr) IVideoUserMappingDo {
	return v.withDO(v.DO.Order(conds...))
}

func (v videoUserMappingDo) Distinct(cols ...field.Expr) IVideoUserMappingDo {
	return v.withDO(v.DO.Distinct(cols...))
}

func (v videoUserMappingDo) Omit(cols ...field.Expr) IVideoUserMappingDo {
	return v.withDO(v.DO.Omit(cols...))
}

func (v videoUserMappingDo) Join(table schema.Tabler, on ...field.Expr) IVideoUserMappingDo {
	return v.withDO(v.DO.Join(table, on...))
}

func (v videoUserMappingDo) LeftJoin(table schema.Tabler, on ...field.Expr) IVideoUserMappingDo {
	return v.withDO(v.DO.LeftJoin(table, on...))
}

func (v videoUserMappingDo) RightJoin(table schema.Tabler, on ...field.Expr) IVideoUserMappingDo {
	return v.withDO(v.DO.RightJoin(table, on...))
}

func (v videoUserMappingDo) Group(cols ...field.Expr) IVideoUserMappingDo {
	return v.withDO(v.DO.Group(cols...))
}

func (v videoUserMappingDo) Having(conds ...gen.Condition) IVideoUserMappingDo {
	return v.withDO(v.DO.Having(conds...))
}

func (v videoUserMappingDo) Limit(limit int) IVideoUserMappingDo {
	return v.withDO(v.DO.Limit(limit))
}

func (v videoUserMappingDo) Offset(offset int) IVideoUserMappingDo {
	return v.withDO(v.DO.Offset(offset))
}

func (v videoUserMappingDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IVideoUserMappingDo {
	return v.withDO(v.DO.Scopes(funcs...))
}

func (v videoUserMappingDo) Unscoped() IVideoUserMappingDo {
	return v.withDO(v.DO.Unscoped())
}

func (v videoUserMappingDo) Create(values ...*model.VideoUserMapping) error {
	if len(values) == 0 {
		return nil
	}
	return v.DO.Create(values)
}

func (v videoUserMappingDo) CreateInBatches(values []*model.VideoUserMapping, batchSize int) error {
	return v.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (v videoUserMappingDo) Save(values ...*model.VideoUserMapping) error {
	if len(values) == 0 {
		return nil
	}
	return v.DO.Save(values)
}

func (v videoUserMappingDo) First() (*model.VideoUserMapping, error) {
	if result, err := v.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.VideoUserMapping), nil
	}
}

func (v videoUserMappingDo) Take() (*model.VideoUserMapping, error) {
	if result, err := v.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.VideoUserMapping), nil
	}
}

func (v videoUserMappingDo) Last() (*model.VideoUserMapping, error) {
	if result, err := v.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.VideoUserMapping), nil
	}
}

func (v videoUserMappingDo) Find() ([]*model.VideoUserMapping, error) {
	result, err := v.DO.Find()
	return result.([]*model.VideoUserMapping), err
}

func (v videoUserMappingDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.VideoUserMapping, err error) {
	buf := make([]*model.VideoUserMapping, 0, batchSize)
	err = v.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (v videoUserMappingDo) FindInBatches(result *[]*model.VideoUserMapping, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return v.DO.FindInBatches(result, batchSize, fc)
}

func (v videoUserMappingDo) Attrs(attrs ...field.AssignExpr) IVideoUserMappingDo {
	return v.withDO(v.DO.Attrs(attrs...))
}

func (v videoUserMappingDo) Assign(attrs ...field.AssignExpr) IVideoUserMappingDo {
	return v.withDO(v.DO.Assign(attrs...))
}

func (v videoUserMappingDo) Joins(fields ...field.RelationField) IVideoUserMappingDo {
	for _, _f := range fields {
		v = *v.withDO(v.DO.Joins(_f))
	}
	return &v
}

func (v videoUserMappingDo) Preload(fields ...field.RelationField) IVideoUserMappingDo {
	for _, _f := range fields {
		v = *v.withDO(v.DO.Preload(_f))
	}
	return &v
}

func (v videoUserMappingDo) FirstOrInit() (*model.VideoUserMapping, error) {
	if result, err := v.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.VideoUserMapping), nil
	}
}

func (v videoUserMappingDo) FirstOrCreate() (*model.VideoUserMapping, error) {
	if result, err := v.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.VideoUserMapping), nil
	}
}

func (v videoUserMappingDo) FindByPage(offset int, limit int) (result []*model.VideoUserMapping, count int64, err error) {
	result, err = v.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = v.Offset(-1).Limit(-1).Count()
	return
}

func (v videoUserMappingDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = v.Count()
	if err != nil {
		return
	}

	err = v.Offset(offset).Limit(limit).Scan(result)
	return
}

func (v videoUserMappingDo) Scan(result interface{}) (err error) {
	return v.DO.Scan(result)
}

func (v videoUserMappingDo) Delete(models ...*model.VideoUserMapping) (result gen.ResultInfo, err error) {
	return v.DO.Delete(models)
}

func (v *videoUserMappingDo) withDO(do gen.Dao) *videoUserMappingDo {
	v.DO = *do.(*gen.DO)
	return v
}