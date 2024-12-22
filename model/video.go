package model

import (
	"time"

	"gorm.io/gorm"
)

// Video 视频表
type Video struct {
	gorm.Model
	CateID        uint      `gorm:"column:cate_id;type:uint;not null;default:0;comment:视频分类ID"`
	Title         string    `gorm:"column:title;type:varchar(255);uniqueIndex;comment:标题"`
	Actress       string    `gorm:"column:actress;type:varchar(100);comment:演员"`
	Size          int64     `gorm:"column:size;type:bigint;comment:大小"`
	Duration      float64   `gorm:"column:duration;type:float;default:0;comment:时长"`
	Poster        string    `gorm:"column:poster;type:varchar(255);comment:封面"`
	Width         int       `gorm:"column:width;type:int;default:0;comment:宽"`
	Height        int       `gorm:"column:height;type:int;default:0;comment:高"`
	CodecName     string    `gorm:"column:codec_name;type:varchar(90);comment:编解码器"`
	ChannelLayout string    `gorm:"column:channel_layout;type:varchar(90);comment:音频声道"`
	CreationTime  time.Time `gorm:"column:creation_time;type:date;comment:时间"`
	Intro         string    `gorm:"column:intro;type:text;comment:简介"`
}
