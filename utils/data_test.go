package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestVideoRename(t *testing.T) {
	var nameMap = map[string]string{
		"rh2048.com@072322_01-10mu-1080p":      "072322_01_萌 Cosplay 淫乱放纵～与顺从的巨乳女仆做爱～_大山美穂",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新":      "060124_001 肉便器育成所 ~人には言えない羞恥記録~  #森田みゆ",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (2)":  "Heyzo-2550 汗だく淫乱美女を弄んでヤッた！  #森田みゆ",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (3)":  "011224_01 出会い系で知り合ったオジサンの精子を吸い取るバリキャリOL  #河合春奈  #森田みゆ",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (4)":  "110223-001 女熱大陸 File.094 ～長めのチンコで奥まで射精して～  #森田みゆ",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (5)":  "040821_001 まんチラの誘惑 〜押しが強いナイスボディのママ友に誘われて〜  #森田みゆ",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (6)":  "031221-001 洗練された大人のいやし亭 ～可愛い狐顔のGカップ美女が、おいなりさんをにぎっておもてなし～  #森田みゆ",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (7)":  "071623-001 かり美びあんず ～女の肌の柔らかさに爆上がり～   #森田みゆ  #加藤えま",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (8)":  "062923_001 騎乗位タイムトライアル！  #森田みゆ",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (9)":  "Heyzo-3074  マザコン店長のおっぱい鑑定～この乳を探し求めていたんだ！～  #森田みゆ",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (10)": "",
		"032022-001-carib~1": "Ｇカップ巨乳痴女が3回精子を抜き取る！",
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
	s := GeneteSQL()
	t.Log(s)
}
