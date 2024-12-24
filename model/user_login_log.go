package model

import (
	"time"

	"gorm.io/gorm"
)

// UserLoginLog 用户登陆日志表
type UserLoginLog struct {
	gorm.Model
	UserID        uint      `gorm:"column:user_id;type:uint;not null;default:0;comment:用户ID" json:"user_id,omitempty"`
	LoginName     string    `gorm:"comment:登录账号" json:"login_name,omitempty"`
	LoginIpaddr   string    `gorm:"comment:登录IP地址" json:"login_ipaddr,omitempty"`
	LoginLocation string    `gorm:"comment:登录地点" json:"login_location,omitempty"`
	Browser       string    `gorm:"comment:浏览器类型" json:"browser,omitempty"`
	Os            string    `gorm:"comment:操作系统" json:"os,omitempty"`
	Status        uint8     `gorm:"comment:登录状态(0成功 1失败)" json:"status,omitempty"`
	LoginTime     time.Time `gorm:"comment:登录时间" json:"login_time,omitempty"`
}
