package model

import (
	"gorm.io/gorm"
)

// User 用户表
type User struct {
	gorm.Model
	Username    string `gorm:"column:username;type:varchar(100);uniqueIndex;comment:账户名" json:"username,omitempty"`
	Password    string `gorm:"column:password;type:varchar(100);comment:账户密码" json:"password,omitempty"`
	Nickname    string `gorm:"column:nickname;type:varchar(30);comment:昵称" json:"nickname,omitempty"`
	Sex         uint8  `gorm:"column:sex;type:uint;default:0;comment:性别 0保密 1男 2女" json:"sex,omitempty"`
	Mobile      string `gorm:"column:mobile;type:string;comment:手机号" json:"mobile,omitempty"`
	Email       string `gorm:"column:email;type:string;comment:邮箱地址" json:"email,omitempty"`
	QQ          string `gorm:"column:qq;type:varchar(20);comment:QQ" json:"qq,omitempty"`
	Avatar      string `gorm:"column:avatar;type:string;comment:用户头像地址" json:"avatar,omitempty"`
	Designation string `gorm:"column:designation;type:string;comment:称号" json:"designation,omitempty"`
	Realname    string `gorm:"column:realname;type:varchar(50);comment:真实姓名" json:"realname,omitempty"`
	Idcard      string `gorm:"column:idcard;type:varchar(100);comment:身份证号" json:"idcard,omitempty"`
	Address     string `gorm:"column:address;type:string;comment:地址" json:"address,omitempty"`
	Intro       string `gorm:"column:intro;type:string;comment:简介" json:"intro,omitempty"`
}
