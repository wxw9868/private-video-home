package config

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var c config

func GetConfig() *config {
	return &c
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Panic("Fatal error config file: ", err)
		} else {
			fmt.Println("Config file error:", err)
		}
	}
	unmarshal()

	viper.OnConfigChange(func(e fsnotify.Event) {
		unmarshal()
		fmt.Println("Config file changed:", e.String())
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
	Alisms   alisms
}

type system struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Cert string `mapstructure:"cert"` //证书
	Key  string `mapstructure:"key"`  //证书
}

type database struct {
	DbType string `mapstructure:"db_type"`
	DbHost string `mapstructure:"db_host"`
	DbPort string `mapstructure:"db_port"`
	DbUser string `mapstructure:"db_user"`
	DbPass string `mapstructure:"db_pass"`
	DbName string `mapstructure:"db_name"`

	SslRootCert string `mapstructure:"ssl_root_cert"`
	SslKey      string `mapstructure:"ssl_key"`
	SslCert     string `mapstructure:"ssl_cert"`
}

type redis struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type alisms struct {
	RegionId        string `mapstructure:"regionId"`
	AccessKeyId     string `mapstructure:"accessKeyId"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
}
