package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gofoundClient "github.com/wxw9868/video/initialize/gofound"
	"github.com/wxw9868/video/service"
)

//type video struct {
//	ID            int     `json:"id"`
//	Title         string  `json:"title"`
//	Actress       string  `json:"actress"`
//	Size          float64 `json:"size"`
//	Duration      string  `json:"duration"`
//	ModTime       string  `json:"mod_time"`
//	Poster        string  `json:"poster"`
//	Width         int     `json:"width"`
//	Height        int     `json:"height"`
//	CodecName     string  `json:"codec_name"`
//	ChannelLayout string  `json:"channel_layout"`
//}

//type actress struct {
//	ID      int    `json:"id"`
//	Actress string `json:"actress"`
//	Avatar  string `json:"avatar"`
//}

const (
	videoDir = "./assets/video"
	//posterDir = "./assets/image/poster"
	//avatarDir = "./assets/image/avatar"
)

var (
	ss    = new(service.SendService)
	us    = new(service.UserService)
	vs    = new(service.VideoService)
	as    = new(service.ActressService)
	ts    = new(service.TagService)
	stock = new(service.StockService)

	client = gofoundClient.GofoundClient()
	list   []string
	//videos          []video
	//actresss        []actress
	//actressList     = make(map[string][]int)
	//actressListSort []string

)

func GetUserID(c *gin.Context) uint {
	return sessions.Default(c).Get("user_id").(uint)
}

type Message struct {
	Code   int         `json:"code"`
	Status bool        `json:"status"`
	Msg    string      `json:"message"`
	Data   interface{} `json:"data"`
}

type Success struct {
	Code   int         `json:"code" example:"1"`
	Status bool        `json:"status" example:"true"`
	Msg    string      `json:"message" example:"status ok"`
	Data   interface{} `json:"data"`
}

type Fail struct {
	Code   int    `json:"code" example:"0"`
	Status bool   `json:"status" example:"false"`
	Msg    string `json:"message" example:"status bad request"`
}

type NotFound struct {
	Code   int    `json:"code" example:"404"`
	Status bool   `json:"status" example:"false"`
	Msg    string `json:"message" example:"status not found"`
}

type ServerError struct {
	Code   int    `json:"code" example:"500"`
	Status bool   `json:"status" example:"false"`
	Msg    string `json:"message" example:"status internal server error"`
}
