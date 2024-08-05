package model

import "gorm.io/gorm"

type Actress struct {
	gorm.Model
	Actress string `gorm:"column:actress;type:varchar(100);uniqueIndex;comment:演员"`
	Avatar  string `gorm:"column:avatar;type:varchar(255);comment:头像"`
}
