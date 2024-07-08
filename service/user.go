package service

import (
	"errors"
	"fmt"
	"strconv"

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
	result := db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("用户不存在！")
	}

	password, _ := util.DataEncryption(oldPassword)
	if password != user.Password {
		return errors.New("原密码输入错误！")
	}

	password, _ = util.DataEncryption(newPassword)
	result = db.Model(&user).Updates(model.User{Password: password})
	if result.Error != nil {
		return errors.New("修改密码失败！")
	}
	return nil
}

type User struct {
	ID          uint
	Username    string `gorm:"column:username;type:varchar(100);uniqueIndex;comment:账户名"`
	Nickname    string `gorm:"column:nickname;type:varchar(30);comment:昵称"`
	Sex         uint8  `gorm:"column:sex;type:uint;default:0;comment:性别 0保密 1男 2女"`
	Mobile      string `gorm:"column:mobile;type:string;comment:手机号"`
	Email       string `gorm:"column:email;type:string;comment:邮箱地址"`
	QQ          string `gorm:"column:qq;type:varchar(20);comment:QQ"`
	Avatar      string `gorm:"column:avatar;type:string;comment:用户头像地址"`
	Designation string `gorm:"column:designation;type:string;comment:称号"`
	Realname    string `gorm:"column:realname;type:varchar(50);comment:真实姓名"`
	Idcard      string `gorm:"column:idcard;type:varchar(100);comment:身份证号"`
	Address     string `gorm:"column:address;type:string;comment:地址"`
	Note        string `gorm:"column:note;type:string;comment:备注"`
}

func (us *UserService) Info(id uint) (*User, error) {
	var user User
	if err := db.First(&user, id).Error; err != nil {
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

func (us *UserService) Updates(id uint, updateUser model.User) error {
	var user model.User

	tx := db.Begin()
	if err := tx.First(&user, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&user).Updates(updateUser).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
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
	rows, err := db.Table("video_UserCollectLog as ucl").
		Select("v.*,l.collect, l.browse, l.zan, l.cai, l.watch").
		Joins("left join video_Video as v on v.id = ucl.video_id").
		Joins("left join video_VideoLog l on l.video_id = ucl.video_id").
		Where("ucl.user_id = ? and ucl.DeletedAt is null", userID).
		Order("ucl.id desc").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var videoInfo VideoInfo
		db.ScanRows(rows, &videoInfo)

		f, _ := strconv.ParseFloat(strconv.FormatInt(videoInfo.Size, 10), 64)
		video := Video{
			ID:            videoInfo.ID,
			Title:         videoInfo.Title,
			Actress:       videoInfo.Actress,
			Size:          f / 1024 / 1024,
			Duration:      utils.ResolveTime(uint32(videoInfo.Duration)),
			ModTime:       videoInfo.CreationTime.Format("2006-01-02 15:04:05"),
			Poster:        videoInfo.Poster,
			Width:         videoInfo.Width,
			Height:        videoInfo.Height,
			CodecName:     videoInfo.CodecName,
			ChannelLayout: videoInfo.ChannelLayout,
			Collect:       videoInfo.Collect,
			Browse:        videoInfo.Browse,
			Zan:           videoInfo.Zan,
			Cai:           videoInfo.Cai,
			Watch:         videoInfo.Watch,
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (us *UserService) BrowseList(userID uint) ([]model.Video, error) {
	var data []model.Video
	result := db.Table("video_UserBrowseLog as ubl").
		Select("v.*").
		Joins("left join video_Video as v on v.id = ubl.video_id").
		Where("ubl.user_id = ?", userID).
		Order("ubl.id desc").
		Find(&data)
	if err := result.Error; err != nil {
		return nil, err
	}
	return data, nil
}
