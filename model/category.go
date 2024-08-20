package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string `gorm:"column:name;type:varchar(255);uniqueIndex;comment:分类名称"`
	Icon        string `gorm:"column:actress;type:varchar(100);comment:分类图标"`
	Sort        int    `gorm:"column:sort;type:int;default:0;comment:排序"`
	Description string `gorm:"column:description;type:size;comment:描述"`
}
