package service

import (
	"strconv"

	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
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
