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

	"github.com/liluoliluoli/gnboot/internal/repo/model"
)

func newAppVersion(db *gorm.DB, opts ...gen.DOOption) appVersion {
	_appVersion := appVersion{}

	_appVersion.appVersionDo.UseDB(db, opts...)
	_appVersion.appVersionDo.UseModel(&model.AppVersion{})

	tableName := _appVersion.appVersionDo.TableName()
	_appVersion.ALL = field.NewAsterisk(tableName)
	_appVersion.ID = field.NewInt64(tableName, "id")
	_appVersion.VersionCode = field.NewString(tableName, "version_code")
	_appVersion.VersionName = field.NewString(tableName, "version_name")
	_appVersion.PublishTime = field.NewTime(tableName, "publish_time")
	_appVersion.Remark = field.NewString(tableName, "remark")
	_appVersion.ForceUpdate = field.NewBool(tableName, "force_update")
	_appVersion.ApkURL = field.NewString(tableName, "apk_url")

	_appVersion.fillFieldMap()

	return _appVersion
}

type appVersion struct {
	appVersionDo appVersionDo

	ALL         field.Asterisk
	ID          field.Int64  // 主键
	VersionCode field.String // 版本数字代号，100
	VersionName field.String // 版本名称，1.0.0
	PublishTime field.Time   // 发布时间
	Remark      field.String // 备注
	ForceUpdate field.Bool   // 是否强制更新
	ApkURL      field.String // 下载地址

	fieldMap map[string]field.Expr
}

func (a appVersion) Table(newTableName string) *appVersion {
	a.appVersionDo.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a appVersion) As(alias string) *appVersion {
	a.appVersionDo.DO = *(a.appVersionDo.As(alias).(*gen.DO))
	return a.updateTableName(alias)
}

func (a *appVersion) updateTableName(table string) *appVersion {
	a.ALL = field.NewAsterisk(table)
	a.ID = field.NewInt64(table, "id")
	a.VersionCode = field.NewString(table, "version_code")
	a.VersionName = field.NewString(table, "version_name")
	a.PublishTime = field.NewTime(table, "publish_time")
	a.Remark = field.NewString(table, "remark")
	a.ForceUpdate = field.NewBool(table, "force_update")
	a.ApkURL = field.NewString(table, "apk_url")

	a.fillFieldMap()

	return a
}

func (a *appVersion) WithContext(ctx context.Context) IAppVersionDo {
	return a.appVersionDo.WithContext(ctx)
}

func (a appVersion) TableName() string { return a.appVersionDo.TableName() }

func (a appVersion) Alias() string { return a.appVersionDo.Alias() }

func (a *appVersion) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *appVersion) fillFieldMap() {
	a.fieldMap = make(map[string]field.Expr, 7)
	a.fieldMap["id"] = a.ID
	a.fieldMap["version_code"] = a.VersionCode
	a.fieldMap["version_name"] = a.VersionName
	a.fieldMap["publish_time"] = a.PublishTime
	a.fieldMap["remark"] = a.Remark
	a.fieldMap["force_update"] = a.ForceUpdate
	a.fieldMap["apk_url"] = a.ApkURL
}

func (a appVersion) clone(db *gorm.DB) appVersion {
	a.appVersionDo.ReplaceConnPool(db.Statement.ConnPool)
	return a
}

func (a appVersion) replaceDB(db *gorm.DB) appVersion {
	a.appVersionDo.ReplaceDB(db)
	return a
}

type appVersionDo struct{ gen.DO }

