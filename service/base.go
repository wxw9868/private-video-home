package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	sqlite "github.com/wxw9868/video/initialize/db"
	gofoundClient "github.com/wxw9868/video/initialize/gofound"
	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
	"gorm.io/gorm"
)

var db = sqlite.DB()
var gofoundCount = 0
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
			pageSize = 1000
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

type VideoActressData struct {
	VideoID   uint   `json:"video_id" gorm:"column:video_id"`
	Actress   string `json:"actress" gorm:"column:actress"`
	ActressID uint   `json:"actress_id" gorm:"column:actress_id"`
}

// 使用联合索引解决问题
func VideoActress() error {
	var sql = "INSERT OR REPLACE INTO video_VideoActress (video_id, actress_id, CreatedAt, UpdatedAt) VALUES "
	var actresss []model.Actress
	var videos []model.Video

	if err := db.Find(&actresss).Error; err != nil {
		return err
	}
	// fmt.Printf("%+v\n", actressData)
	if len(actresss) > 0 {
		for _, actress := range actresss {
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

	err := db.Transaction(func(tx *gorm.DB) error {
		// 删除数据
		if err := tx.Exec("DELETE FROM video_VideoActress").Error; err != nil {
			return err
		}
		// 重置主键
		if err := tx.Exec("UPDATE SQLITE_SEQUENCE SET seq = 0 WHERE name = 'video_VideoActress'").Error; err != nil {
			return err
		}
		if err := tx.Exec(sql).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func VideoActressAdditionalInformation() error {
	var actresss []model.Actress
	if err := db.Find(&actresss).Error; err != nil {
		return err
	}
	for _, v := range actresss {
		m := v

		elems := make([]string, 3)
		elems[0] = "https://xslist.org/search?query="
		elems[1] = v.Actress
		elems[2] = "&lg=zh"
		doc, err := utils.GetWebDocument(strings.Join(elems, ""))
		if err != nil {
			return err
		}
		href, _ := doc.Find("a").Attr("href")

		doc, err = utils.GetWebDocument(href)
		if err != nil {
			return err
		}
		sss1 := doc.Find("#sss1")
		// actress := sss1.Find("header").Text()
		alias := sss1.Find("p").Text()
		avatar, _ := sss1.Find("img").Attr("src")

		saveFile := ""
		utils.DownloadImage(avatar, saveFile)

		m.Alias = alias
		m.Avatar = avatar
		doc.Find("h2").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				personal, _ := s.Next().Html()
				personal = strings.Replace(strings.Replace(strings.Replace(personal, "<span itemprop=\"height\">", "", -1), "<span itemprop=\"nationality\">", "", -1), "</span>", "", -1)
				personals := strings.Split(personal, "<br/>")
				birth := personals[0]                  // 出生
				measurements := personals[1]           // 三围
				cup_size := personals[2]               // 罩杯
				debut_date := personals[3]             // 出道日期
				star_sign := personals[4]              // 星座
				blood_group := personals[5]            // 血型
				stature := personals[6]                // 身高
				nationality := personals[7]            // 国籍
				introduction := s.Next().Next().Text() // 简介

				m.Birth = birth
				m.Measurements = measurements
				m.CupSize = cup_size
				m.DebutDate = debut_date
				m.StarSign = star_sign
				m.BloodGroup = blood_group
				m.Stature = stature
				m.Nationality = nationality
				m.Introduction = introduction
			}
		})
		if err := db.Save(&m).Error; err != nil {
			return err
		}
	}

	return nil
}
