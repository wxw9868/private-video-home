package service

import (
	sqlite "github.com/wxw9868/video/initialize/db"
	"github.com/wxw9868/video/model"
)

var db = sqlite.DB()

var videoDir = "./assets/video"
var posterDir = "./assets/image/poster"
var avatarDir = "./assets/image/avatar"

var list []string
var videos []model.Video
var actressList = make(map[string][]int)
var actressListSort []string
