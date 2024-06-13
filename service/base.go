package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	sqlite "github.com/wxw9868/video/initialize/db"
	"github.com/wxw9868/video/initialize/httpclient"
	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
	"gorm.io/gorm"
)

var db = sqlite.DB()
var mutex = new(sync.Mutex)

func Paginate(page, pageSize, count int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > count:
			pageSize = count
		case pageSize <= 0:
			pageSize = 16
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func VideoImport(videoDir string) error {
	files, err := os.ReadDir(videoDir)
	if err != nil {
		return err
	}

	var avatarDir = "./assets/image/avatar"
	var posterDir = "./assets/image/poster"
	var actressList = make(map[string]struct{})
	var videoSql = "INSERT OR REPLACE INTO video_Video (title, actress, size, duration, poster, width, height, codec_name, channel_layout, creation_time, CreatedAt, UpdatedAt) VALUES "
	var actressSql = "INSERT OR REPLACE INTO video_Actress (actress, avatar, CreatedAt, UpdatedAt) VALUES "

	for _, file := range files {
		filename := file.Name()
		ext := filepath.Ext(filename)
		if ext == ".mp4" {
			title := strings.Split(filename, ".")[0]
			arr := strings.Split(title, "_")
			actress := arr[len(arr)-1]

			mutex.Lock()
			if _, ok := actressList[actress]; !ok {
				actressList[actress] = struct{}{}
			}
			mutex.Unlock()

			filePath := videoDir + "/" + filename
			posterPath := posterDir + "/" + title + ".jpg"
			_, err = os.Stat(posterPath)
			if os.IsNotExist(err) {
				if err = utils.ReadFrameAsJpeg(filePath, posterPath, "00:1:58"); err != nil {
					return err
				}
			}
			videoInfo, err := utils.VideoInfo(filePath)
			if err != nil {
				return err
			}

			videoSql += fmt.Sprintf("('%s', '%s', %d, %f, '%s', %d, %d, '%s', '%s', '%v', '%v', '%v'), ", title, actress, videoInfo["size"].(int64), videoInfo["duration"].(float64), posterPath, videoInfo["width"].(int64), videoInfo["height"].(int64), fmt.Sprintf("%s,%s", videoInfo["codec_name0"].(string), videoInfo["codec_name1"].(string)), videoInfo["channel_layout"].(string), videoInfo["creation_time"].(time.Time), time.Now().Local(), time.Now().Local())
		}
	}

	if len(actressList) > 0 {
		for actress, _ := range actressList {
			avatarPath := avatarDir + "/" + actress + ".png"

			_, err := os.Stat(avatarPath)
			if os.IsNotExist(err) {
				nameSlice := []rune(actress)
				if err := utils.GenerateAvatar(string(nameSlice[0]), avatarPath); err != nil {
					return err
				}
			}

			actressSql += fmt.Sprintf("('%s', '%s', '%v', '%v'), ", actress, avatarPath, time.Now().Local(), time.Now().Local())
		}
	}

	videoSqlBytes := []byte(videoSql)
	actressSqlBytes := []byte(actressSql)
	videoSql = string(videoSqlBytes[:len(videoSqlBytes)-2])
	actressSql = string(actressSqlBytes[:len(actressSqlBytes)-2])

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(videoSql).Error; err != nil {
			return err
		}
		if err := tx.Exec(actressSql).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func ImportActress() error {
	var avatarDir = "./assets/image/avatar"
	var actressMap = make(map[string]struct{})

	utils.ReadFileToMap("actress.json", &actressMap)
	// fmt.Printf("map: %+v\n", actressMap)

	var actressSql = "INSERT OR REPLACE INTO video_Actress (actress, avatar, CreatedAt, UpdatedAt) VALUES "
	for actress, _ := range actressMap {
		avatarPath := avatarDir + "/" + actress + ".png"
		_, err := os.Stat(avatarPath)
		if os.IsNotExist(err) {
			nameSlice := []rune(actress)
			if err := utils.GenerateAvatar(string(nameSlice[0]), avatarPath); err != nil {
				return err
			}
		}
		actressSql += fmt.Sprintf("('%s', '%s', '%v', '%v'), ", actress, avatarPath, time.Now().Local(), time.Now().Local())
	}
	actressSqlBytes := []byte(actressSql)
	actressSql = string(actressSqlBytes[:len(actressSqlBytes)-2])

	if err := db.Exec(actressSql).Error; err != nil {
		return err
	}
	return nil
}

func Post(url string, body io.Reader) error {
	resp, err := httpclient.HttpClient().POST(url, "application/json", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	robots, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(robots))
	return nil
}

type VideoActressData struct {
	VideoID   uint   `json:"video_id" gorm:"column:video_id"`
	Actress   string `json:"actress" gorm:"column:actress"`
	ActressID uint   `json:"actress_id" gorm:"column:actress_id"`
}

func VideoActress() error {
	var sql = "INSERT OR REPLACE INTO video_VideoActress (video_id, actress_id, CreatedAt, UpdatedAt) VALUES "
	var actressData []model.Actress
	var videos []model.Video
	var actressWhere []string

	dbVideo := db.Table("video_Video as v")
	dbVideo = dbVideo.Select("v.id as video_id, v.actress, a.id as actress_id").Joins("left join video_Actress a on a.actress = v.actress")
	rows, err := dbVideo.Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var data VideoActressData
		db.ScanRows(rows, &data)

		sql += fmt.Sprintf("(%d, %d, '%v', '%v'), ", data.VideoID, data.ActressID, time.Now().Local(), time.Now().Local())
		actressWhere = append(actressWhere, data.Actress)
	}

	db.Where("actress NOT IN ?", actressWhere).Find(&actressData)
	// fmt.Printf("%+v\n", actressData)
	if len(actressData) > 0 {
		for _, actress := range actressData {
			db.Where("title LIKE ?", "%"+actress.Actress+"%").Find(&videos)
			if len(videos) > 0 {
				for _, video := range videos {
					sql += fmt.Sprintf("(%d, %d, '%v', '%v'), ", video.ID, actress.ID, time.Now().Local(), time.Now().Local())
				}
			}
		}
	}

	sqlBytes := []byte(sql)
	sql = string(sqlBytes[:len(sqlBytes)-2])
	if err := db.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}
