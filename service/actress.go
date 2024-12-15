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
		"西条沙羅": `别名: Sara Saijo, 東条紗奈, 果山ゆら, 西条沙羅, 西條沙羅, 高橋由香利
出生: 1989年11月29日
三围: B95 / W59 / H92
罩杯: H Cup
出道日期: 2015年02月
星座: Sagittarius
血型: A
身高: 160cm
国籍: 日本
简介: 暂无关于西条沙羅(Sara Saijo/35岁)的介绍`,
		"一条リオン": `别名: YURINA, YURINA​, 一条りおん, 山田あゆ, 琴音りあ, 鈴木里緒奈
出生: 1994年06月08日
三围: B95 / W60 / H86
罩杯: E Cup
出道日期: n/a
星座: Gemini
血型: A
身高: 157cm
国籍: 日本
简介: 暂无关于一条リオン(Rion Ichijo/30岁)的介绍。`,
		"松本まりな": `别名: 奥田美須絵
松本まりな(Marina Matsumoto)个人资料:
出生: 1969年06月08日
三围: B82 / W59 / H85
罩杯: C Cup
出道日期: n/a
星座: Gemini
血型: A
身高: 158cm
国籍: 日本
简介: 暂无关于松本まりな(Marina Matsumoto)的介绍。`,
		"成宮はるあ": `别名: 一ノ木ありさ, 乃木はるか, 成宮はるか, 春宮はるな, 東美奈, 芦原亮子, 葉月絢音, 葵律, 陽咲希美
出生: 1992年07月29日
三围: B97 / W58 / H87
罩杯: H Cup
出道日期: n/a
星座: Leo
血型: B
身高: 163cm
国籍: 日本
简介: 暂无关于成宮はるあ(Harua Narumiya/32岁)的介绍。`,
		"咲乃柑菜": `别名: 平井絵里, 蘭華
出生: 1996年06月02日
三围: B82 / W58 / H84
罩杯: C Cup
出道日期: 2015年12月
星座: Gemini
血型: O
身高: 157cm
国籍: 日本
简介: 暂无关于咲乃柑菜(Kanna Sakuno/28岁)的介绍。`,
		"大場ゆい": `别名: 小松なつ, 竹内美羽
出生: 1987年06月06日
三围: B83 / W60 / H86
罩杯: E Cup
出道日期: n/a
星座: Gemini
血型: O
身高: 170cm
国籍: 日本
简介: 暂无关于大場ゆい(Yui Oba)的介绍。`,
		"真田春香": `别名: Haruka Sanada
出生: 1988年06月05日
三围: B95 / W58 / H88
罩杯: G Cup
出道日期: 2006年12月
星座: Gemini
血型: A
身高: 165cm
国籍: 日本
简介: 暂无关于真田春香(Haruka Sanada)的介绍。`,
		"青山未来": `别名: 今井杏樹
出生: 1993年06月18日
三围: B85 / W59 / H90
罩杯: D Cup
出道日期: 2014年02月
星座: Gemini
血型: B
身高: 152cm
国籍: 日本
简介: 暂无关于青山未来(Miku Aoyama/31岁)的介绍。`,
		"木村つな": `别名: 保坂祐美子, 和泉千佳, 大橋みく, 新島優, 日向ももか, 紺野いろは, 霧宮てん
出生: 1993年03月14日
三围: B81 / W58 / H86
罩杯: B Cup
出道日期: 2012年02月
星座: Pisces
血型: A
身高: 149cm
国籍: 日本
简介: 暂无关于木村つな(Tsuna Kimura/31岁)的介绍。`,
		"木村美羽": `出生: 1993年11月29日
三围: 86-58-89 (cm)
罩杯: E Cup
出道日期: 2014年10月
星座: Sagittarius
血型: n/a
身高: 168 cm
国籍: 日本
简介: 暂无关于木村美羽(Miu Kimura/31岁)的介绍。`,
		"香坂澪": `别名: 宇佐美玲奈, 寺野愛, 能登亜実花, 風谷音緒, 風谷音緒・風谷ねお, 香坂香織
出生: 1984年11月25日
三围: B88 / W59 / H85
罩杯: D Cup
出道日期: n/a
星座: Sagittarius
血型: O
身高: 165cm
国籍: 日本
简介: 暂无关于香坂澪(Mio Kosaka)的介绍。`,
		"藤咲りさ": `出生: 1986年04月10日
三围: B82 / W58 / H82
罩杯: C Cup
出道日期: 2006年12月
星座: Aries
血型: O
身高: 160cm
国籍: 日本
简介: 暂无关于藤咲りさ(Risa Fujisaki)的介绍。`,
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
