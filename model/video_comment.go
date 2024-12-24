package model

import "gorm.io/gorm"

// VideoComment 视频评论表
type VideoComment struct {
	gorm.Model
	ParentId      uint   `gorm:"column:parent_id;type:uint;not null;default:0;comment:父级评论的ID"`
	VideoId       uint   `gorm:"column:video_id;type:uint;not null;default:0;comment:被评论的视频ID"`
	UserId        uint   `gorm:"column:user_id;type:uint;not null;default:0;comment:评论人的ID"`
	Nickname      string `gorm:"column:nickname;type:varchar(13);null;comment:评论人的昵称"`
	Avatar        string `gorm:"column:avatar;type:varchar(255);null;comment:评论人的头像地址"`
	Status        string `gorm:"column:status;type:check(status in ('VERIFYING','APPROVED','REJECT','DELETED'));default:'VERIFYING';comment:评论的状态"`
	LikeNum       uint   `gorm:"column:like_num;type:uint;not null;default:0;comment:点赞人数"`
	LikeUserid    string `gorm:"column:like_userid;type:varchar(255);null;comment:"`
	ReplyNum      uint   `gorm:"column:reply_num;type:uint;not null;default:0;comment:评论回复数"`
	IsAnonymous   uint   `gorm:"column:is_anonymous;type:uint;not null;default:0;comment:是否匿名评价 0是 1不是"`
	Content       string `gorm:"column:content;type:text;not null;comment:评论内容"`
	Remark        string `gorm:"column:remark;type:varchar(100);not null;comment:备注（审核不通过时添加）"`
	LikesCount    uint   `gorm:"column:likes_count;type:uint;not null;default:0;comment:点赞量"`
	DislikesCount uint   `gorm:"column:dislikes_count;type:uint;not null;default:0;comment:点踩量"`
	IsShow        uint8  `gorm:"column:is_show;type:uint;not null;default:1;comment:是否显示 0不显示 1显示"`
}
