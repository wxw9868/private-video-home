package model

import (
	"time"

	"gorm.io/gorm"
)

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
