package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestVideoFileRename(t *testing.T) {
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
	var nameSlice = []string{"无码频道_tg关注_@AVWUMAYUANPIAN_每天更新_", "_tg关注_@AVWUMAYUANPIAN", "_一本道_无码AV_無碼AV", "_一本道_无码AV", "_加勒比_无码AV_無碼AV", "_加勒比_无码AV", "_人妻paco_无码AV", "_天然素人_无码AV", "_#Heyzo_无码AV", "#", " "}
	var actressSlice = []string{"佐々木かな", "Heyzo-", "Vol.", "File."}
	if err := VideoFileRename(nameMap, nameSlice, actressSlice); err != nil {
		log.Fatal(err)
	}
	fmt.Println("SUCCESS")
}

func TestVideoRename(t *testing.T) {
	var nameMap = map[string]string{
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新":      "072624-001 極上泡姫物語 Vol.127   #青山茉悠",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (2)":  "071324_001 呼べば性欲処理しに来てくれる愛人  #青山茉悠",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (3)":  "032324_001 働きウーマン ~エッチな要望に寄り添うコンセルジュ~  #青山茉悠",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (4)":  "031524-001 口コミ評価満点！指名殺到の家事代行お姉さん ～これは私だけのサービスです！～  #青山茉悠",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (5)":  "050124_001 人妻の禁断不倫SEX  #青山茉悠",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (6)":  "122123-001 女熱大陸 File.096   #青山茉悠",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (7)":  "Heyzo-1627  家賃のかたにヤラれちゃう若妻  #杉浦花音",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (8)":  "Heyzo-1571 花音が教えてアゲル！～ウブな男にSEX指導～  #杉浦花音",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (9)":  "Heyzo-1519  トイレに逝ってきます～会社でオナっちゃう淫乱OL～  #杉浦花音",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (10)": "",
	}
	var nameSlice = []string{
		"无码频道_tg关注_@AVWUMAYUANPIAN_每天更新_", "_tg关注_@AVWUMAYUANPIAN",
		"_一本道_无码AV_無碼AV", "_一本道_无码AV",
		"_加勒比_无码AV_無碼AV", "_加勒比_无码AV",
		"_人妻paco_无码AV", "_天然素人_无码AV", "_#Heyzo_无码AV", "_TG频道@TBBAD", "#", " ", "_茂野美嘉_片平美嘉"}
	var actressSlice = []string{"青山茉悠", "杉浦花音", "本宮あすか", "Heyzo-", "Vol.", "File.", "No."}
	if err := VideoRename("D:/ta", nameMap, nameSlice, actressSlice); err != nil {
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
