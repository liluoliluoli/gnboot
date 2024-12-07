package main

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

var tableMap = make(map[string]string)

func init() {
	tableMap["actor"] = "Actor"
	tableMap["app_version"] = "AppVersion"
	tableMap["episode"] = "Episode"
	tableMap["genre"] = "Genre"
	tableMap["keyword"] = "Keyword"
	tableMap["movie"] = "Movie"
	tableMap["season"] = "Season"
	tableMap["series"] = "Series"
	tableMap["studio"] = "Studio"
	tableMap["user"] = "User"
	tableMap["video_actor_mapping"] = "VideoActorMapping"
	tableMap["video_genre_mapping"] = "VideoGenreMapping"
	tableMap["video_keyword_mapping"] = "VideoKeywordMapping"
	tableMap["video_studio_mapping"] = "VideoStudioMapping"
	tableMap["video_subtitle_mapping"] = "VideoSubtitleMapping"
	tableMap["video_user_mapping"] = "VideoUserMapping"
}

func main() {
	// 加载配置
	dsn := "gnvideo:7H5rdOAIMA815pXz@tcp(mysql.sqlpub.com:3306)/gnvideo?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local&timeout=10000ms"

	// 初始化DB连接
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Errorf("cannot establish mysql connection: %w", err))
	}

	// 初始化代码生成器实例
	g := gen.NewGenerator(gen.Config{
		OutPath:      "../internal/repo/gen", // 生成的DAL类文件系统路径
		ModelPkgPath: "model",                // DataObject模型包名
		// WithDefaultQuery 生成默认查询结构体(作为全局变量使用), 即`Q`结构体和其字段(各表模型)
		// WithoutContext 生成没有context调用限制的代码供查询
		// WithQueryInterface 生成interface形式的查询代码(可导出), 如`Where()`方法返回的就是一个可导出的接口类型
		Mode: gen.WithDefaultQuery | gen.WithQueryInterface,
		// 表字段可为 null 值时, 对应结体字段使用指针类型
		FieldNullable: true,
		// 表字段默认值与模型结构体字段零值不一致的字段, 在插入数据时需要赋值该字段值为零值的, 结构体字段须是指针类型才能成功, 即`FieldCoverable:true`配置下生成的结构体字段.
		// 因为在插入时遇到字段为零值的会被GORM赋予默认值. 如字段`age`表默认值为10, 即使你显式设置为0最后也会被GORM设为10提交.
		// 如果该字段没有上面提到的插入时赋零值的特殊需要, 则字段为非指针类型使用起来会比较方便.
		FieldCoverable: false,
		// 模型结构体字段的数字类型的符号表示是否与表字段的一致, `false`指示都用有符号类型
		FieldSignable: false,
		// 生成 gorm 标签的字段索引属性
		FieldWithIndexTag: false,
		// 生成 gorm 标签的字段类型属性
		FieldWithTypeTag: true,
	})

	g.WithJSONTagNameStrategy(func(c string) string {
		return toCamelCase(c)
	})

	g.UseDB(db)

	// 通用配置
	commonOpts := commonFieldOpts()

	// 配置业务数据表
	models := addBizTable(g, commonOpts)
	//g.ApplyBasic(models...)
	g.ApplyInterface(func(Querier) {}, models...)
	// 生成代码
	g.Execute()
}

type Querier interface {
	// SELECT * FROM @@table WHERE id=@id
	SelectByID(id int64) (*gen.T, error) // GetByID query data by id and return it as *struct*
}

// 将下划线命名转换为驼峰命名
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	result := strings.Join(parts, "")
	result = strings.ToLower(result[:1]) + result[1:]
	return result
}

func commonFieldOpts() []gen.ModelOpt {
	// 自动更新时间戳
	autoUpdateTimeField := gen.FieldGORMTag("update_time", func(tag field.GormTag) field.GormTag {
		tag.Set(field.TagKeyGormColumn, "update_time")
		tag.Set(field.TagKeyGormType, "int unsigned")
		tag.Set("autoUpdateTime", "")
		return tag
	})
	autoCreateTimeField := gen.FieldGORMTag("create_time", func(tag field.GormTag) field.GormTag {
		tag.Set(field.TagKeyGormColumn, "create_time")
		tag.Set(field.TagKeyGormType, "int unsigned")
		tag.Set("autoCreateTime", "")
		return tag
	})

	return []gen.ModelOpt{
		autoUpdateTimeField,
		autoCreateTimeField,
	}
}

func addBizTable(g *gen.Generator, commonOpts []gen.ModelOpt) []interface{} {
	models := make([]interface{}, 32)

	for table, entity := range tableMap {
		// 书本
		model := g.GenerateModelAs(table, entity,
			append(
				commonOpts,
				// 在下面添加其他配置
			)...,
		)
		// 这里添加自定义SQL
		//g.ApplyInterface(func(bookMapper) {}, book)
		models = append(models, model)
	}
	return models
}
