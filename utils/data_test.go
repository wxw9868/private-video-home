package utils

import (
	"fmt"
	"log"
	"testing"
)

var nameMap = map[string]string{
	"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新":      "",
	"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (2)":  "",
	"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (3)":  "",
	"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (4)":  "",
	"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (5)":  "",
	"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (6)":  "",
	"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (7)":  "",
	"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (8)":  "",
	"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (9)":  "",
	"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (10)": "",
}

func TestVideoRename(t *testing.T) {
	var nameMap = map[string]string{
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新":      "060624_001_エロ家政婦、小泉真希の掃除クリーニングSEX！_小泉真希",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (2)":  "030224_001_肉便器育成所 ～いいなり絶対服従～_小泉真希",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (3)":  "021424_001_バイト先の人妻にお願い ～マッサージ師になりたいんです！～_小泉真希",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (4)":  "Heyzo_3251_結局、感じてしまう人妻Vol2_小泉真希",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (5)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (6)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (7)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (8)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (9)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (10)": "",
	}

	var nameSlice = []string{"无码频道_tg关注_@AVWUMAYUANPIAN_每天更新_", "_tg关注_@AVWUMAYUANPIAN", "_一本道_无码AV_無碼AV", "_#Heyzo_无码AV", "#"}
	if err := VideoRename("D:/GoLang/ta", nameMap, nameSlice); err != nil {
		log.Fatal(err)
	}
	fmt.Println("SUCCESS")
}

func TestCutVideoForGif(t *testing.T) {
	filePath := "D:/GoLang/video/assets/video/lc.mp4"
	posterPath := "D:/GoLang/video/assets/lc.gif"
	_ = CutVideoForGif(filePath, posterPath, "00:2:58")
}

func TestGeneteSQL(t *testing.T) {
	s := GeneteSQL()
	t.Log(s)
}
