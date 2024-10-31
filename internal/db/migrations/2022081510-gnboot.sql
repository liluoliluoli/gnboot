-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `actor`
(
    id           bigint unsigned auto_increment comment '主键' primary key,
    name         varchar(128)         null comment '名字',
    original_ame varchar(128)         null comment '原名',
    adult        tinyint(1) default 1 not null comment '是否成年',
    gender       tinyint(1)           not null comment '性别',
    profile      varchar(1024)        null comment '演员的照片或头像的URL'
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment '演职人员';

CREATE TABLE IF NOT EXISTS `app_version`
(
    id                   bigint unsigned auto_increment comment '主键' primary key,
    tag_name             varchar(8)    not null comment '版本名称，1.0.0三位',
    published_at         datetime      not null comment '发布时间',
    body                 text          null comment '备注',
    name                 varchar(256)  not null comment '资源名称',
    browser_download_url varchar(1024) not null comment '下载地址'
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci  comment '应用版本';

CREATE TABLE IF NOT EXISTS `episode`
(
    id          bigint  unsigned auto_increment comment '主键' primary key,
    season_id   bigint            null comment '季id',
    episode     int               not null comment '第几集',
    skip_intro  bigint            null comment '片头跳过秒数',
    skip_ending bigint            null comment '片尾跳过秒数',
    url         varchar(1024)     not null comment '影片地址',
    downloaded  tinyint default 0 not null comment '是否能下载',
    ext         varchar(1024)     null comment '扩展参数',
    file_size   bigint            null comment '文件大小'
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment '集';

CREATE TABLE IF NOT EXISTS `genre`
(
    id   bigint unsigned auto_increment comment '主键' primary key,
    name varchar(32) not null comment '名称'
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment '流派';

CREATE TABLE IF NOT EXISTS `keyword`
(
    id   bigint unsigned auto_increment comment '主键' primary key,
    name varchar(128) not null comment '词名称'
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment '关键词';

CREATE TABLE IF NOT EXISTS `movie`
(
    id             bigint unsigned auto_increment comment '主键' primary key,
    original_title varchar(1024)        not null comment '标题',
    status         varchar(64)          not null comment '状态，Returning Series, Ended, Released, Unknown',
    vote_average   float                null comment '平均评分',
    vote_count     bigint               null comment '评分数',
    country        varchar(32)          null comment '国家',
    trailer        varchar(1024)        null comment '预告片地址',
    url            varchar(2048)        not null comment '影片地址',
    downloaded     tinyint(1) default 0 not null comment '是否可以下载',
    file_size      bigint               null comment '文件大小',
    filename       varchar(256)         null comment '文件名',
    ext            varchar(1024)        null comment '扩展参数',
    platform       varchar(45)          NULL COMMENT '1=i4k',
    year           varchar(45)          NULL COMMENT '年份',
    definition     varchar(45)          NULL COMMENT '清晰度（1=720p,2=1080P，3=4k）',
    promotional    varchar(2048)        NULL COMMENT '封面',
    external       varchar(45)          NULL COMMENT '外部id'
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment '电影';

CREATE TABLE IF NOT EXISTS `season`
(
    id            bigint unsigned auto_increment comment '主键' primary key,
    series_id     bigint       not null comment '连续剧id',
    season        int          not null comment '第几季',
    series_title  varchar(256) not null comment '季名称',
    skip_intro    bigint       null comment '片头跳过秒数',
    skip_ending   bigint       null comment '片尾跳过秒数',
    episode_count int          not null comment '总集数'
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment '季';

CREATE TABLE IF NOT EXISTS `series`
(
    id             bigint unsigned auto_increment comment '主键' primary key,
    original_title varchar(1024) not null comment '标题',
    status         varchar(64)   not null comment '状态，Returning Series, Ended, Released, Unknown',
    vote_average   float         null comment '平均评分',
    vote_count     bigint        null comment '评分数',
    country        varchar(32)   null comment '国家',
    trailer        varchar(1024) null comment '预告片地址',
    skip_intro     bigint        null comment '片头跳过秒数',
    skip_ending    bigint        null comment '片尾跳过秒数',
    file_size      bigint        null comment '文件大小',
    filename       varchar(256)  null comment '文件名'
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment '连续剧';

CREATE TABLE IF NOT EXISTS `studio`
(
    id      bigint unsigned auto_increment comment '主键' primary key,
    name    varchar(32)   not null comment '名称',
    country varchar(32)   not null comment '国家',
    logo    varchar(1024) null comment 'logo的地址'
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment '出品方';

CREATE TABLE IF NOT EXISTS `user`
(
    id bigint auto_increment comment '主键' primary key
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment '用户';

CREATE TABLE IF NOT EXISTS `video_actor_mapping`
(
    id          bigint unsigned auto_increment comment '主键' primary key,
    video_type  varchar(32)  not null comment '影片类型，movie,series,season,episode',
    video_id    bigint       not null comment '影片id，根据video_type类型分别来自movie,series,season,episode表',
    actor_id    bigint       not null comment '演职人员id',
    `character` varchar(128) null comment '饰演角色名称'
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment '影片-演职人员关系表';

CREATE TABLE IF NOT EXISTS `video_genre_mapping`
(
    id         bigint unsigned auto_increment comment '主键' primary key,
    video_type varchar(32) null comment '影片类型，movie,series,season,episode',
    video_id   bigint      null comment '影片id，根据video_type类型分别来自movie,series,season,episode表',
    genre_id   bigint      not null comment '流派id'
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment '影片-流派关联表';

CREATE TABLE IF NOT EXISTS `video_keyword_mapping`
(
    id         bigint unsigned auto_increment comment '主键' primary key,
    video_type varchar(32) not null comment '影片类型，movie,series,season,episode',
    video_id   bigint      not null comment '影片id，根据video_type类型分别来自movie,series,season,episode表',
    keyword_id bigint      not null comment '关键词id'
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment '影片关键词关联表';

CREATE TABLE IF NOT EXISTS `video_studio_mapping`
(
    id         bigint unsigned auto_increment comment '主键' primary key,
    video_type varchar(32) not null comment '影片类型，movie,series,season,episode',
    video_id   bigint      not null comment '影片id，根据video_type类型分别来自movie,series,season,episode表',
    studio_id  bigint      not null comment '出品方id'
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment '影片出品方关联表';

CREATE TABLE IF NOT EXISTS `video_subtitle_mapping`
(
    id         bigint unsigned auto_increment comment '主键' primary key,
    video_type varchar(32)   not null comment '影片类型，movie,series,season,episode',
    video_id   bigint        not null comment '影片id，根据video_type类型分别来自movie,series,season,episode表',
    url        varchar(1024) not null comment '字幕地址',
    title      varchar(1024) not null comment '字幕标题',
    language   varchar(32)   not null comment '字幕语言',
    mime_type  varchar(256)  not null comment '字幕格式'
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment '影片字幕关联表';

CREATE TABLE IF NOT EXISTS `video_user_mapping`
(
    id                   bigint unsigned auto_increment comment '主键' primary key,
    video_type           varchar(32)          not null comment '影片类型，movie,series,season,episode',
    video_id             bigint               not null comment '影片id，根据video_type类型分别来自movie,series,season,episode表',
    last_played_position bigint               null comment '影片上次播放位置，第n秒',
    last_played_time     bigint               null comment '上次播放时候',
    favorited            tinyint(1) default 0 not null comment '是否收藏喜欢'
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci comment '影片用户关联表';


