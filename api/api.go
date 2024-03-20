package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wxw9868/video/service"
	"github.com/wxw9868/video/utils"
)

func LoginApi(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "登录",
	})
}

func DoLoginApi(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	fmt.Println("s: ", email, password)

	if email != "" && password != "" {
		session := sessions.Default(c)
		session.Set("email", email)
		session.Set("password", password)
		session.Save()

		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
}

func LogoutApi(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, nil)
}

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
		writeToCache()
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
	intId, _ := strconv.Atoi(id)
	name := list[intId]

	videoUrl := videoDir + "/" + name

	c.HTML(http.StatusOK, "play.html", gin.H{
		"title":     "视频播放",
		"videoUrl":  videoUrl,
		"videoName": name,
	})
}

func VideoActress(c *gin.Context) {
	// actresss, _ := as.Find()
	// fmt.Printf("%+v\n", actresss)
	// return

	actressList := actressListSort
	actressSlice := make([]actress, len(actressList))

	if len(actressList) > 0 {
		for index, name := range actressListSort {
			nameSlice := []rune(name)
			avatarPath := avatarDir + "/" + name + ".png"

			_, err := os.Stat(avatarPath)
			if os.IsNotExist(err) {
				utils.GenerateAvatar(string(nameSlice[0]), avatarPath)
			}

			actressSlice[index] = actress{
				ID:      index + 1,
				Actress: name,
				Avatar:  avatarPath,
			}
		}
	}

	actressBytes, _ := json.Marshal(actressSlice)

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

func writeToCache() error {
	files, err := os.ReadDir(videoDir)
	if err != nil {
		return err
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
			duration, _ := utils.GetMP4Duration(fil)

			videoInfo, _ := utils.VideoInfo(filePath)

			posterPath := posterDir + "/" + title + ".jpg"
			_, err = os.Stat(posterPath)
			if os.IsNotExist(err) {
				_ = utils.ReadFrameAsJpeg(filePath, posterPath, "00:02:30")
			}

			//snapshotPath := snapshotDir + "/" + title + ".gif"
			//_ = CutVideoForGif(filePath, posterPath)

			list[k] = filename
			videos[k] = video{
				ID:            k + 1,
				Title:         title,
				Actress:       actress,
				Size:          size / 1024 / 1024,
				Duration:      utils.ResolveTime(duration),
				ModTime:       fi.ModTime().Format("2006-01-02 15:04:05"),
				Poster:        posterPath,
				Width:         int(videoInfo["width"].(int64)),
				Height:        int(videoInfo["height"].(int64)),
				CodecName:     fmt.Sprintf("%s,%s", videoInfo["codec_name0"].(string), videoInfo["codec_name1"].(string)),
				ChannelLayout: videoInfo["channel_layout"].(string),
			}
		}
	}
	return nil
}
