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
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新":      "123123_001_蝶が如く ～ピンク通りの二輪車ソープランド23～_白川麻衣_石川さとみ",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (2)":  "082223_001_ボクを狂わせる家庭教師のおねえさん ～胸に触れた瞬間に何かが崩壊しました～_白川麻衣",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (3)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (4)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (5)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (6)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (7)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (8)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (9)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (10)": "",
	}

	var nameSlice = []string{"无码频道_tg关注_@AVWUMAYUANPIAN_每天更新_", "_tg关注_@AVWUMAYUANPIAN", "_一本道_无码AV_無碼AV", "_加勒比_无码AV", "_加勒比_无码AV", "_#Heyzo_无码AV", "#"}
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
