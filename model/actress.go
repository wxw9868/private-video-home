package model

import "gorm.io/gorm"

// Actress 演员表
type Actress struct {
	gorm.Model
	Actress      string `gorm:"column:actress;type:varchar(100);uniqueIndex;comment:演员名称" json:"actress,omitempty"`
	Alias        string `gorm:"column:alias;type:varchar(255);comment:演员别称" json:"alias,omitempty"`
	Avatar       string `gorm:"column:avatar;type:varchar(255);comment:头像" json:"avatar,omitempty"`
	Birth        string `gorm:"column:birth;type:varchar(30);comment:出生" json:"birth,omitempty"`
	Measurements string `gorm:"column:measurements;type:varchar(30);comment:三围" json:"measurements,omitempty"`
	CupSize      string `gorm:"column:cup_size;type:varchar(10);comment:罩杯" json:"cup_size,omitempty"`
	DebutDate    string `gorm:"column:debut_date;type:varchar(30);comment:出道日期" json:"debut_date,omitempty"`
	StarSign     string `gorm:"column:star_sign;type:varchar(20);comment:星座" json:"star_sign,omitempty"`
	BloodGroup   string `gorm:"column:blood_group;type:varchar(5);comment:血型" json:"blood_group,omitempty"`
	Stature      string `gorm:"column:stature;type:varchar(5);comment:身高" json:"stature,omitempty"`
	Nationality  string `gorm:"column:nationality;type:varchar(255);comment:国籍" json:"nationality,omitempty"`
	Intro        string `gorm:"column:intro;type:text;comment:简介" json:"intro,omitempty"`
}
