package model

import "gorm.io/gorm"

// VideoLog 视频日志信息表
type VideoLog struct {
	gorm.Model
	VideoID          uint `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	CollectionVolume uint `gorm:"column:collection_volume;type:uint;not null;default:0;comment:收藏量"`
	PageViews        uint `gorm:"column:page_views;type:uint;not null;default:0;comment:浏览量"`
	LikesCount       uint `gorm:"column:likes_count;type:uint;not null;default:0;comment:点赞量"`
	DislikesCount    uint `gorm:"column:dislikes_count;type:uint;not null;default:0;comment:点踩量"`
	ViewsCount       uint `gorm:"column:views_count;type:uint;not null;default:0;comment:观看次数"`
}
