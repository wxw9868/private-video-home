package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func videoList(c *gin.Context) {
	if c.Query("actress_id") != "" {
		indexInt, _ := strconv.Atoi(c.Query("actress_id"))
		actress := actressListSort[indexInt]
		indexs := actressList[actress]
		actressVideos := make([]video, len(indexs))
		for k, index := range indexs {
			actressVideos[k] = videos[index]
		}

		videosBytes, _ := json.Marshal(actressVideos)

		c.HTML(http.StatusOK, "list.html", gin.H{
			"title":       "视频列表",
			"data":        string(videosBytes),
			"actressList": actressListSort,
			"actressID":   indexInt,
		})
		return
	}

	if videos == nil && list == nil {
		files, err := os.ReadDir(videoDir)
		if err != nil {
			log.Fatal(err)
		}

		list = make([]string, len(files))
		videos = make([]video, len(files))

		for k, file := range files {
			filename := file.Name()
			ext := filepath.Ext(filename)
			if ext == ".mp4" {
				strs := strings.Split(filename, ".")
				title := strs[0]
				arrs := strings.Split(strs[0], "_")
				actress := arrs[len(arrs)-1]
				if _, ok := actressList[actress]; !ok {
					actressListSort = append(actressListSort, actress)
				}
				actressList[actress] = append(actressList[actress], k)

				fi, _ := file.Info()
				size, _ := strconv.ParseFloat(strconv.FormatInt(fi.Size(), 10), 64)

				filePath := videoDir + "/" + filename
				fil, _ := os.Open(filePath)
				duration, _ := GetMP4Duration(fil)

				posterPath := posterDir + "/" + title + ".gif"
				//_ = ReadFrameAsJpeg(filePath, posterPath, "00:00:55")

				//snapshotPath := snapshotDir + "/" + title + ".gif"
				//_ = CutVideoForGif(filePath, posterPath)

				list[k] = filename
				videos[k] = video{
					ID:       k + 1,
					Title:    title,
					Actress:  actress,
					Size:     size / 1024 / 1024,
					Duration: ResolveTime(duration),
					ModTime:  fi.ModTime().Format("2006-01-02 15:04:05"),
					Poster:   posterPath,
				}
			}
		}
	}

	videosBytes, _ := json.Marshal(videos)

	c.HTML(http.StatusOK, "list.html", gin.H{
		"title":       "视频列表",
		"data":        string(videosBytes),
		"actressList": actressListSort,
		"actressID":   -1,
	})
}
func videoPlay(c *gin.Context) {
	id := c.Query("id")
	intId, _ := strconv.Atoi(id)
	name := list[intId]

	c.HTML(http.StatusOK, "play.html", gin.H{
		"title":     "视频播放",
		"video_url": name,
	})
}

var replaceName = map[string]struct{}{
	"_tg关注_@AVWUMAYUANPIAN": {},
	"无码频道_每天更新_":            {},
}

func videoRename(c *gin.Context) {
	files, err := os.ReadDir(videoDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		filename := file.Name()
		oldpath := videoDir + "/" + filename
		filename = strings.Replace(filename, "无码频道_每天更新_", "", -1)
		newpath := videoDir + "/" + filename
		os.Rename(oldpath, newpath)
	}
}
