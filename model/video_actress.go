package model

import "gorm.io/gorm"

// VideoActress 视频演员表
type VideoActress struct {
	gorm.Model
	VideoId   uint `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	ActressId uint `gorm:"column:actress_id;type:uint;not null;default:0;comment:演员ID"`
}
