package model

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Title         string    `gorm:"column:title;type:varchar(255);comment:标题"`
	Actress       string    `gorm:"column:actress;type:varchar(100);comment:演员"`
	Size          int64     `gorm:"column:size;type:bigint;comment:大小"`
	Duration      float64   `gorm:"column:duration;type:float;default:0;comment:时长"`
	Poster        string    `gorm:"column:poster;type:varchar(255);comment:封面"`
	Width         int       `gorm:"column:width;type:int;default:0;comment:宽"`
	Height        int       `gorm:"column:height;type:int;default:0;comment:高"`
	CodecName     string    `gorm:"column:codec_name;type:varchar(90);comment:编解码器"`
	ChannelLayout string    `gorm:"column:channel_layout;type:varchar(90);comment:音频声道"`
	CreationTime  time.Time `gorm:"column:creation_time;type:date;comment:时间"`
}

type Actress struct {
	gorm.Model
	Actress string `gorm:"column:actress;type:varchar(100);comment:演员"`
	Avatar  string `gorm:"column:avatar;type:varchar(255);comment:头像"`
}

type VideoLog struct {
	gorm.Model
	VideoID uint `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	Collect uint `gorm:"column:collect;type:uint;not null;default:0;comment:收藏"`
	Browse  uint `gorm:"column:browse;type:uint;not null;default:0;comment:浏览"`
	Zan     uint `gorm:"column:zan;type:uint;not null;default:0;comment:赞"`
	Cai     uint `gorm:"column:cai;type:uint;not null;default:0;comment:踩"`
	Watch   uint `gorm:"column:watch;type:uint;not null;default:0;comment:观看"`
}
