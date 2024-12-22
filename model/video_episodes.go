package model

import (
	"time"

	"gorm.io/gorm"
)

// VideoEpisodes 视频剧集信息表
type VideoEpisodes struct {
	gorm.Model
	VideoId       uint      `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	Episode       uint      `gorm:"column:episode;type:uint;not null;default:0;comment:剧集"`
	Size          int64     `gorm:"column:size;type:bigint;comment:大小"`
	Duration      float64   `gorm:"column:duration;type:float;default:0;comment:时长"`
	Width         int       `gorm:"column:width;type:int;default:0;comment:宽"`
	Height        int       `gorm:"column:height;type:int;default:0;comment:高"`
	CodecName     string    `gorm:"column:codec_name;type:varchar(90);comment:编解码器"`
	ChannelLayout string    `gorm:"column:channel_layout;type:varchar(90);comment:音频声道"`
	CreationTime  time.Time `gorm:"column:creation_time;type:date;comment:时间"`
}
