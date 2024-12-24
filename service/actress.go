package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/gocolly/colly/v2"
	"github.com/redis/go-redis/v9"
	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/model/request"
	"github.com/wxw9868/video/utils"
	"gorm.io/gorm"
)

type ActressService struct{}

func (as *ActressService) Create(req request.CreateActress) error {
	if req.Avatar == "" {
		req.Avatar = "assets/image/avatar/anonymous.png"
	}
	result := db.Where(model.Actress{Actress: req.Name}).FirstOrCreate(&model.Actress{
		Actress:      req.Name,
		Alias:        req.Alias,
		Avatar:       req.Avatar,
		Birth:        req.Birth,
		Measurements: req.Measurements,
		CupSize:      req.CupSize,
		DebutDate:    req.DebutDate,
		StarSign:     req.StarSign,
		BloodGroup:   req.BloodGroup,
		Stature:      req.Stature,
		Nationality:  req.Nationality,
		Intro:        req.Intro,
	})
	if result.RowsAffected == 1 {
		return nil
	}
	return errors.New("演员存在")
}

func (as *ActressService) Updates(req request.UpdateActress) error {
	actress := model.Actress{
		Actress:      req.Name,
		Alias:        req.Alias,
		Avatar:       req.Avatar,
		Birth:        req.Birth,
		Measurements: req.Measurements,
		CupSize:      req.CupSize,
		DebutDate:    req.DebutDate,
		StarSign:     req.StarSign,
		BloodGroup:   req.BloodGroup,
		Stature:      req.Stature,
		Nationality:  req.Nationality,
		Intro:        req.Intro,
	}
	if err := db.Model(&model.Actress{}).Where("id = ?", req.Id).Updates(actress).Error; err != nil {
		return err
	}
	return nil
}

func (as *ActressService) Delete(id uint) error {
	if err := db.Delete(&model.Actress{}, id).Error; err != nil {
		return err
	}
	return nil
}

type Actress struct {
	ID      uint   `json:"id"`
	Actress string `gorm:"column:actress" json:"actress"`
	Avatar  string `gorm:"column:avatar" json:"avatar"`
	Count   uint32 `gorm:"column:count" json:"count"`
}

func (as *ActressService) List(req request.SearchActress) (data map[string]interface{}, err error) {
	var ids []uint

	f := func(ids []uint, totalCount int) (map[string]interface{}, error) {
		actresss := make([]Actress, len(ids))
		for i, id := range ids {
			result := rdb.HGetAll(ctx, utils.Join("video_actress_", strconv.Itoa(int(id)))).Val()
			count, err := strconv.Atoi(result["count"])
			if err != nil {
				return nil, err
			}
			actresss[i] = Actress{
				ID:      id,
				Actress: result["actress"],
				Avatar:  result["avatar"],
				Count:   uint32(count),
			}
		}
		return map[string]interface{}{"list": actresss, "total": totalCount}, nil
	}

	if req.Actress != "" {
		var adb = db.Model(&model.Actress{}).Where("actress like ?", "%"+req.Actress+"%")
		var count int64
		if err = adb.Count(&count).Error; err != nil {
			return nil, err
		}
		adb.Pluck("id", &ids)
		return f(ids, int(count))
	}

	var key string
	var sql = "SELECT a.id, a.actress, a.avatar, count(va.video_id) as count FROM video_Actress a left join video_VideoActress va on a.id = va.actress_id group by 1,2,3"
	if req.Column != "" && req.Order != "" {
		sql += utils.Join(" order by ", req.Column, " ", req.Order)
	}

	var totalCount int64
	var adb = db.Model(&model.Actress{})
	switch req.Column {
	case "a.CreatedAt":
		adb = adb.Order(utils.Join("CreatedAt", " ", req.Order))
		key = "video_actress_createdAt"
	case "a.actress":
		adb = adb.Order(utils.Join("actress", " ", req.Order))
		key = "video_actress_actress"
	case "count":
		adb = db.Table("(?)", db.Raw(sql))
		if err = adb.Count(&totalCount).Error; err != nil {
			return nil, err
		}
		if err = adb.Pluck("id", &ids).Error; err != nil {
			return nil, err
		}
		key = "video_actress_count"
	default:
		key = "video_actress"
	}

	if req.Column != "count" {
		if err = adb.Count(&totalCount).Error; err != nil {
			return nil, err
		}
		if err = adb.Pluck("id", &ids).Error; err != nil {
			return nil, err
		}
	}

	bytes, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}
	bytes = []byte{}
	result, _ := rdb.HGet(ctx, key, "ids").Result()
	if strings.Compare(string(bytes), result) == 0 && result != "" {
		return f(ids, int(totalCount))
	}

	var count int64
	if err = db.Model(&model.Actress{}).Count(&count).Error; err != nil {
		return nil, err
	}

	var actresss []Actress
	if err = db.Raw(sql).Scopes(Paginate(req.Page, req.Size, int(count))).Scan(&actresss).Error; err != nil {
		return nil, err
	}

	keys := make([]string, len(actresss)+1)
	keys[0] = key
	for i, a := range actresss {
		keys[i+1] = utils.Join("video_actress_", strconv.Itoa(int(a.ID)))
	}

	txf := func(tx *redis.Tx) error {
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HSet(ctx, key, "len", len(ids), "ids", string(bytes))
			for _, a := range actresss {
				pipe.HSet(ctx, utils.Join("video_actress_", strconv.Itoa(int(a.ID))), "id", a.ID, "actress", a.Actress, "avatar", a.Avatar, "count", a.Count)
			}
			return nil
		})
		return err
	}
	if err = rdb.Watch(ctx, txf, keys...); errors.Is(err, redis.TxFailedErr) {
		return nil, err
	}

	return map[string]interface{}{"list": actresss, "total": count}, nil
}

