package service

import "github.com/wxw9868/video/model"

type ActressService struct{}

type Actress struct {
	ID      uint   `json:"id"`
	Actress string `json:"actress"`
	Avatar  string `json:"avatar"`
}

func (as *ActressService) Find() ([]Actress, error) {
	var actresss []Actress
	if err := db.Model(&model.Actress{}).Unscoped().Find(&actresss).Error; err != nil {
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
