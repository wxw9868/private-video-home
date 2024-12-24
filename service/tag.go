package service

import (
	"errors"

	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
)

type TagService struct{}

func (ts *TagService) Create(name string) error {
	result := db.Where(model.Tag{Name: name}).FirstOrCreate(&model.Tag{Name: name})
	if result.RowsAffected == 1 {
		return nil
	}
	return errors.New("标签已存在")
}

func (ts *TagService) List(page, pageSize int, column, order string) (map[string]interface{}, error) {
	rdb := db.Model(&model.Tag{})
	if column != "" && order != "" {
		rdb = rdb.Order(utils.Join(column, " ", order))
	}

	var total int64
	if err := rdb.Count(&total).Error; err != nil {
		return nil, err
	}

	var tags []model.Tag
	if err := rdb.Scopes(Paginate(page, pageSize, int(total))).Find(&tags).Error; err != nil {
		return nil, err
	}
	return map[string]interface{}{"list": tags, "total": total}, nil
}

func (ts *TagService) Info(id uint) (*model.Tag, error) {
	var tag model.Tag
	if err := db.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (ts *TagService) Update(id uint, name string) error {
	if err := db.Model(&model.Tag{}).Where("id = ?", id).Updates(model.Tag{Name: name}).Error; err != nil {
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