func (as *ActressService) Info(id uint) (*model.Actress, error) {
	var actress model.Actress
	if err := db.First(&actress, id).Error; err != nil {
		return nil, err
	}
	return &actress, nil
}

// SaveActress 补充信息
func (as *ActressService) SaveActress() error {
	var strs = map[string]string{
		//		"小野静香": `别名: 秋野早苗
		//出生: 1995年07月31日
		//三围: B80 / W60 / H87
		//罩杯: E Cup
		//出道日期: 2015年01月
		//星座: Leo
		//血型: n/a
		//身高: 154
		//国籍: 日本
		//简介: 暂无关于秋野早苗(Sanae Akino/29岁)的介绍。`,
		//"": ``,
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for actress, str := range strs {
			m := model.Actress{}
			m.Avatar = "assets/image/avatar/" + actress + ".jpg"
			results := strings.Split(str, "\n")
			for _, result := range results {
				arr := strings.Split(result, ":")
				fmt.Printf("actress: %s,arr: %+v\n", actress, arr)
				column := strings.Trim(arr[0], " ")
				value := strings.Replace(strings.Trim(arr[1], " "), "n/a", "", -1)
				switch column {
				case "别名":
					m.Alias = value
				case "出生":
					m.Birth = value
				case "三围":
					m.Measurements = value
				case "罩杯":
					m.CupSize = value
				case "出道日期":
					m.DebutDate = value
				case "星座":
					m.StarSign = value
				case "血型":
					m.BloodGroup = value
				case "身高":
					m.Stature = value
				case "国籍":
					m.Nationality = value
				case "简介":
					m.Intro = value
				}
			}
			// 根据 `struct` 更新属性，只会更新非零值的字段
			if err := tx.Model(&model.Actress{}).Where("actress = ?", actress).Updates(m).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (as *ActressService) DownAvatar() error {
	c := colly.NewCollector(
		colly.UserAgent(browser.Random()),
		colly.AllowedDomains("javmenu.com"),
	)

	c.OnHTML(".model", func(e *colly.HTMLElement) {
		src, _ := e.DOM.Find("img").Attr("src")
		name := e.DOM.Find(".model_name").Text()
		fmt.Printf("actress: %s, src:%s, ext:%s\n", name, src, path.Ext(src))

		savePath := "assets/image/pic"
		saveFile := utils.Join(name, path.Ext(src))
		err := utils.DownloadImage(src, savePath, saveFile)
		fmt.Printf("error: %s\n", err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Response %s: %d bytes\n", r.Request.URL, len(r.Body))
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error %s: %v\n", r.Request.URL, err)
	})

	var actresses []string
	if err := db.Model(&model.Actress{}).Pluck("actress", &actresses).Error; err != nil {
		return err
	}

	for _, actress := range actresses {
		err := c.Visit(utils.Join("https://ggjav.com/main/search?string=", url.QueryEscape(actress)))
		if err != nil {
			return err
		}
	}
	return nil
}

func (as *ActressService) CopyAvatar() error {
	var actresses []string
	if err := db.Model(&model.Actress{}).Where("avatar = ?", "assets/image/avatar/defaultgirl.png").Pluck("actress", &actresses).Error; err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, actress := range actresses {
			_, err := os.Stat(path.Join("D:/video/assets/image/pic", utils.Join(actress, ".jpg")))
			if err == nil {
				m := model.Actress{}
				m.Avatar = "assets/image/avatar/" + actress + ".jpg"
				if err = tx.Model(&model.Actress{}).Where("actress = ?", actress).Updates(m).Error; err != nil {
					return err
				}
				err = os.Rename(path.Join("D:/video/assets/image/pic", utils.Join(actress, ".jpg")), path.Join("D:/video/assets/image/na", utils.Join(actress, ".jpg")))
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
