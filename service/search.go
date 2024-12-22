package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	gofoundClient "github.com/wxw9868/video/initialize/gofound"
	"github.com/wxw9868/video/utils"
)

type VideoData struct {
	ID            uint    `json:"id"`
	Title         string  `json:"title"`         // 视频标题
	Poster        string  `json:"poster"`        // 视频封面
	Duration      string  `json:"duration"`      // 视频时长
	Size          float64 `json:"size"`          // 视频大小
	CreationTime  string  `json:"creationTime"`  // 视频创建时间
	Width         int     `json:"width"`         // 视频宽度
	Height        int     `json:"height"`        // 视频长度
	CodecName     string  `json:"codecName"`     // 视频编解码器
	ChannelLayout string  `json:"channelLayout"` // 视频音频声道
	CollectNum    uint    `json:"collectNum"`    // 收藏次数
	BrowseNum     uint    `json:"browseNum"`     // 浏览次数
	WatchNum      uint    `json:"watchNum"`      // 观看次数
	LikeNum       uint    `json:"likeNum"`       // 视频赞次数
	DisLikeNum    uint    `json:"dislikeNum"`    // 视频踩次数
}

type Index struct {
	Id       uint        `json:"id" binding:"required"`
	Text     string      `json:"text" binding:"required"`
	Document interface{} `json:"document" binding:"required"`
}

func VideoWriteGoFound(query string) error {
	vdb := db.Table("video_Video as v").
		Select("v.*,l.collect, l.browse, l.like, l.dislike, l.watch").
		Joins("left join video_VideoLog l on l.video_id = v.id")
	if query != "" {
		vdb = vdb.Where(query)
	}
	rows, err := vdb.Rows()
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

		indexBatch = append(indexBatch, Index{
			Id:   videoInfo.ID,
			Text: videoInfo.Title,
			Document: VideoData{
				ID:            videoInfo.ID,
				Title:         videoInfo.Title,
				Poster:        videoInfo.Poster,
				Duration:      utils.ResolveTime(uint32(videoInfo.Duration)),
				Size:          float64(videoInfo.Size) / 1024 / 1024,
				CreationTime:  videoInfo.CreationTime.Format("2006-01-02 15:04:05"),
				Width:         videoInfo.Width,
				Height:        videoInfo.Height,
				CodecName:     videoInfo.CodecName,
				ChannelLayout: videoInfo.ChannelLayout,
				CollectNum:    videoInfo.Collect,
				BrowseNum:     videoInfo.Browse,
				WatchNum:      videoInfo.Watch,
				LikeNum:       videoInfo.Like,
				DisLikeNum:    videoInfo.Dislike,
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

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("resp: %s\n", string(b))

	return nil
}
