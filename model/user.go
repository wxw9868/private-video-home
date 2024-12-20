package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	gorm.Model
	Username    string `gorm:"column:username;type:varchar(100);uniqueIndex;comment:账户名"`
	Password    string `gorm:"column:password;type:varchar(100);comment:账户密码"`
	Nickname    string `gorm:"column:nickname;type:varchar(30);comment:昵称"`
	Sex         uint8  `gorm:"column:sex;type:uint;default:0;comment:性别 0保密 1男 2女"`
	Mobile      string `gorm:"column:mobile;type:string;comment:手机号"`
	Email       string `gorm:"column:email;type:string;comment:邮箱地址"`
	QQ          string `gorm:"column:qq;type:varchar(20);comment:QQ"`
	Avatar      string `gorm:"column:avatar;type:string;comment:用户头像地址"`
	Designation string `gorm:"column:designation;type:string;comment:称号"`
	Realname    string `gorm:"column:realname;type:varchar(50);comment:真实姓名"`
	Idcard      string `gorm:"column:idcard;type:varchar(100);comment:身份证号"`
	Address     string `gorm:"column:address;type:string;comment:地址"`
	Intro       string `gorm:"column:intro;type:string;comment:简介"`
}

// UserLoginLog 用户登陆日志表
type UserLoginLog struct {
	gorm.Model
	UserID        uint      `gorm:"column:user_id;type:uint;not null;default:0;comment:用户ID"`
	LoginName     string    `gorm:"comment:登录账号"`
	LoginIpaddr   string    `gorm:"comment:登录IP地址"`
	LoginLocation string    `gorm:"comment:登录地点"`
	Browser       string    `gorm:"comment:浏览器类型"`
	Os            string    `gorm:"comment:操作系统"`
	Status        uint8     `gorm:"comment:登录状态(0成功 1失败)"`
	LoginTime     time.Time `gorm:"comment:登录时间"`
}

// UserCollectLog 用户收藏日志表
type UserCollectLog struct {
	gorm.Model
	UserID  uint `gorm:"column:user_id;type:uint;not null;default:0;comment:用户ID"`
	VideoID uint `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
}

// UserBrowseLog 用户浏览日志表
type UserBrowseLog struct {
	gorm.Model
	UserID  uint `gorm:"column:user_id;type:uint;not null;default:0;comment:用户ID"`
	VideoID uint `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	Number  uint `gorm:"column:number;type:uint;not null;default:0;comment:数"`
}

// UserCommentLog 用户评论日志表
type UserCommentLog struct {
	gorm.Model
	UserID    uint  `gorm:"column:user_id;type:uint;not null;default:0;comment:用户ID"`
	VideoID   uint  `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	CommentID uint  `gorm:"column:comment_id;type:uint;not null;default:0;comment:评论ID"`
	Support   *uint `gorm:"column:support;type:uint;not null;default:0;comment:支持（赞）"`
	Oppose    *uint `gorm:"column:oppose;type:uint;not null;default:0;comment:反对（踩）"`
}
