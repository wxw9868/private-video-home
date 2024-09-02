package httpclient

import (
	"fmt"
	"log"

	"github.com/wxw9868/video/config"
	"github.com/wxw9868/video/utils"
)

func GofoundClient() *utils.GofoundClient {
	conf := config.Config().Gofound
	api := fmt.Sprintf("http://%s:%d/api", conf.Host, conf.Port)
	return utils.NewGofoundClient(api)
}

func init() {
	resp, err := GofoundClient().GET("/status")
	if err != nil {
		log.Fatalf("GoFound连接失败: %s", err)
	}
	defer resp.Body.Close()
	log.Println("GoFound服务连接成功")
}
