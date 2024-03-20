package model

import (
	"math/big"
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Title    string    `gorm:"column:title;type:varchar(255);comment:标题"`
	Actress  string    `gorm:"column:actress;type:varchar(100);comment:演员"`
	Size     big.Int   `gorm:"column:size;type:bigint;comment:大小"`
	Duration int       `gorm:"column:duration;type:int;default:0;comment:时长"`
	ModTime  time.Time `gorm:"column:mod_time;type:time;comment:修改时间"`
	Poster   string    `gorm:"column:poster;type:varchar(255);comment:封面"`
}

type Actress struct {
	gorm.Model
	Actress string `gorm:"column:actress;type:varchar(100);comment:演员"`
	Avatar  string `gorm:"column:poster;type:varchar(255);comment:头像"`
}
