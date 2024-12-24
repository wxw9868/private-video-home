package model

import "gorm.io/gorm"

// UserCommentLog 用户评论日志表
type UserCommentLog struct {
	gorm.Model
	UserID    uint  `gorm:"column:user_id;type:uint;not null;default:0;comment:用户ID"`
	VideoID   uint  `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	CommentID uint  `gorm:"column:comment_id;type:uint;not null;default:0;comment:评论ID"`
	Like      *int8 `gorm:"column:like;not null;default:0;comment:支持（赞）"`
	Dislike   *int8 `gorm:"column:dislike;not null;default:0;comment:反对（踩）"`
}
