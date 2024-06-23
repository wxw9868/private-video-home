package service

import (
	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
)

type ActressService struct{}

type Actress struct {
	ID      uint   `json:"id"`
	Actress string `gorm:"column:actress" json:"actress"`
	Avatar  string `gorm:"column:avatar" json:"avatar"`
	Count   uint32 `gorm:"column:count" json:"count"`
}

func (as *ActressService) Find(page, pageSize int, action, sort string) ([]Actress, error) {
	var actresss []Actress
	sql := "SELECT a.id, a.actress, a.avatar, count(va.video_id) as count FROM video_Actress a left join video_VideoActress va on a.id = va.actress_id group by 1,2,3"
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
	return actresss, nil
}

func (as *ActressService) List() ([]model.Actress, error) {
	var actresss []model.Actress
	if err := db.Find(&actresss).Error; err != nil {
		return nil, err
	}
	return actresss, nil
}

func (as *ActressService) Create(actresss []model.Actress) error {
	if err := db.Create(&actresss).Error; err != nil {
		return err
	}
	return nil
}
