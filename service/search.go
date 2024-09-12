package service

import (
	"bytes"
	"database/sql"
	"encoding/json"

	"github.com/wxw9868/video/utils"
)

type Index struct {
	Id       uint32      `json:"id" binding:"required"`
	Text     string      `json:"text" binding:"required"`
	Document interface{} `json:"document" binding:"required"`
}

func VideoWriteGoFound() error {
	rows, err := db.Table("video_Video as v").Select("v.*,l.collect, l.browse, l.zan, l.cai, l.watch").Joins("left join video_VideoLog l on l.video_id = v.id").Rows()
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var indexBatch []Index
	for rows.Next() {
		var videoInfo VideoInfo
		if err = db.ScanRows(rows, &videoInfo); err != nil {
			return err
		}

		//f, _ := strconv.ParseFloat(strconv.FormatInt(videoInfo.Size, 10), 64)
		video := Video{
			ID:       videoInfo.ID,
			Title:    videoInfo.Title,
			Poster:   videoInfo.Poster,
			Duration: utils.ResolveTime(uint32(videoInfo.Duration)),
			Browse:   videoInfo.Browse,
			Watch:    videoInfo.Watch,
			//Actress:       videoInfo.Actress,
			//Size:          f / 1024 / 1024,
			//ModTime:       videoInfo.CreationTime.Format("2006-01-02 15:04:05"),
			//Width:         videoInfo.Width,
			//Height:        videoInfo.Height,
			//CodecName:     videoInfo.CodecName,
			//ChannelLayout: videoInfo.ChannelLayout,
			//Collect:       videoInfo.Collect,
			//Zan:           videoInfo.Zan,
			//Cai:           videoInfo.Cai,
		}

		indexBatch = append(indexBatch, Index{
			Id:       uint32(videoInfo.ID),
			Text:     videoInfo.Title,
			Document: video,
		})
	}

	b, err := json.Marshal(&indexBatch)
	if err != nil {
		return err
	}
	if err = Post(utils.Join("/index/batch", "?", "database=video"), bytes.NewReader(b)); err != nil {
		return err
	}

	return nil
}
