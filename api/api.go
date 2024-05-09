package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
	"github.com/wxw9868/video/service"
)

func VideoIndex(c *gin.Context) {
	cards := make([]string, 14)
	for i := 0; i < 14; i++ {
		cardPath := "./assets/image/card/card" + strconv.Itoa(i+1) + ".jpeg"
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
		log.Panicln(err)
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
	} else if player == "xgplayer" {
		name = "xgplayer.html"
	} else {
		name = "player.html"
	}

	vs, err := vs.First(id)
	if err != nil {
		log.Printf("%s\n", err)
	}

	c.HTML(http.StatusOK, name, gin.H{
		"title":        "视频播放",
		"videoID":      vs.ID,
		"videoTitle":   vs.Title,
		"videoActress": vs.Actress,
		"videoUrl":     videoDir + "/" + vs.Title + ".mp4",
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

type VideoCollect struct {
	VideoID uint `form:"video_id" json:"video_id" binding:"required"`
	Collect int  `form:"collect" json:"collect" binding:"required,oneof=1 -1"`
}

// 收藏
func VideoCollectApi(c *gin.Context) {
	var bind VideoCollect
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := vs.Collect(bind.VideoID, bind.Collect); err != nil {
		c.JSON(http.StatusOK, util.Fail(err.Error()))
	}

	msg := "收藏成功"
	if bind.Collect == -1 {
		msg = "取消收藏"
	}
	c.JSON(http.StatusOK, util.Success(msg, nil))
}

func VideoImport(c *gin.Context) {
	service.VideoImport()
	c.JSON(http.StatusOK, gin.H{
		"message": "SUCCESS",
	})
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
		filename = strings.Replace(filename, "_tg关注_@AVWUMAYUANPIAN", "", -1)
		filename = strings.Replace(filename, "无码频道_每天更新_", "", -1)
		newpath := videoDir + "/" + filename
		os.Rename(oldpath, newpath)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "SUCCESS",
	})
}
