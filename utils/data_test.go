package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestVideoRename(t *testing.T) {
	var nameMap = map[string]string{
		"2":     "123014_001 爆乳激震！超絶潮吹きオマンコに連続中出し  #赤井美月  #折原ほのか",
		"2 (2)": "090122_001 性欲が満たされない人妻と隣人の禁断関係  #折原ほのか ",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新":      "041324-001 かり美びあんず ～レンタルルームで貝合わせ巨乳カップル～   #折原ほのか  #小衣くるみ",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (2)":  "092623-001 恍惚 ～ひとたび挿入されたら離れられない～  #折原ほのか  ",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (3)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (4)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (5)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (6)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (7)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (8)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (9)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (10)": "",
	}
	var nameSlice = []string{"无码频道_tg关注_@AVWUMAYUANPIAN_每天更新_", "_tg关注_@AVWUMAYUANPIAN", "_#Heyzo_无码AV", "#"}
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
	GeneteSQL()
}
