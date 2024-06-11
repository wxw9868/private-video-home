package utils

import "testing"

func TestMain(t *testing.T) {
	VideoRename("D:/GoLang/ta")
}

func TestCutVideoForGif(t *testing.T) {
	filePath := "/Users/v_weixiongwei/go/src/video/assets/video/lc.mp4"
	posterPath := "/Users/v_weixiongwei/go/src/video/assets/video/lc.gif"
	_ = CutVideoForGif(filePath, posterPath, "00:1:58")
}