type IAppVersionDo interface {
	gen.SubQuery
	Debug() IAppVersionDo
	WithContext(ctx context.Context) IAppVersionDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IAppVersionDo
	WriteDB() IAppVersionDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IAppVersionDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IAppVersionDo
	Not(conds ...gen.Condition) IAppVersionDo
	Or(conds ...gen.Condition) IAppVersionDo
	Select(conds ...field.Expr) IAppVersionDo
	Where(conds ...gen.Condition) IAppVersionDo
	Order(conds ...field.Expr) IAppVersionDo
	Distinct(cols ...field.Expr) IAppVersionDo
	Omit(cols ...field.Expr) IAppVersionDo
	Join(table schema.Tabler, on ...field.Expr) IAppVersionDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IAppVersionDo
	RightJoin(table schema.Tabler, on ...field.Expr) IAppVersionDo
	Group(cols ...field.Expr) IAppVersionDo
	Having(conds ...gen.Condition) IAppVersionDo
	Limit(limit int) IAppVersionDo
	Offset(offset int) IAppVersionDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IAppVersionDo
	Unscoped() IAppVersionDo
	Create(values ...*model.AppVersion) error
	CreateInBatches(values []*model.AppVersion, batchSize int) error
	Save(values ...*model.AppVersion) error
	First() (*model.AppVersion, error)
	Take() (*model.AppVersion, error)
	Last() (*model.AppVersion, error)
	Find() ([]*model.AppVersion, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.AppVersion, err error)
	FindInBatches(result *[]*model.AppVersion, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.AppVersion) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IAppVersionDo
	Assign(attrs ...field.AssignExpr) IAppVersionDo
	Joins(fields ...field.RelationField) IAppVersionDo
	Preload(fields ...field.RelationField) IAppVersionDo
	FirstOrInit() (*model.AppVersion, error)
	FirstOrCreate() (*model.AppVersion, error)
	FindByPage(offset int, limit int) (result []*model.AppVersion, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IAppVersionDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	SelectByID(id int64) (result *model.AppVersion, err error)
}

// SELECT * FROM @@table WHERE id=@id
func (a appVersionDo) SelectByID(id int64) (result *model.AppVersion, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("SELECT * FROM app_version WHERE id=? ")

	var executeSQL *gorm.DB
	executeSQL = a.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (a appVersionDo) Debug() IAppVersionDo {
	return a.withDO(a.DO.Debug())
}

func (a appVersionDo) WithContext(ctx context.Context) IAppVersionDo {
	return a.withDO(a.DO.WithContext(ctx))
}

func (a appVersionDo) ReadDB() IAppVersionDo {
	return a.Clauses(dbresolver.Read)
}

func (a appVersionDo) WriteDB() IAppVersionDo {
	return a.Clauses(dbresolver.Write)
}

func (a appVersionDo) Session(config *gorm.Session) IAppVersionDo {
	return a.withDO(a.DO.Session(config))
}

func (a appVersionDo) Clauses(conds ...clause.Expression) IAppVersionDo {
	return a.withDO(a.DO.Clauses(conds...))
}

func (a appVersionDo) Returning(value interface{}, columns ...string) IAppVersionDo {
	return a.withDO(a.DO.Returning(value, columns...))
}

func (a appVersionDo) Not(conds ...gen.Condition) IAppVersionDo {
	return a.withDO(a.DO.Not(conds...))
}

func (a appVersionDo) Or(conds ...gen.Condition) IAppVersionDo {
	return a.withDO(a.DO.Or(conds...))
}

func (a appVersionDo) Select(conds ...field.Expr) IAppVersionDo {
	return a.withDO(a.DO.Select(conds...))
}

func (a appVersionDo) Where(conds ...gen.Condition) IAppVersionDo {
	return a.withDO(a.DO.Where(conds...))
}

func (a appVersionDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IAppVersionDo {
	return a.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (a appVersionDo) Order(conds ...field.Expr) IAppVersionDo {
	return a.withDO(a.DO.Order(conds...))
}

func (a appVersionDo) Distinct(cols ...field.Expr) IAppVersionDo {
	return a.withDO(a.DO.Distinct(cols...))
}

func (a appVersionDo) Omit(cols ...field.Expr) IAppVersionDo {
	return a.withDO(a.DO.Omit(cols...))
}

func (a appVersionDo) Join(table schema.Tabler, on ...field.Expr) IAppVersionDo {
	return a.withDO(a.DO.Join(table, on...))
}

func (a appVersionDo) LeftJoin(table schema.Tabler, on ...field.Expr) IAppVersionDo {
	return a.withDO(a.DO.LeftJoin(table, on...))
}

func (a appVersionDo) RightJoin(table schema.Tabler, on ...field.Expr) IAppVersionDo {
	return a.withDO(a.DO.RightJoin(table, on...))
}

func (a appVersionDo) Group(cols ...field.Expr) IAppVersionDo {
	return a.withDO(a.DO.Group(cols...))
}

func (a appVersionDo) Having(conds ...gen.Condition) IAppVersionDo {
	return a.withDO(a.DO.Having(conds...))
}

func (a appVersionDo) Limit(limit int) IAppVersionDo {
	return a.withDO(a.DO.Limit(limit))
}

func (a appVersionDo) Offset(offset int) IAppVersionDo {
	return a.withDO(a.DO.Offset(offset))
}

func (a appVersionDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IAppVersionDo {
	return a.withDO(a.DO.Scopes(funcs...))
}

func (a appVersionDo) Unscoped() IAppVersionDo {
	return a.withDO(a.DO.Unscoped())
}

func (a appVersionDo) Create(values ...*model.AppVersion) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Create(values)
}

func (a appVersionDo) CreateInBatches(values []*model.AppVersion, batchSize int) error {
	return a.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (a appVersionDo) Save(values ...*model.AppVersion) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Save(values)
}

func (a appVersionDo) First() (*model.AppVersion, error) {
	if result, err := a.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.AppVersion), nil
	}
}

func (a appVersionDo) Take() (*model.AppVersion, error) {
	if result, err := a.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.AppVersion), nil
	}
}

