package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestVideoFileRename(t *testing.T) {
	var nameMap = map[string]string{
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新":      "022324-001 シェアガール ～令和版あなたの街のかわいい何でも屋さん～  #日南りん",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (2)":  "081823-001 逆お持ち帰り ～酒でスケベが増した2人組の美人上司に朝までヌキ尽くされました！～   #日南りん  #佐々木かな",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (3)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (4)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (5)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (6)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (7)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (8)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (9)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (10)": "",
	}
	var nameSlice = []string{
		"#", " ",
		// "无码频道_tg关注_@AVWUMAYUANPIAN_每天更新_", "_tg关注_@AVWUMAYUANPIAN",
		// "_一本道_无码AV_無碼AV", "_一本道_无码AV",
		// "_加勒比_无码AV_無碼AV", "_加勒比_无码AV",
		// "_人妻paco_无码AV", "_天然素人_无码AV", "_#Heyzo_无码AV",
	}
	var actressSlice = []string{"高垣楓", "西川ゆい", "Vol."}
	if err := VideoFileRename(nameMap, nameSlice, actressSlice); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", nameMap)
	fmt.Println("SUCCESS")
}

func TestVideoRename(t *testing.T) {
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

	var nameSlice = []string{
		"无码频道_tg关注_@AVWUMAYUANPIAN_每天更新_", "_tg关注_@AVWUMAYUANPIAN",
		"_一本道_无码AV_無碼AV", "_一本道_无码AV",
		"_加勒比_无码AV_無碼AV", "_加勒比_无码AV",
		"_人妻paco_无码AV", "_天然素人_无码AV", "_#Heyzo_无码AV", "#", " "}
	var actressSlice = []string{"三浦凛", "宮澤さおり"}
	if err := VideoRename("C:/Users/wxw9868/Downloads/ta", nameMap, nameSlice, actressSlice); err != nil {
		log.Fatal(err)
	}
	fmt.Println("SUCCESS")
}

func TestCutVideoForGif(t *testing.T) {
	filePath := "D:/video/assets/video/lc.mp4"
	posterPath := "D:/video/assets/lc.gif"
	_ = CutVideoForGif(filePath, posterPath, "00:2:58")
}

func TestGeneteSQL(t *testing.T) {
	s := GeneteSQL()
	t.Log(s)
}
