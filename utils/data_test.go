package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestVideoRename(t *testing.T) {
	var nameMap = map[string]string{
		"(1)": "",
		"(2)": "",
		"(3)": "",
		"(4)": "",
		"(5)": "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新":      "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (2)":  "092923_001_保健室の美人先生に調教されたい_小美川まゆ",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (3)":  "071423_001_乳首をハムハム！授乳プレイ ～赤ちゃんにオッパイをあげてみたいの～_小美川まゆ",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (4)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (5)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (6)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (7)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (8)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (9)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (10)": "",
	}
	if err := VideoRename("D:/GoLang/ta", nameMap); err != nil {
		log.Fatal(err)
	}
	fmt.Println("SUCCESS")
}
