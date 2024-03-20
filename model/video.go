package model

import (
	"math/big"
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Title         string    `gorm:"column:title;type:varchar(255);comment:标题"`
	Actress       string    `gorm:"column:actress;type:varchar(100);comment:演员"`
	Size          big.Int   `gorm:"column:size;type:bigint;comment:大小"`
	Duration      int       `gorm:"column:duration;type:int;default:0;comment:时长"`
	ModTime       string    `gorm:"column:mod_time;type:varchar(255);comment:修改时间"`
	Poster        string    `gorm:"column:poster;type:varchar(255);comment:封面"`
	Width         int       `gorm:"column:width;type:int;default:0;comment:宽"`
	Height        int       `gorm:"column:height;type:int;default:0;comment:高"`
	CodecName     string    `gorm:"column:title;type:varchar(90);comment:编解码器"`
	ChannelLayout string    `gorm:"column:channel_layout;type:varchar(90);comment:音频声道"`
	CreationTime  time.Time `gorm:"column:creation_time;type:time;comment:创建时间"`
}

type Actress struct {
	gorm.Model
	Actress string `gorm:"column:actress;type:varchar(100);comment:演员"`
	Avatar  string `gorm:"column:avatar;type:varchar(255);comment:头像"`
}