func (a appVersionDo) Last() (*model.AppVersion, error) {
	if result, err := a.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.AppVersion), nil
	}
}

func (a appVersionDo) Find() ([]*model.AppVersion, error) {
	result, err := a.DO.Find()
	return result.([]*model.AppVersion), err
}

func (a appVersionDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.AppVersion, err error) {
	buf := make([]*model.AppVersion, 0, batchSize)
	err = a.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (a appVersionDo) FindInBatches(result *[]*model.AppVersion, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return a.DO.FindInBatches(result, batchSize, fc)
}

func (a appVersionDo) Attrs(attrs ...field.AssignExpr) IAppVersionDo {
	return a.withDO(a.DO.Attrs(attrs...))
}

func (a appVersionDo) Assign(attrs ...field.AssignExpr) IAppVersionDo {
	return a.withDO(a.DO.Assign(attrs...))
}

func (a appVersionDo) Joins(fields ...field.RelationField) IAppVersionDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Joins(_f))
	}
	return &a
}

func (a appVersionDo) Preload(fields ...field.RelationField) IAppVersionDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Preload(_f))
	}
	return &a
}

func (a appVersionDo) FirstOrInit() (*model.AppVersion, error) {
	if result, err := a.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.AppVersion), nil
	}
}

func (a appVersionDo) FirstOrCreate() (*model.AppVersion, error) {
	if result, err := a.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.AppVersion), nil
	}
}

func (a appVersionDo) FindByPage(offset int, limit int) (result []*model.AppVersion, count int64, err error) {
	result, err = a.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = a.Offset(-1).Limit(-1).Count()
	return
}

func (a appVersionDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = a.Count()
	if err != nil {
		return
	}

	err = a.Offset(offset).Limit(limit).Scan(result)
	return
}

func (a appVersionDo) Scan(result interface{}) (err error) {
	return a.DO.Scan(result)
}

func (a appVersionDo) Delete(models ...*model.AppVersion) (result gen.ResultInfo, err error) {
	return a.DO.Delete(models)
}

func (a *appVersionDo) withDO(do gen.Dao) *appVersionDo {
	a.DO = *do.(*gen.DO)
	return a
}
