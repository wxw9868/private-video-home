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
	var nameSlice = []string{
		"无码频道_tg关注_@AVWUMAYUANPIAN_每天更新_", "_tg关注_@AVWUMAYUANPIAN",
		"_一本道_无码AV_無碼AV", "_一本道_无码AV",
		"_加勒比_无码AV_無碼AV", "_加勒比_无码AV",
		"_人妻paco_无码AV", "_天然素人_无码AV", "_#Heyzo_无码AV", "#", " "}
	var actressSlice = []string{"佐々木かな", "Heyzo-", "Vol.", "File."}
	if err := VideoFileRename(nameMap, nameSlice, actressSlice); err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%+v\n", nameMap)
	fmt.Println("SUCCESS")
}

// #川相千里  #大沢まなみ 山中麗子

func TestVideoRename(t *testing.T) {
	var nameMap = map[string]string{
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新":      "010424-001 女熱大陸 File.097   #佐々木かな",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (2)":  "Heyzo-3361 セレブ妻とオフパコ！  #セレブ妻Ｚさん  #星野さやか  #白咲花  #小林杏",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (3)":  "033024_100 スケベ椅子持参！玄関で即尺してくれる宅配ソープ嬢  #星野さやか  #白咲花  #小林杏",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (4)":  "050524_100 夫に電話をさせながら人妻をハメる ~性欲をセフレで爆発させる不倫妻~  #星野さやか  #白咲花  #小林杏",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (5)":  "Heyzo-3307 夫がリストラ！でAV出演しちゃいます！！  #星野さやか  #白咲花  #小林杏",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (6)":  "021724_985 熟女の恥ずかしい性癖  #星野さやか  #白咲花  #小林杏",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (7)":  "020124_001 隣に引っ越してきたド助平な奥さん！～ノーブラノーパンで挑発  #星野さやか  #白咲花  #小林杏",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (8)":  "031624_001 M痴女  #星野さやか  #白咲花  #小林杏",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (9)":  "110223_932 エッチがご無沙汰なエキゾチックな熟女 ~何をされてもカメラ目線~  #星野さやか",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (10)": "012024_973 連続アクメでイキまくりの豊満美巨乳な熟女をとことんヤリまくる  #星野さやか  #白咲花  #小林杏",
	}
	var nameSlice = []string{
		"无码频道_tg关注_@AVWUMAYUANPIAN_每天更新_", "_tg关注_@AVWUMAYUANPIAN",
		"_一本道_无码AV_無碼AV", "_一本道_无码AV",
		"_加勒比_无码AV_無碼AV", "_加勒比_无码AV",
		"_人妻paco_无码AV", "_天然素人_无码AV", "_#Heyzo_无码AV", "#", " ", "星野さやか白咲花"}
	var actressSlice = []string{"佐々木かな", "小林杏", "Heyzo-", "Vol.", "File."}
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
