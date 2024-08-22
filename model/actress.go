package model

import "gorm.io/gorm"

type Actress struct {
	gorm.Model
	Actress      string `gorm:"column:actress;type:varchar(100);uniqueIndex;comment:演员名称"`
	Alias        string `gorm:"column:alias;type:varchar(255);comment:演员别称"`
	Avatar       string `gorm:"column:avatar;type:varchar(255);comment:头像"`
	Birth        string `gorm:"column:birth;type:varchar(30);comment:出生"`
	Measurements string `gorm:"column:measurements;type:varchar(30);comment:三围"`
	CupSize      string `gorm:"column:cup_size;type:varchar(10);comment:罩杯"`
	DebutDate    string `gorm:"column:debut_date;type:varchar(30);comment:出道日期"`
	StarSign     string `gorm:"column:star_sign;type:varchar(20);comment:星座"`
	BloodGroup   string `gorm:"column:blood_group;type:varchar(5);comment:血型"`
	Stature      string `gorm:"column:stature;type:varchar(5);comment:身高"`
	Nationality  string `gorm:"column:nationality;type:varchar(255);comment:国籍"`
	Introduction string `gorm:"column:introduction;type:text;comment:简介"`
}
