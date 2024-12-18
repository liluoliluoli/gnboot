// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameVideoGenreMapping = "video_genre_mapping"

// VideoGenreMapping mapped from table <video_genre_mapping>
type VideoGenreMapping struct {
	ID         int64     `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement:true;comment:主键" json:"id"`
	VideoType  string    `gorm:"column:video_type;type:varchar(32);not null;comment:影片类型，movie,series,season,episode" json:"videoType"`
	VideoID    int64     `gorm:"column:video_id;type:bigint;not null;comment:影片id，根据video_type类型分别来自movie,series,season,episode表" json:"videoId"`
	GenreID    int64     `gorm:"column:genre_id;type:bigint;not null;comment:流派id" json:"genreId"`
	CreateTime time.Time `gorm:"column:create_time;type:int unsigned;not null;autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;type:int unsigned;not null;autoUpdateTime" json:"updateTime"`
}

// TableName VideoGenreMapping's table name
func (*VideoGenreMapping) TableName() string {
	return TableNameVideoGenreMapping
}
