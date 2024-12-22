package model

import "gorm.io/gorm"

// UserCollectLog 用户收藏日志表
type UserCollectLog struct {
	gorm.Model
	UserID  uint `gorm:"column:user_id;type:uint;not null;default:0;comment:用户ID"`
	VideoID uint `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
}
