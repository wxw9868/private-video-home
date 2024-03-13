package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Title    string  `gorm:"column:title;type:varchar(255);comment:标题"`
	Actress  string  `gorm:"column:actress;type:varchar(100);comment:演员"`
	Size     float64 `gorm:"column:size;type:float;comment:大小"`
	Duration int     `gorm:"column:duration;type:int;default:0;comment:时长"`
	ModTime  string  `gorm:"column:mod_time;type:varchar(255);comment:修改时间"`
	Poster   string  `gorm:"column:poster;type:varchar(255);comment:封面"`
}
