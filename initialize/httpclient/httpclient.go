package httpclient

import "github.com/wxw9868/video/utils"

func HttpClient() *utils.HttpClient {
	return utils.NewHttpClient("http://127.0.0.1:5678/api")
}
