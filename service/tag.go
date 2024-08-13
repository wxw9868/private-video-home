package service

import (
	"errors"

	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
)

type TagService struct{}

func (ts *TagService) Add(name string) error {
	result := db.Where(model.Tag{Name: name}).FirstOrCreate(&model.Tag{Name: name})
	if result.RowsAffected == 1 {
		return nil
	}
	return errors.New("标签已存在")
}

func (ts *TagService) Edit(id uint, name string) error {
	var tag model.Tag
	tag.ID = id
	if err := db.Model(&tag).Updates(model.Tag{Name: "hello"}).Error; err != nil {
		return err
	}
	return nil
}

func (ts *TagService) Delete(id uint) error {
	if err := db.Delete(&model.Tag{}, id).Error; err != nil {
		return err
	}
	return nil
}

type Tags struct {
	ID     uint   `json:"id"`
	Name   string `gorm:"column:name" json:"name"`
	Avatar string `gorm:"column:avatar" json:"avatar"`
	Count  uint32 `gorm:"column:count" json:"count"`
}

func (as *TagService) List(page, pageSize int, action, sort string) ([]Actress, error) {
	var actresss []Actress
	sql := "SELECT t.id, t.name, t.avatar, count(vt.video_id) as count FROM video_Tag t left join video_VideoTag vt on t.id = vt.tag_id group by 1,2,3"
	if action != "" && sort != "" {
		sql += utils.Join(" order by ", action, " ", sort)
	}

	var count int64
	if err := db.Table("video_Tag t").Count(&count).Error; err != nil {
		return nil, err
	}

	if err := db.Raw(sql).Scopes(Paginate(page, pageSize, int(count))).Scan(&actresss).Error; err != nil {
		return nil, err
	}
	return actresss, nil
}
