package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestVideoRename(t *testing.T) {
	var nameMap = map[string]string{
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新":      "050423_001_行列のできる私のマンコ ～お口とマンコで３連続ザーメン採取～_今田美玲",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (2)":  "030924_001_かり美びあんず ～同棲するガチ勢はとにかくお盛ん～_今田美玲_上山奈々",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (3)":  "Heyzo_3165_男の夢！ウハウハ逆3P！Vol_12_上山奈々_今田美玲",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (4)":  "082623_001_積極的なオンナ_今田美玲",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (5)":  "100723_001_肉便器育成所 ～射精3連発～_今田美玲",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (6)":  "Heyzo_3140_続々生中～美魔女のエロボディを味わい尽くす～_今田美玲",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (7)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (8)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (9)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (10)": "",
	}

	var nameSlice = []string{"无码频道_tg关注_@AVWUMAYUANPIAN_每天更新_", "_tg关注_@AVWUMAYUANPIAN", "_一本道_无码AV_無碼AV", "_加勒比_无码AV", "_加勒比_无码AV", "_#Heyzo_无码AV", "#"}
	if err := VideoRename("D:/ta", nameMap, nameSlice); err != nil {
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
