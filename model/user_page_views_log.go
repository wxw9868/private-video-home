package model

import "gorm.io/gorm"

// UserPageViewsLog 用户页面浏览量表
type UserPageViewsLog struct {
	gorm.Model
	UserID    uint `gorm:"column:user_id;type:uint;not null;default:0;comment:用户ID"`
	VideoID   uint `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	PageViews uint `gorm:"column:page_views;type:uint;not null;default:0;comment:页面浏览量"`
}
