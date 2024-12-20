package model

import "gorm.io/gorm"

// Tag 标签表
type Tag struct {
	gorm.Model
	Name   string `gorm:"column:name;type:varchar(255);uniqueIndex;comment:演员"`
	Avatar string `gorm:"column:avatar;type:varchar(255);comment:头像"`
}
