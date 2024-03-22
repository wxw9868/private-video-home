package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wxw9868/video/service"
)

func VideoIndex(c *gin.Context) {
	cards := make([]string, 6)
	for i := 0; i < 6; i++ {
		cardPath := "./assets/image/bizhi/card" + strconv.Itoa(i+1) + ".jpeg"
		cards[i] = cardPath
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "首页",
		"data":  cards,
	})
}

func VideoList(c *gin.Context) {
	actressID := c.Query("actress_id")
	videos, err := vs.Find(actressID)
	if err != nil {
		log.Fatal(err)
	}
	videosBytes, _ := json.Marshal(videos)

	c.HTML(http.StatusOK, "list.html", gin.H{
		"title":       "视频列表",
		"data":        string(videosBytes),
		"actressList": actressListSort,
		"actressID":   -1,
	})
}

func VideoPlay(c *gin.Context) {
	id := c.Query("id")
	player := c.Query("player")

	var name string
	if player == "ckplayer" {
		name = "ckplayer.html"
	} else {
		name = "xgplayer.html"
	}

	vs, err := vs.First(id)
	if err != nil {
		log.Fatal(err)
	}

	c.HTML(http.StatusOK, name, gin.H{
		"title":     "视频播放",
		"videoUrl":  videoDir + "/" + vs.Title + ".mp4",
		"videoName": vs.Actress,
	})
}

func VideoActress(c *gin.Context) {
	actresss, err := as.Find()
	if err != nil {
		log.Fatal(err)
	}

	actressBytes, _ := json.Marshal(actresss)

	c.HTML(http.StatusOK, "actress.html", gin.H{
		"title":       "演员列表",
		"actressList": string(actressBytes),
	})
}

func VideoImport(c *gin.Context) {
	service.VideoImport()
	c.JSON(http.StatusOK, "SUCCESS")
}

func VideoRename(c *gin.Context) {
	var videoDir = "C:/Users/wxw9868/Downloads/Telegram Desktop"
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
