package service

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
)

type ActressService struct{}

func (as *ActressService) Add(name string) error {
	result := db.Where(model.Actress{Actress: name, Avatar: "assets/image/avatar/anonymous.png"}).FirstOrCreate(&model.Actress{Actress: name})
	if result.RowsAffected == 1 {
		return nil
	}
	return errors.New("演员存在")
}

func (as *ActressService) Edit(id uint, name string) error {
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

func (as *ActressService) List(page, pageSize int, action, sort, actress string) ([]Actress, error) {
	var ids []uint

	f := func(ids []uint) ([]Actress, error) {
		actresss := make([]Actress, len(ids))
		for i, id := range ids {
			data := rdb.HGetAll(ctx, utils.Join("video_actress_", strconv.Itoa(int(id)))).Val()
			count, err := strconv.Atoi(data["count"])
			if err != nil {
				return nil, err
			}
			actresss[i] = Actress{
				ID:      id,
				Actress: data["actress"],
				Avatar:  data["avatar"],
				Count:   uint32(count),
			}
		}
		return actresss, nil
	}

	if actress != "" {
		db.Model(&model.Actress{}).Where("actress = ?", actress).Pluck("id", &ids)
		return f(ids)
	}

	var key string
	var sql = "SELECT a.id, a.actress, a.avatar, count(va.video_id) as count FROM video_Actress a left join video_VideoActress va on a.id = va.actress_id group by 1,2,3"
	if action != "" && sort != "" {
		sql += utils.Join(" order by ", action, " ", sort)
	}

	switch action {
	case "a.CreatedAt":
		if err := db.Model(&model.Actress{}).Order(utils.Join("CreatedAt", " ", sort)).Pluck("id", &ids).Error; err != nil {
			return nil, err
		}
		key = "video_actress_createdAt"
	case "a.actress":
		if err := db.Model(&model.Actress{}).Order(utils.Join("actress", " ", sort)).Pluck("id", &ids).Error; err != nil {
			return nil, err
		}
		key = "video_actress_actress"
	case "count":
		if err := db.Table("(?)", db.Raw(sql)).Pluck("id", &ids).Error; err != nil {
			return nil, err
		}
		key = "video_actress_count"
	default:
		if err := db.Model(&model.Actress{}).Pluck("id", &ids).Error; err != nil {
			return nil, err
		}
		key = "video_actress"
	}

	bytes, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}
	//bytes = []byte{}
	result, _ := rdb.HGet(ctx, key, "ids").Result()
	if strings.Compare(string(bytes), result) == 0 && result != "" {
		return f(ids)
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
	if err = rdb.Watch(ctx, txf, keys...); err == redis.TxFailedErr {
		return nil, err
	}
	return actresss, nil
}
