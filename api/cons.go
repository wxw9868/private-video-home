package api

type video struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Actress  string  `json:"actress"`
	Size     float64 `json:"size"`
	Duration string  `json:"duration"`
	ModTime  string  `json:"mod_time"`
	Poster   string  `json:"poster"`
}

type actress struct {
	ID      int    `json:"id"`
	Actress string `json:"actress"`
	Avatar  string `json:"avatar"`
}

var videoDir = "./assets/video"
var posterDir = "./assets/image/poster"
var avatarDir = "./assets/image/avatar"

// var snapshotDir = "./assets/image/snapshot"
var list []string
var videos []video
var actressList = make(map[string][]int)
var actressListSort []string
