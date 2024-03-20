package api

import "github.com/wxw9868/video/service"

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

type card struct {
	ID   int    `json:"id"`
	Card string `json:"card"`
}

var videoDir = "./assets/video"
var posterDir = "./assets/image/poster"
var avatarDir = "./assets/image/avatar"

// var snapshotDir = "./assets/image/snapshot"
var list []string
var videos []video
var actressList = make(map[string][]int)
var actressListSort []string

var vs = new(service.VideoService)
var as = new(service.ActressService)
