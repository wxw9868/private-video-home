package service

import (
	"errors"

	"github.com/wxw9868/util"
	"github.com/wxw9868/video/model"
	"gorm.io/gorm"
)

type UserService struct{}

type APIUser struct {
	ID       uint
	Email    string `gorm:"column:email;type:string;comment:邮箱"`
	Password string `gorm:"column:password;type:string;comment:账户密码"`
}

func (us *UserService) Login(email, password string) (*APIUser, error) {
	password, err := util.DataEncryption(password)
	if err != nil {
		return nil, err
	}
	var user APIUser
	result := db.Model(&model.User{}).Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("用户不存在！")
	}
	if password != user.Password {
		return nil, errors.New("用户密码错误！")
	}
	return &user, nil
}

func (us *UserService) CollectLog(userID uint, videoID uint) (*model.UserCollectLog, error) {
	var data model.UserCollectLog
	result := db.Where("user_id = ? and video_id = ?", userID, videoID).First(&data)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("收藏记录不存在！")
	}
	return &data, nil
}
