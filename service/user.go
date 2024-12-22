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
func (us *UserService) ForgotPassword(email, newPassword string) error {
	var user model.User
	if errors.Is(db.Where("email = ?", email).First(&user).Error, gorm.ErrRecordNotFound) {
		return errors.New("邮箱错误！")
	}
	password, _ := util.DataEncryption(newPassword)
	result := db.Model(&user).Updates(model.User{Password: password})
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
	CollectNum  uint   `gorm:"column:collect_num;type:uint;default:0;comment:"`
	BrowseNum   uint   `gorm:"column:browse_num;type:uint;default:0;comment:"`
}

func (us *UserService) Info(id uint) (*User, error) {
	var user User

	ucl := db.Table("video_UserCollectLog as ucl").Select("ucl.user_id,count(ucl.video_id) as collect_num").Where("ucl.DeletedAt is null and ucl.user_id = ?", id).Group("ucl.user_id")
	ubl := db.Table("video_UserBrowseLog as ubl").Select("ubl.user_id,sum(ubl.number) as browse_num").Where("ubl.DeletedAt is null and ubl.user_id = ?", id).Group("ubl.user_id")
	if err := db.Table("video_User as users").
		Select("users.id,users.username,users.nickname,users.sex,users.mobile,users.email,users.qq,users.avatar,users.designation,users.realname,users.idcard,users.address,users.note,collect_num,browse_num").
		Joins("left join (?) as ucl on ucl.user_id = users.id", ucl).
		Joins("left join (?) as ubl on ubl.user_id = users.id", ubl).
		Where("users.id = ?", id).
		Group("users.id,users.username,users.nickname,users.sex,users.mobile,users.email,users.qq,users.avatar,users.designation,users.realname,users.idcard,users.address,users.note").
		First(&user).Error; err != nil {
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
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err := tx.Model(&user).Updates(updateUser).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (us *UserService) CollectLog(userID uint, videoID uint) (*model.UserCollectLog, error) {
	var data model.UserCollectLog
	result := db.Where("user_id = ? and video_id = ?", userID, videoID).First(&data)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("收藏记录不存在！")
	}
	return &data, nil
}

func (us *UserService) FavoriteList(userID uint) ([]Video, error) {
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

		videos = append(videos, Video{
			ID:       videoInfo.ID,
			Title:    videoInfo.Title,
			Poster:   videoInfo.Poster,
			Duration: videoInfo.Duration,
			Browse:   videoInfo.Browse,
			Collect:  videoInfo.Collect,
		})
	}

	return videos, nil
}

type VideoBrowse struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
	Poster   string `json:"poster"`
}

func (us *UserService) BrowseList(userID uint) ([]VideoBrowse, error) {
	var videos []model.Video
	result := db.Table("video_UserBrowseLog as a").
		Select("v.id,v.title,v.duration,v.poster").
		Joins("left join video_Video as v on v.id = a.video_id").
		Where("a.user_id = ? and a.UpdatedAt >= ?", userID, utils.NowTime().StringToTime("yesterday")).
		Order("a.UpdatedAt desc").
		Find(&videos)
	if err := result.Error; err != nil {
		return nil, err
	}
	data := make([]VideoBrowse, len(videos))
	for k, video := range videos {
		data[k] = VideoBrowse{
			ID:       video.ID,
			Title:    video.Title,
			Duration: utils.ResolveTime(uint32(video.Duration)),
			Poster:   video.Poster,
		}
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
