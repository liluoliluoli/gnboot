// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameGenre = "genre"

// Genre 流派
type Genre struct {
	ID   uint64 `gorm:"column:id;primaryKey;autoIncrement:true;comment:主键" json:"id,string"` // 主键
	Name string `gorm:"column:name;not null;comment:名称" json:"name"`                         // 名称
}

// TableName Genre's table name
func (*Genre) TableName() string {
	return TableNameGenre
}