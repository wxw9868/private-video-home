package service

import (
	"errors"
	"strconv"

	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
	"gorm.io/gorm"
)

type VideoService struct{}

type Video struct {
	ID            uint    `json:"id"`
	Title         string  `json:"title"`
	Actress       string  `json:"actress"`
	Size          float64 `json:"size"`
	Duration      string  `json:"duration"`
	ModTime       string  `json:"mod_time"`
	Poster        string  `json:"poster"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	CodecName     string  `json:"codec_name"`
	ChannelLayout string  `json:"channel_layout"`
}

func (as *VideoService) Find(actressID string) ([]Video, error) {
	dbVideo := db.Model(&model.Video{})
	if actressID != "" {
		var actress model.Actress
		if err := db.Select("Actress").Where("id = ?", actressID).First(&actress).Error; err != nil {
			return nil, err
		}
		dbVideo = dbVideo.Where("actress = ?", actress.Actress)
	}

	rows, err := dbVideo.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var modelVideo model.Video
		db.ScanRows(rows, &modelVideo)

		f, _ := strconv.ParseFloat(strconv.FormatInt(modelVideo.Size, 10), 64)
		videos = append(videos, Video{
			ID:            modelVideo.ID,
			Title:         modelVideo.Title,
			Actress:       modelVideo.Actress,
			Size:          f / 1024 / 1024,
			Duration:      utils.ResolveTime(uint32(modelVideo.Duration)),
			ModTime:       modelVideo.CreationTime.Format("2006-01-02 15:04:05"),
			Poster:        modelVideo.Poster,
			Width:         modelVideo.Width,
			Height:        modelVideo.Height,
			CodecName:     modelVideo.CodecName,
			ChannelLayout: modelVideo.ChannelLayout,
		})
	}
	return videos, nil
}

func (vs *VideoService) First(id string) (model.Video, error) {
	var video model.Video
	if err := db.Where("id = ?", id).First(&video).Error; err != nil {
		return video, err
	}
	return video, nil
}

type VideoInfo struct {
	gorm.Model
	VideoID uint `gorm:"column:video_id;type:uint;not null;default:0;comment:视频ID"`
	Collect uint `gorm:"column:collect;type:uint;not null;default:0;comment:收藏"`
	Browse  uint `gorm:"column:browse;type:uint;not null;default:0;comment:浏览"`
	Zan     uint `gorm:"column:zan;type:uint;not null;default:0;comment:赞"`
	Cai     uint `gorm:"column:cai;type:uint;not null;default:0;comment:踩"`
	Watch   uint `gorm:"column:watch;type:uint;not null;default:0;comment:观看"`
	Video   Video
}

func (vs *VideoService) Info(id string) (model.Video, error) {
	var video model.Video
	if err := db.Where("id = ?", id).First(&video).Error; err != nil {
		return video, err
	}
	return video, nil
}

func (vs *VideoService) List() ([]model.Video, error) {
	var videos []model.Video
	if err := db.Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

func (vs *VideoService) Create(videos []model.Video) error {
	if err := db.Create(&videos).Error; err != nil {
		return err
	}
	return nil
}

func (vs *VideoService) Collect(videoID uint, collect int) error {
	var video model.Video
	result := db.First(&video, videoID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("视频不存在！")
	}

	var expr string
	if collect == 1 {
		// 增加1
		expr = "collect + 1"
	} else {
		// 减少1
		expr = "collect - 1"
	}
	result = db.Model(&model.VideoLog{}).Where("video_id = ?", videoID).Update("collect", gorm.Expr(expr))
	if result.Error != nil {
		return errors.New("更新失败！")
	}
	return nil
}
