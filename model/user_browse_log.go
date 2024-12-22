package model

import "gorm.io/gorm"

// UserBrowseLog 用户浏览日志表
type UserBrowseLog struct {
	gorm.Model
	UserID  uint `gorm:"column:user_id;type:uint;not null;default:0;comment:用户ID"`
	VideoID uint `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	Number  uint `gorm:"column:number;type:uint;not null;default:0;comment:数"`
}
