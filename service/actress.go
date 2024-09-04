package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
	"strconv"
	"strings"
)

type ActressService struct{}

func (as *ActressService) Add(name string) error {
	result := db.Where(model.Actress{Actress: name}).FirstOrCreate(&model.Actress{Actress: name})
	if result.RowsAffected == 1 {
		return nil
	}
	return errors.New("头像已存在")
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
	var actresss []Actress
	var ids []uint
	var key string
	var sql string
	ctx := context.Background()

	selectSQL := "SELECT a.id, a.actress, a.avatar, count(va.video_id) as count FROM video_Actress a left join video_VideoActress va on a.id = va.actress_id"
	groupSQL := " group by 1,2,3"

	if actress != "" {
		sql = utils.Join(" where a.actress = ", "'", actress, "'")
		if err := db.Raw(utils.Join(selectSQL, sql, groupSQL)).Scan(&actresss).Error; err != nil {
			return nil, err
		}
		return actresss, nil
	}

	if action == "null" || sort == "null" {
		action = ""
		sort = ""
	}
	if action != "" && sort != "" {
		sql = utils.Join(" order by ", action, " ", sort)
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
		if err := db.Table("(?)", db.Raw(utils.Join(selectSQL, groupSQL, sql))).Pluck("id", &ids).Error; err != nil {
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
	result, _ := rdb.HGet(ctx, key, "ids").Result()
	//fmt.Println("compare: ", strings.Compare(string(bytes), result))
	if strings.Compare(string(bytes), result) == 0 && result != "" {
		for _, id := range ids {
			data := rdb.HGetAll(ctx, utils.Join("video_actress_", strconv.Itoa(int(id)))).Val()
			count, _ := strconv.Atoi(data["count"])
			actresss = append(actresss, Actress{
				ID:      id,
				Actress: data["actress"],
				Avatar:  data["avatar"],
				Count:   uint32(count),
			})
		}
		return actresss, nil
	}

	var count int64
	if err = db.Model(&model.Actress{}).Count(&count).Error; err != nil {
		return nil, err
	}

	if err = db.Raw(utils.Join(selectSQL, groupSQL, sql)).Scopes(Paginate(page, pageSize, int(count))).Scan(&actresss).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = rdb.HSet(ctx, key, "len", len(ids), "ids", string(bytes)).Err()
	if err != nil {
		return nil, err
	}

	for _, a := range actresss {
		rdb.HSet(ctx, utils.Join("video_actress_", strconv.Itoa(int(a.ID))), "id", a.ID, "actress", a.Actress, "avatar", a.Avatar, "count", a.Count)
	}

	return actresss, nil
}
