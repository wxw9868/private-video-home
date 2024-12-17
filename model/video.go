package model

import (
	"time"

	"gorm.io/gorm"
)

// Video 视频基本信息
type Video struct {
	gorm.Model
	CateID        uint      `gorm:"column:cate_id;type:uint;not null;default:0;comment:视频分类ID"`
	Title         string    `gorm:"column:title;type:varchar(255);uniqueIndex;comment:标题"`
	Actress       string    `gorm:"column:actress;type:varchar(100);comment:演员"`
	Size          int64     `gorm:"column:size;type:bigint;comment:大小"`
	Duration      float64   `gorm:"column:duration;type:float;default:0;comment:时长"`
	Poster        string    `gorm:"column:poster;type:varchar(255);comment:封面"`
	Width         int       `gorm:"column:width;type:int;default:0;comment:宽"`
	Height        int       `gorm:"column:height;type:int;default:0;comment:高"`
	CodecName     string    `gorm:"column:codec_name;type:varchar(90);comment:编解码器"`
	ChannelLayout string    `gorm:"column:channel_layout;type:varchar(90);comment:音频声道"`
	CreationTime  time.Time `gorm:"column:creation_time;type:date;comment:时间"`
	Introduction  string    `gorm:"column:introduction;type:text;comment:简介"`
}

// VideoEpisodes 视频剧集信息
type VideoEpisodes struct {
	gorm.Model
	VideoId       uint      `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	Episode       uint      `gorm:"column:episode;type:uint;not null;default:0;comment:剧集"`
	Size          int64     `gorm:"column:size;type:bigint;comment:大小"`
	Duration      float64   `gorm:"column:duration;type:float;default:0;comment:时长"`
	Width         int       `gorm:"column:width;type:int;default:0;comment:宽"`
	Height        int       `gorm:"column:height;type:int;default:0;comment:高"`
	CodecName     string    `gorm:"column:codec_name;type:varchar(90);comment:编解码器"`
	ChannelLayout string    `gorm:"column:channel_layout;type:varchar(90);comment:音频声道"`
	CreationTime  time.Time `gorm:"column:creation_time;type:date;comment:时间"`
}

// 视频日志信息
type VideoLog struct {
	gorm.Model
	VideoID uint `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	Collect uint `gorm:"column:collect;type:uint;not null;default:0;comment:收藏"`
	Browse  uint `gorm:"column:browse;type:uint;not null;default:0;comment:浏览"`
	Zan     uint `gorm:"column:zan;type:uint;not null;default:0;comment:赞"`
	Cai     uint `gorm:"column:cai;type:uint;not null;default:0;comment:踩"`
	Watch   uint `gorm:"column:watch;type:uint;not null;default:0;comment:观看"`
}

type VideoComment struct {
	gorm.Model
	ParentId    uint   `gorm:"column:parent_id;type:uint;not null;default:0;comment:父级评论的ID"`
	VideoId     uint   `gorm:"column:video_id;type:uint;not null;default:0;comment:被评论的视频ID"`
	UserId      uint   `gorm:"column:user_id;type:uint;not null;default:0;comment:评论人的ID"`
	Nickname    string `gorm:"column:nickname;type:varchar(13);null;comment:评论人的昵称"`
	Avatar      string `gorm:"column:avatar;type:varchar(255);null;comment:评论人的头像地址"`
	Status      string `gorm:"column:status;type:check(status in ('VERIFYING','APPROVED','REJECT','DELETED'));default:'VERIFYING';comment:评论的状态"`
	ZanNum      uint   `gorm:"column:zan_num;type:uint;not null;default:0;comment:点赞人数"`
	ZanUserid   string `gorm:"column:zan_userid;type:varchar(255);not null;comment:"`
	ReplyNum    uint   `gorm:"column:reply_num;type:uint;not null;default:0;comment:评论回复数"`
	IsAnonymous uint   `gorm:"column:is_anonymous;type:uint;not null;default:0;comment:是否匿名评价 0是 1不是"`
	Content     string `gorm:"column:content;type:text;not null;comment:评论内容"`
	Remark      string `gorm:"column:remark;type:varchar(100);not null;comment:备注（审核不通过时添加）"`
	Support     uint   `gorm:"column:support;type:uint;not null;default:0;comment:支持（赞）"`
	Oppose      uint   `gorm:"column:oppose;type:uint;not null;default:0;comment:反对（踩）"`
	IsShow      uint8  `gorm:"column:is_show;type:uint;not null;default:1;comment:是否显示 0不显示 1显示"`
	// Status      string `gorm:"column:status;type:enum('VERIFYING','APPROVED','REJECT','DELETED');default:'VERIFYING';comment:评论的状态"`
}

type VideoActress struct {
	gorm.Model
	VideoId   uint `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	ActressId uint `gorm:"column:actress_id;type:uint;not null;default:0;comment:演员ID"`
}

type VideoTag struct {
	gorm.Model
	VideoId uint `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	TagId   uint `gorm:"column:tag_id;type:uint;not null;default:0;comment:标签ID"`
}

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
