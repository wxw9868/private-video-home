package model

import "gorm.io/gorm"

// VideoLog 视频日志信息表
type VideoLog struct {
	gorm.Model
	VideoID uint `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	Collect uint `gorm:"column:collect;type:uint;not null;default:0;comment:收藏"`
	Browse  uint `gorm:"column:browse;type:uint;not null;default:0;comment:浏览"`
	Like    uint `gorm:"column:like;type:uint;not null;default:0;comment:赞"`
	Dislike uint `gorm:"column:dislike;type:uint;not null;default:0;comment:踩"`
	Watch   uint `gorm:"column:watch;type:uint;not null;default:0;comment:观看"`
}
