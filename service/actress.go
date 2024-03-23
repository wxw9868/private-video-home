package service

import "github.com/wxw9868/video/model"

type ActressService struct{}

type Actress struct {
	ID      uint   `json:"id"`
	Actress string `gorm:"column:actress" json:"actress"`
	Avatar  string `gorm:"column:avatar" json:"avatar"`
	Count   uint32 `json:"count"`
}

func (as *ActressService) Find() ([]Actress, error) {
	var actresss []Actress
	rows, err := db.Model(&model.Actress{}).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var count int64
	for rows.Next() {
		var modelActress model.Actress
		db.ScanRows(rows, &modelActress)
		db.Model(&model.Video{}).Where("actress = ?", modelActress.Actress).Count(&count)

		actresss = append(actresss, Actress{
			ID:      modelActress.ID,
			Actress: modelActress.Actress,
			Avatar:  modelActress.Avatar,
			Count:   uint32(count),
		})
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
