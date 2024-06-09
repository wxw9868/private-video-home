package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	sqlite "github.com/wxw9868/video/initialize/db"
	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
	"gorm.io/gorm"
)

var db = sqlite.DB()

var videoDir = "./assets/video"
var posterDir = "./assets/image/poster"
var avatarDir = "./assets/image/avatar"

var list []string
var videos []model.Video
var actresss []model.Actress
var actressList = make(map[string][]int)
var actressListSort []string
var mux = new(sync.Mutex)
var once = new(sync.Once)

func VideoImport() error {
	var videoDir = "D:/GoLang/ta"
	files, err := os.ReadDir(videoDir)
	if err != nil {
		log.Fatal(err)
	}

	list = make([]string, len(files))
	videos = make([]model.Video, len(files))

	// INSERT OR IGNORE INTO
	// INSERT OR REPLACE INTO
	var videoSql, actressSql string
	videoSql = "INSERT OR REPLACE INTO video_Video (title, actress, size, duration, poster, width, height, codec_name, channel_layout, creation_time, CreatedAt, UpdatedAt) VALUES "
	actressSql = "INSERT OR REPLACE INTO video_Actress (actress, avatar, CreatedAt, UpdatedAt) VALUES "

	for k, file := range files {
		filename := file.Name()
		ext := filepath.Ext(filename)
		if ext == ".mp4" {
			strs := strings.Split(filename, ".")
			title := strs[0]
			arrs := strings.Split(strs[0], "_")
			actress := arrs[len(arrs)-1]
			mux.Lock()
			if _, ok := actressList[actress]; !ok {
				actressListSort = append(actressListSort, actress)
			}
			actressList[actress] = append(actressList[actress], k)
			mux.Unlock()

			filePath := videoDir + "/" + filename
			posterPath := posterDir + "/" + title + ".jpg"
			_, err = os.Stat(posterPath)
			if os.IsNotExist(err) {
				_ = utils.ReadFrameAsJpeg(filePath, posterPath, "00:02:00")
			}
			videoInfo, _ := utils.VideoInfo(filePath)

			//snapshotPath := snapshotDir + "/" + title + ".gif"
			//_ = CutVideoForGif(filePath, posterPath)

			videoSql += fmt.Sprintf("('%s', '%s', %d, %f, '%s', %d, %d, '%s', '%s', '%v', '%v', '%v'), ",
				title,
				actress,
				videoInfo["size"].(int64),
				videoInfo["duration"].(float64),
				posterPath,
				videoInfo["width"].(int64),
				videoInfo["height"].(int64),
				fmt.Sprintf("%s,%s",
					videoInfo["codec_name0"].(string),
					videoInfo["codec_name1"].(string)),
				videoInfo["channel_layout"].(string),
				videoInfo["creation_time"].(time.Time),
				time.Now().Local(),
				time.Now().Local(),
			)

			// list[k] = filename
			// videos[k] = model.Video{
			// 	Title:         title,
			// 	Actress:       actress,
			// 	Size:          videoInfo["size"].(int64),
			// 	Duration:      videoInfo["duration"].(float64),
			// 	Poster:        posterPath,
			// 	Width:         int(videoInfo["width"].(int64)),
			// 	Height:        int(videoInfo["height"].(int64)),
			// 	CodecName:     fmt.Sprintf("%s,%s", videoInfo["codec_name0"].(string), videoInfo["codec_name1"].(string)),
			// 	ChannelLayout: videoInfo["channel_layout"].(string),
			// 	CreationTime:  videoInfo["creation_time"].(time.Time),
			// }
		}
	}

	actresss = make([]model.Actress, len(actressListSort))
	if len(actressListSort) > 0 {
		for _, name := range actressListSort {
			nameSlice := []rune(name)
			avatarPath := avatarDir + "/" + name + ".png"

			_, err := os.Stat(avatarPath)
			if os.IsNotExist(err) {
				if err := utils.GenerateAvatar(string(nameSlice[0]), avatarPath); err != nil {
					log.Fatal(err)
				}
			}

			actressSql += fmt.Sprintf("('%s', '%s', '%v', '%v'), ", name, avatarPath, time.Now().Local(), time.Now().Local())
			// actresss[index] = model.Actress{
			// 	Actress: name,
			// 	Avatar:  avatarPath,
			// }
		}
	}

	// fmt.Println(videoSql)
	// fmt.Println(actressSql)
	videoSqlBytes := []byte(videoSql)
	actressSqlBytes := []byte(actressSql)
	videoSql = string(videoSqlBytes[:len(videoSqlBytes)-2])
	actressSql = string(actressSqlBytes[:len(actressSqlBytes)-2])
	// fmt.Println(string(videoSqlBytes[:len(videoSqlBytes)-2]))
	// fmt.Println(string(actressSqlBytes[:len(actressSqlBytes)-2]))

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(videoSql).Error; err != nil {
			return err
		}
		if err := tx.Exec(actressSql).Error; err != nil {
			return err
		}

		// if err := tx.CreateInBatches(videos, 30).Error; err != nil {
		// 	return err
		// }
		// if err := tx.CreateInBatches(actresss, 30).Error; err != nil {
		// 	return err
		// }
		return nil
	})
	return err
}

func ImportActress() {
	var data []model.Actress
	var actressMap = make(map[string]struct{})

	ReadFileToMap("actress.json", &actressMap)
	fmt.Printf("map: %+v\n", actressMap)

	rows, err := db.Model(&model.Actress{}).Rows()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var modelActress model.Actress
		db.ScanRows(rows, &modelActress)
		if _, ok := actressMap[modelActress.Actress]; ok {
			delete(actressMap, modelActress.Actress)
		}
	}

	fmt.Printf("map: %+v\n", actressMap)
	for k, _ := range actressMap {
		nameSlice := []rune(k)
		avatarPath := avatarDir + "/" + k + ".png"

		_, err := os.Stat(avatarPath)
		if os.IsNotExist(err) {
			if err := utils.GenerateAvatar(string(nameSlice[0]), avatarPath); err != nil {
				log.Fatal(err)
			}
		}
		data = append(data, model.Actress{
			Actress: k,
			Avatar:  avatarPath,
		})
	}
	fmt.Printf("data: %+v\n", data)

	err = db.CreateInBatches(data, 30).Error
	fmt.Printf("err: %+v\n", err)
}

// 读取文件数据到 map
func ReadFileToMap(name string, v any) error {
	bytes, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bytes, &v); err != nil {
		return err
	}
	return nil
}
