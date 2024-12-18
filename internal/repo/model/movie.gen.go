// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameMovie = "movie"

// Movie mapped from table <movie>
type Movie struct {
	ID            int64      `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement:true;comment:主键" json:"id"`
	ExternalID    *string    `gorm:"column:external_id;type:varchar(45);comment:外部id" json:"externalId"`
	OriginalTitle string     `gorm:"column:original_title;type:varchar(1024);not null;comment:原名" json:"originalTitle"`
	Status        string     `gorm:"column:status;type:varchar(64);not null;comment:状态，Returning Series, Ended, Released, Unknown" json:"status"`
	VoteAverage   *float32   `gorm:"column:vote_average;type:float;comment:平均评分" json:"voteAverage"`
	VoteCount     *int32     `gorm:"column:vote_count;type:int;comment:评分数" json:"voteCount"`
	Country       *string    `gorm:"column:country;type:varchar(32);comment:国家" json:"country"`
	Trailer       *string    `gorm:"column:trailer;type:varchar(1024);comment:预告片地址" json:"trailer"`
	URL           string     `gorm:"column:url;type:varchar(2048);not null;comment:影片地址" json:"url"`
	Downloaded    bool       `gorm:"column:downloaded;type:tinyint(1);not null;comment:是否可以下载" json:"downloaded"`
	FileSize      *int32     `gorm:"column:file_size;type:int;comment:文件大小" json:"fileSize"`
	Filename      *string    `gorm:"column:filename;type:varchar(256);comment:文件名" json:"filename"`
	Ext           *string    `gorm:"column:ext;type:varchar(1024);comment:扩展参数" json:"ext"`
	Platform      *string    `gorm:"column:platform;type:varchar(45);comment:1=i4k" json:"platform"`
	Year          *string    `gorm:"column:year;type:varchar(45);comment:年份" json:"year"`
	Definition    *string    `gorm:"column:definition;type:varchar(45);comment:清晰度（1=720p,2=1080P，3=4k）" json:"definition"`
	Promotional   *string    `gorm:"column:promotional;type:varchar(2048);comment:封面地址" json:"promotional"`
	CreateTime    time.Time  `gorm:"column:create_time;type:int unsigned;not null;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdateTime    time.Time  `gorm:"column:update_time;type:int unsigned;not null;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Title         *string    `gorm:"column:title;type:varchar(1024);comment:标题" json:"title"`
	Poster        *string    `gorm:"column:poster;type:varchar(1024);comment:海报" json:"poster"`
	Logo          *string    `gorm:"column:logo;type:varchar(1024);comment:logo" json:"logo"`
	AirDate       *time.Time `gorm:"column:air_date;type:datetime;comment:开播时间" json:"airDate"`
	Overview      *string    `gorm:"column:overview;type:varchar(2048);comment:简介" json:"overview"`
}

// TableName Movie's table name
func (*Movie) TableName() string {
	return TableNameMovie
}
