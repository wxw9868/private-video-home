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
	"github.com/wxw9868/video/utils"
	"gorm.io/gorm"
)

type ActressService struct{}

func (as *ActressService) Create(name string) error {
	result := db.Where(model.Actress{Actress: name}).FirstOrCreate(&model.Actress{Actress: name, Avatar: "assets/image/avatar/anonymous.png"})
	if result.RowsAffected == 1 {
		return nil
	}
	return errors.New("演员存在")
}

func (as *ActressService) Updates(id uint, name string) error {
	var actress model.Actress
	actress.ID = id
	if err := db.Model(&actress).Updates(model.Actress{Actress: name}).Error; err != nil {
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

func (as *ActressService) List(page, pageSize int, action, sort, actress string) (data map[string]interface{}, err error) {
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
		data = map[string]interface{}{
			"list":  actresss,
			"count": totalCount,
		}
		return data, nil
	}

	if actress != "" {
		var adb = db.Model(&model.Actress{}).Where("actress like ?", "%"+actress+"%")
		var count int64
		if err = adb.Count(&count).Error; err != nil {
			return nil, err
		}
		adb.Pluck("id", &ids)
		return f(ids, int(count))
	}

	var key string
	var sql = "SELECT a.id, a.actress, a.avatar, count(va.video_id) as count FROM video_Actress a left join video_VideoActress va on a.id = va.actress_id group by 1,2,3"
	if action != "" && sort != "" {
		sql += utils.Join(" order by ", action, " ", sort)
	}

	var totalCount int64
	var adb = db.Model(&model.Actress{})
	switch action {
	case "a.CreatedAt":
		adb = adb.Order(utils.Join("CreatedAt", " ", sort))
		key = "video_actress_createdAt"
	case "a.actress":
		adb = adb.Order(utils.Join("actress", " ", sort))
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

	if action != "count" {
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
	if err = db.Raw(sql).Scopes(Paginate(page, pageSize, int(count))).Scan(&actresss).Error; err != nil {
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

	data = map[string]interface{}{
		"list":  actresss,
		"count": count,
	}

	return data, nil
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
		"望月しおん": `别名: 桜井ひなた,井坂由希恵
三围: B80 / W55 / H82
身高: 152cm`,
		"橘ひなた": `别名: 小日向ひかり, 牧瀬美央, 高橋陽子
出生: 1990年08月11日
三围: B83 / W60 / H86
罩杯: D Cup
出道日期: n/a
星座: Leo
血型: A
身高: 160cm
国籍: 日本
简介: 暂无关于橘ひなた(Hinata Tachibana/34岁)的介绍。`,
		"中原あきな": `别名: あまねなのは
出生: 1990年02月20日
三围: B84 / W58 / H84
罩杯: E Cup
出道日期: 2008年09月
星座: Pisces
血型: n/a
身高: n/a
国籍: 日本
简介: 暂无关于あまねなのは(Nanoha Amane/34岁)的介绍。`,
		"小野麻里亜": `别名: Maria Ono, 中村真理亜, 小野まり, 小野まりあ, 小野麻理亜, 尾崎ゆりあ, 椎名綺更, 神代凛, 神代凜
出生: 1989年03月16日
三围: B86 / W60 / H87
罩杯: E Cup
出道日期: n/a
星座: Pisces
血型: B
身高: 150cm
国籍: 日本
简介: 暂无关于小野麻里亜(Maria Ono/35岁)的介绍。`,
		"渋谷まなか": `别名: 彩葉みおり, 豊田愛菜, 川村里穂, 大橋愛菜
出生: 1995年05月11日
三围: B95 / W60 / H87
罩杯: H Cup
出道日期: 2017年11月
星座: Taurus
血型: n/a
身高: 168 cm
国籍: 日本
简介: 暂无关于豊田愛菜(Mana Toyota/29岁)的介绍。`,
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
