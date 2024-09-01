package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wxw9868/video/initialize/rdb"
	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
	"strconv"
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
	sql := "SELECT a.id, a.actress, a.avatar, count(va.video_id) as count FROM video_Actress a left join video_VideoActress va on a.id = va.actress_id"
	if actress != "" {
		sql += utils.Join(" where a.actress = ", "'", actress, "'")
	}
	sql += " group by 1,2,3"
	if action == "null" || sort == "null" {
		action = ""
		sort = ""
	}
	if action != "" && sort != "" {
		sql += utils.Join(" order by ", action, " ", sort)
	}

	var count int64
	if err := db.Table("video_Actress a").Count(&count).Error; err != nil {
		return nil, err
	}

	if err := db.Raw(sql).Scopes(Paginate(page, pageSize, int(count))).Scan(&actresss).Error; err != nil {
		return nil, err
	}

	var ids []uint
	for _, a := range actresss {
		ids = append(ids, a.ID)
		rdb.Rdb().HSet(context.Background(), utils.Join("video_actress_", strconv.Itoa(int(a.ID))), "id", a.ID, "actress", a.Actress, "avatar", a.Avatar, "count", a.Count)
	}
	bytes, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}
	err = rdb.Rdb().HSet(context.Background(), "video_actress", "len", len(actresss), "ids", bytes).Err()
	if err != nil {
		return nil, err
	}
	
	return actresss, nil
}
