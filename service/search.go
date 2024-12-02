package service

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"

	gofoundClient "github.com/wxw9868/video/initialize/gofound"
	"github.com/wxw9868/video/utils"
)

type VideoData struct {
	ID            uint    `json:"id"`
	Title         string  `json:"title"`          // 视频标题
	Poster        string  `json:"poster"`         // 视频封面
	Duration      string  `json:"duration"`       // 视频时长
	Size          float64 `json:"size"`           // 视频大小
	CreationTime  string  `json:"creation_time"`  // 视频创建时间
	Width         int     `json:"width"`          // 视频宽度
	Height        int     `json:"height"`         // 视频长度
	CodecName     string  `json:"codec_name"`     // 视频编解码器
	ChannelLayout string  `json:"channel_layout"` // 视频音频声道
	CollectNum    uint    `json:"collect_num"`    // 收藏次数
	BrowseNum     uint    `json:"browse_num"`     // 浏览次数
	WatchNum      uint    `json:"watch_num"`      // 观看次数
	ZanNum        uint    `json:"zan_num"`        // 视频赞次数
	CaiNum        uint    `json:"cai_num"`        // 视频踩次数
}

type Index struct {
	Id       uint        `json:"id" binding:"required"`
	Text     string      `json:"text" binding:"required"`
	Document interface{} `json:"document" binding:"required"`
}

func VideoWriteGoFound() error {
	rows, err := db.Table("video_Video as v").
		Select("v.*,l.collect, l.browse, l.zan, l.cai, l.watch").
		Joins("left join video_VideoLog l on l.video_id = v.id").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	var indexBatch []Index
	for rows.Next() {
		var videoInfo VideoInfo
		if err = db.ScanRows(rows, &videoInfo); err != nil {
			return err
		}

		f, _ := strconv.ParseFloat(strconv.FormatInt(videoInfo.Size, 10), 64)
		indexBatch = append(indexBatch, Index{
			Id:   videoInfo.ID,
			Text: videoInfo.Title,
			Document: VideoData{
				ID:            videoInfo.ID,
				Title:         videoInfo.Title,
				Poster:        videoInfo.Poster,
				Duration:      utils.ResolveTime(uint32(videoInfo.Duration)),
				Size:          f / 1024 / 1024,
				CreationTime:  videoInfo.CreationTime.Format("2006-01-02 15:04:05"),
				Width:         videoInfo.Width,
				Height:        videoInfo.Height,
				CodecName:     videoInfo.CodecName,
				ChannelLayout: videoInfo.ChannelLayout,
				CollectNum:    videoInfo.Collect,
				BrowseNum:     videoInfo.Browse,
				WatchNum:      videoInfo.Watch,
				ZanNum:        videoInfo.Zan,
				CaiNum:        videoInfo.Cai,
			},
		})
	}

	b, err := json.Marshal(&indexBatch)
	if err != nil {
		return err
	}
	if err = Post(utils.Join("/index/batch", "?", "database=", "private-video"), bytes.NewReader(b)); err != nil {
		return err
	}

	return nil
}

func Post(url string, body io.Reader) error {
	resp, err := gofoundClient.GofoundClient().POST(url, "application/json", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
