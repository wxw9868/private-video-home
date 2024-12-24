package service

import (
	"errors"
	"fmt"

	"github.com/wxw9868/util"
	"github.com/wxw9868/util/randomname"
	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
	"gorm.io/gorm"
)

type UserService struct{}

func (us *UserService) Register(username, email, password string) error {
	if !errors.Is(db.Where("email = ?", email).First(&model.User{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("邮箱已存在！")
	}
	if !errors.Is(db.Where("username = ?", username).First(&model.User{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("用户已存在！")
	}
	password, err := util.DataEncryption(password)
	if err != nil {
		return err
	}

	if err := db.Create(&model.User{Username: username, Nickname: randomname.GenerateName(1), Password: password, Email: email, Avatar: "assets/image/avatar/avatar.png"}).Error; err != nil {
		return fmt.Errorf("注册失败: %s", err)
	}
	return nil
}

type APIUser struct {
	ID       uint
	Email    string `gorm:"column:email;type:string;comment:邮箱"`
	Password string `gorm:"column:password;type:string;comment:账户密码"`
}

func (us *UserService) Login(email, password string) (*model.User, error) {
	password, err := util.DataEncryption(password)
	if err != nil {
		return nil, err
	}
	var user model.User
	result := db.Model(&model.User{}).Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("用户不存在！")
	}
	if password != user.Password {
		return nil, errors.New("用户密码错误！")
	}
	return &user, nil
}

func (us *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	var user model.User
	if err := db.First(&user, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("用户不存在！")
	}

	oldPassword, _ = util.DataEncryption(oldPassword)
	if oldPassword != user.Password {
		return errors.New("原密码输入错误！")
	}

	newPassword, _ = util.DataEncryption(newPassword)
	if err := db.Model(&user).Updates(model.User{Password: newPassword}).Error; err != nil {
		return errors.New("修改密码失败！")
	}
	return nil
}
func (us *UserService) ForgotPassword(email, newPassword string) error {
	var user model.User
	if errors.Is(db.Where("email = ?", email).First(&user).Error, gorm.ErrRecordNotFound) {
		return errors.New("邮箱错误！")
	}

	newPassword, _ = util.DataEncryption(newPassword)
	if err := db.Model(&user).Updates(model.User{Password: newPassword}).Error; err != nil {
		return errors.New("修改密码失败！")
	}
	return nil
}

type User struct {
	ID          uint   `json:"id,omitempty"`
	Username    string `gorm:"column:username;type:varchar(100);uniqueIndex;comment:账户名" json:"username,omitempty"`
	Nickname    string `gorm:"column:nickname;type:varchar(30);comment:昵称" json:"nickname,omitempty"`
	Sex         uint8  `gorm:"column:sex;type:uint;default:0;comment:性别 0保密 1男 2女" json:"sex,omitempty"`
	Mobile      string `gorm:"column:mobile;type:string;comment:手机号" json:"mobile,omitempty"`
	Email       string `gorm:"column:email;type:string;comment:邮箱地址" json:"email,omitempty"`
	QQ          string `gorm:"column:qq;type:varchar(20);comment:QQ" json:"qq,omitempty"`
	Avatar      string `gorm:"column:avatar;type:string;comment:用户头像地址" json:"avatar,omitempty"`
	Designation string `gorm:"column:designation;type:string;comment:称号" json:"designation,omitempty"`
	Realname    string `gorm:"column:realname;type:varchar(50);comment:真实姓名" json:"realname,omitempty"`
	Idcard      string `gorm:"column:idcard;type:varchar(100);comment:身份证号" json:"idcard,omitempty"`
	Address     string `gorm:"column:address;type:string;comment:地址" json:"address,omitempty"`
	Note        string `gorm:"column:note;type:string;comment:备注" json:"note,omitempty"`
	CollectNum  uint   `gorm:"column:collect_num;type:uint;default:0;comment:" json:"collect_num,omitempty"`
	BrowseNum   uint   `gorm:"column:browse_num;type:uint;default:0;comment:" json:"browse_num,omitempty"`
}

func (us *UserService) Info(id uint) (*User, error) {
	var user User
	a := db.Table("video_UserCollectLog as a").Select("a.user_id,count(a.video_id) as collect_num").Where("a.DeletedAt is null and a.user_id = ?", id).Group("a.user_id")
	b := db.Table("video_UserPageViewsLog as b").Select("b.user_id,sum(b.page_views) as browse_num").Where("b.DeletedAt is null and b.user_id = ?", id).Group("b.user_id")
	err := db.Table("video_User as users").
		Select("users.id,users.username,users.nickname,users.sex,users.mobile,users.email,users.qq,users.avatar,users.designation,users.realname,users.idcard,users.address,users.note,collect_num,browse_num").
		Joins("left join (?) as a on a.user_id = users.id", a).
		Joins("left join (?) as b on b.user_id = users.id", b).
		Where("users.id = ?", id).
		Group("users.id,users.username,users.nickname,users.sex,users.mobile,users.email,users.qq,users.avatar,users.designation,users.realname,users.idcard,users.address,users.note").
		First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UserService) Update(id uint, column string, value interface{}) error {
	if err := db.Model(&model.User{}).Where("id = ?", id).Update(column, value).Error; err != nil {
		return err
	}
	return nil
}

func (us *UserService) Updates(user model.User) error {
	if err := db.Model(&user).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

func (us *UserService) CollectLog(userID uint, videoID uint) (*model.UserCollectLog, error) {
	var data model.UserCollectLog
	result := db.Where("user_id = ? and video_id = ?", userID, videoID).First(&data)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("收藏记录不存在！")
	}
	return &data, nil
}

func (us *UserService) CollectList(userID uint) ([]Video, error) {
	var videos []Video
	err := db.Table("video_UserCollectLog as a").
		Select("v.id, v.title, v.poster, v.duration, l.collection_volume, l.page_views, l.likes_count, l.dislikes_count").
		Joins("left join video_Video as v on v.id = a.video_id").
		Joins("left join video_VideoLog l on l.video_id = a.video_id").
		Where("a.user_id = ? and a.DeletedAt is null", userID).
		Order("a.id desc").Scan(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

type VideoPageViews struct {
	ID       uint    `json:"id"`
	Title    string  `gorm:"column:title;type:varchar(255);uniqueIndex;comment:标题" json:"title"`
	Duration float64 `gorm:"column:duration;type:float;default:0;comment:时长" json:"duration"`
	Poster   string  `gorm:"column:poster;type:varchar(255);comment:封面" json:"poster"`
}

func (us *UserService) VideoPageViewsList(userID uint) ([]VideoPageViews, error) {
	var data []VideoPageViews
	err := db.Table("video_UserPageViewsLog as a").
		Select("v.id,v.title,v.duration,v.poster").
		Joins("left join video_Video as v on v.id = a.video_id").
		Where("a.user_id = ? and a.UpdatedAt >= ?", userID, utils.NowTime().StringToTime("yesterday")).
		Order("a.UpdatedAt desc").Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (us *UserService) CreateUserLoginLog(m model.UserLoginLog) error {
	if err := db.Create(&m).Error; err != nil {
		return err
	}
	return nil
}

func (us *UserService) UserLoginLogList(id uint) ([]model.UserLoginLog, error) {
	var data []model.UserLoginLog
	if err := db.Model(&model.UserLoginLog{}).Where("user_id = ?", id).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
