package db

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/wxw9868/video/config"
	"github.com/wxw9868/video/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	var err error
	db, err = gorm.Open(sqlite.Open(config.Config().Database.DBPath), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "video_",                        // table name prefix, table for `User` would be `t_users`
			SingularTable: true,                            // use singular table name, table for `User` would be `user` with this option enabled
			NoLowerCase:   true,                            // skip the snake_casing of names
			NameReplacer:  strings.NewReplacer("ID", "id"), // use name replacer to change struct/field name before convert it to db name
		},
		PrepareStmt: true,
		Logger:      newLogger,
	})
	if err != nil {
		log.Fatalf("数据库链接失败: %s\n", err)
	}

	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("数据库链接失败: %s\n", err)
	}

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("数据库连接成功")

	// 迁移 schema
	if err := db.AutoMigrate(
		&model.User{},
		&model.UserLoginLog{},
		&model.UserCollectLog{},
		&model.UserPageViewsLog{},
		&model.UserCommentLog{},
		&model.Actress{},
		&model.Video{},
		&model.VideoLog{},
		&model.VideoComment{},
		&model.VideoActress{},
		&model.VideoTag{},
		&model.VideoDanmu{},
		&model.Tag{},
		&model.Liquidation{}, &model.TradingRecords{},
	); err != nil {
		log.Fatalf("迁移 schema 失败: %s\n", err)
	}
}

func DB() *gorm.DB {
	return db
}
