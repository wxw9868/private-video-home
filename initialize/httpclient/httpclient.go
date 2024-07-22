package httpclient

import "github.com/wxw9868/video/utils"

func HttpClient() *utils.HttpClient {
	// api := "http://127.0.0.1:5678/api"
	api := "http://192.168.0.9:5678/api"
	return utils.NewHttpClient(api)
}
