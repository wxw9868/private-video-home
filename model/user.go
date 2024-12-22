package model

import (
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
