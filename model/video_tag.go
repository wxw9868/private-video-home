package model

import "gorm.io/gorm"

// VideoTag 视频标签表
type VideoTag struct {
	gorm.Model
	VideoId uint `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	TagId   uint `gorm:"column:tag_id;type:uint;not null;default:0;comment:标签ID"`
}
