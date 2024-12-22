package model

import "gorm.io/gorm"

// VideoDanmu 视频弹幕表
type VideoDanmu struct {
	gorm.Model
	VideoId uint    `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	UserId  uint    `gorm:"column:user_id;type:uint;not null;default:0;comment:弹幕人的ID"`
	Text    string  `gorm:"column:text;type:text;not null;comment:弹幕文本"`
	Time    float64 `gorm:"column:time;type:double;not null;comment:弹幕时间, 默认为当前播放器时间"`
	Mode    uint8   `gorm:"column:mode;type:uint;not null;default:0;comment:弹幕模式: 0: 滚动(默认)，1: 顶部，2: 底部"`
	Color   string  `gorm:"column:color;type:text;not null;comment:弹幕颜色，默认为白色"`
	Border  bool    `gorm:"column:border;type:bool;not null;default:false;comment:弹幕是否有描边, 默认为 false"`
	Style   string  `gorm:"column:style;type:text;not null;comment:弹幕自定义样式, 默认为空对象"`
}
