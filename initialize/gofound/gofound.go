package httpclient

import (
	"fmt"
	"github.com/wxw9868/video/config"
	"github.com/wxw9868/video/utils"
)

func GofoundClient() *utils.GofoundClient {
	conf := config.Config().Gofound
	api := fmt.Sprintf("http://%s:%d/api", conf.Host, conf.Port)
	return utils.NewGofoundClient(api)
}
