package config

import (
	"errors"
	"log"
	"path"
	"path/filepath"
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var c config

func Config() *config {
	return &c
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path.Join(AbsPath("")))
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Panic("Fatal error config file: ", err)
		}
	}
	unmarshal()

	viper.OnConfigChange(func(e fsnotify.Event) {
		unmarshal()
		log.Println("Config file changed:", e.String())
	})
	viper.WatchConfig()
}

func AbsPath(path string) string {
	if path == "" {
		_, filename, _, _ := runtime.Caller(0)
		return filepath.Dir(filename)
	} else {
		s, _ := filepath.Abs(path)
		if filepath.IsAbs(s) {
			return s
		} else {
			return ""
		}
	}
}

func unmarshal() {
	if err := viper.Unmarshal(&c); err != nil {
		log.Panic("unable to decode into struct,", err)
	}
}

type config struct {
	System   system
	Database database
	Redis    redis
	Gofound  gofound
	Alisms   alisms
}

type system struct {
	Host string `mapstructure:"host"`
	Port uint   `mapstructure:"port"`
}

type database struct {
	DBPath string `mapstructure:"db_path"`
}

type redis struct {
	Host     string `mapstructure:"host"`
	Port     uint   `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type gofound struct {
	Host string `mapstructure:"host"`
	Port uint   `mapstructure:"port"`
}

type alisms struct {
	RegionId        string `mapstructure:"regionId"`
	AccessKeyId     string `mapstructure:"accessKeyId"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
}
