// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameUser = "user"

// User 用户
type User struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement:true;comment:主键" json:"id,string"` // 主键
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}