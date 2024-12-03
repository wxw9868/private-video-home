package service

import (
	"context"

	sqlite "github.com/wxw9868/video/initialize/db"
	redis "github.com/wxw9868/video/initialize/rdb"
	"gorm.io/gorm"
)

var (
	db  = sqlite.DB()
	rdb = redis.Rdb()
	ctx = context.Background()
)

func Paginate(page, pageSize, count int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > count:
			pageSize = count
		case pageSize <= 0:
			pageSize = 20
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// ResetTable 使用联合索引解决问题
func ResetTable(table string) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 删除数据
		if err := tx.Exec("DELETE FROM ?", table).Error; err != nil {
			return err
		}
		// 重置主键
		if err := tx.Exec("UPDATE SQLITE_SEQUENCE SET seq = 0 WHERE name = ?", table).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
