package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wxw9868/video/service"
)

type video struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Actress       string  `json:"actress"`
	Size          float64 `json:"size"`
	Duration      string  `json:"duration"`
	ModTime       string  `json:"mod_time"`
	Poster        string  `json:"poster"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	CodecName     string  `json:"codec_name"`
	ChannelLayout string  `json:"channel_layout"`
}

type actress struct {
	ID      int    `json:"id"`
	Actress string `json:"actress"`
	Avatar  string `json:"avatar"`
}

const (
	videoDir  = "./assets/video"
	posterDir = "./assets/image/poster"
	avatarDir = "./assets/image/avatar"
)

var (
	us = new(service.UserService)
	vs = new(service.VideoService)
	as = new(service.ActressService)
	ss = new(service.SendService)

	// snapshotDir = "./assets/image/snapshot"
	list            []string
	videos          []video
	actresss        []actress
	actressList     = make(map[string][]int)
	actressListSort []string
)

func GetUserID(c *gin.Context) uint {
	return sessions.Default(c).Get("userID").(uint)
}
