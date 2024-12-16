package service

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
)

type ActressService struct{}

func (as *ActressService) Add(name string) error {
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
		"輝月あんり": `别名: Anri Kizuki, あんり（専門学生）, 天木ゆう, 輝月あんり
出生: 1994年02月03日
三围: B76 / W59 / H87
罩杯: C Cup
出道日期: n/a
星座: Aquarius
血型: A
身高: 164cm
国籍: 日本
简介: 暂无关于輝月あんり(Anri Kizuki/30岁)的介绍。`,
		"麻生希": `出生: 1988年12月19日
三围: B88 / W58 / H89
罩杯: E Cup
出道日期: 2016年06月
星座: Sagittarius
血型: A
身高: 170cm
国籍: 日本
简介: 暂无关于麻生希(Nozomi Aso)的介绍。`,
		"京野ななか": `别名: Emiri Aizawa, Nana Koizumi, あやめさくら, かな（居酒屋）, 小泉奈々, 早坂かな, 相沢えみり
出生: 1992年06月07日
三围: B84 / W60 / H88
罩杯: D Cup
出道日期: n/a
星座: Gemini
血型: A
身高: 158cm
国籍: 日本
简介: 暂无关于京野ななか(Nanaka Kyono/32岁)的介绍。`,
		"あざみねね": `别名: 伊藤英玲奈, 吉澤留美, 春日もな, 立花えれな, 鈴木ワカ, 鈴木ワコ
出生: 1990年10月10日
三围: B91 / W58 / H85
罩杯: J Cup
出道日期: n/a
星座: Libra
血型: B
身高: 153cm
国籍: 日本
简介: 暂无关于あざみねね(Nene Asami/34岁)的介绍。`,
		"杏堂なつ": `别名: 安藤なつみ, 松山なつみ, 榊れいな
出生: 1987年12月22日
三围: B93 / W59 / H87
罩杯: n/a
出道日期: 2006年12月
星座: Capricorn
血型: n/a
身高: 163cm
国籍: 日本
简介: 暂无关于杏堂なつ(Natsu Ando)的介绍。`,
		"七瀬ななみ": `别名: Nanami Nagase
出生: 1981年08月27日
三围: 83-60-87 (cm)
罩杯: C-70 Cup
出道日期: n/a
星座: Virgo
血型: A
身高: 160 cm
国籍: 日本
简介: 暂无关于七瀬ななみ(Nanami Nanase)的介绍。`,
		"雨音わかな": `别名: 奥野光香
出生: n/a
三围: B89 / W59 / H90
罩杯: F Cup
出道日期: 2016年06月
星座: n/a
血型: A
身高: 166cm
国籍: 日本
简介: 暂无关于雨音わかな(Wakana Amane)的介绍。`,
		//"": ``,
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for actress, str := range strs {
			m := model.Actress{}
			m.Avatar = "assets/image/avatar/" + actress + ".jpg"
			results := strings.Split(str, "\n")
			for _, result := range results {
				arr := strings.Split(result, ":")
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
					m.Introduction = value
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
