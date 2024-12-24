package model

import "gorm.io/gorm"

// Tag 标签表
type Tag struct {
	gorm.Model
	Name string `gorm:"column:name;type:varchar(255);uniqueIndex;comment:标签名称"`
	Icon string `gorm:"column:icon;type:varchar(255);comment:标签图标"`
}
