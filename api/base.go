package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gofoundClient "github.com/wxw9868/video/initialize/gofound"
	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/service"
)

var (
	sendService    = new(service.SendService)
	userService    = new(service.UserService)
	videoService   = new(service.VideoService)
	actressService = new(service.ActressService)
	tagService     = new(service.TagService)
	stockService   = new(service.StockService)
	utilService    = new(service.UtilService)

	client = gofoundClient.GofoundClient()

	// list            []string
	// videos          []video
	// actresss        []actress
	// actressList     = make(map[string][]int)
	// actressListSort []string
)

func GetUser(c *gin.Context) *model.User {
	session := sessions.Default(c)
	user := new(model.User)
	user.ID = session.Get("user_id").(uint)
	user.Avatar = session.Get("user_avatar").(string)
	user.Username = session.Get("user_username").(string)
	user.Nickname = session.Get("user_nickname").(string)
	user.Email = session.Get("user_email").(string)
	user.Mobile = session.Get("user_mobile").(string)
	user.Designation = session.Get("user_designation").(string)
	return user
}

func GetUserID(c *gin.Context) uint {
	return GetUser(c).ID
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
