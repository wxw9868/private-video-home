package utils

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
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

func TestDownloadImage(t *testing.T) {
	url := "https://www.golang-mix.com/imgs/user.png"
	savePath := "/Users/v_weixiongwei/go/src/video/assets/image/avatar/"
	_, err := DownloadImage(url, savePath)
	t.Logf("err is %s\n", err)
}

// https://github.com/PuerkitoBio/goquery
// https://github.com/gocolly/colly
// https://xslist.org/search?query=小野寺梨紗&lg=zh
// https://www.9sex.tv/cn/search?_token=eKfWaNle2cSL9iHl65TplHFLXjRmMIxzTgkYFaf0&type=actresses&query=
// https://cn.airav.wiki/?search_type=actors&lng=zh-CN&search=

func TestPachong(t *testing.T) {
	Pachong2()
	Pachong3()
}

func Pachong3() {
	elems := make([]string, 2)
	elems[0] = "https://cn.airav.wiki/?search_type=actors&lng=zh-CN&search="
	elems[1] = "小野寺梨紗"
	url := strings.Join(elems, "")
	fmt.Printf("url is %s \n", url)

	doc := getDoc(url)
	fmt.Println(doc)
	return
	href, _ := doc.Find("a").Attr("href")

	doc = getDoc(href)
	doc.Find("h2").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			title := s.Text()
			fmt.Printf("title is %s \n", title)
		}
	})
}

func Pachong2() {
	elems := make([]string, 2)
	elems[0] = "https://www.9sex.tv/cn/search?_token=eKfWaNle2cSL9iHl65TplHFLXjRmMIxzTgkYFaf0&type=actresses&query="
	elems[1] = "小野寺梨紗"
	url := strings.Join(elems, "")
	fmt.Printf("url is %s \n", url)

	doc := getDoc(url)
	fmt.Println(doc)
	doc.Find("li > a").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			title := s.Text()
			fmt.Printf("title is %s \n", title)
			return
		}
	})

	href, _ := doc.Find("a").Attr("href")
	doc = getDoc(href)
	doc.Find("h2").Each(func(i int, s *goquery.Selection) {
		if i == 0 {

		}
	})
}

func Pachong1() {
	elems := make([]string, 3)
	elems[0] = "https://xslist.org/search?query="
	elems[1] = "小野寺梨紗"
	elems[2] = "&lg=zh"
	url := strings.Join(elems, "")

	doc := getDoc(url)
	href, _ := doc.Find("a").Attr("href")

	doc = getDoc(href)
	actress := doc.Find("#sss1").Find("header").Text()
	alias := doc.Find("#sss1").Find("p").Text()
	avatar, _ := doc.Find("#sss1").Find("img").Attr("src")
	fmt.Printf("actress is %s \n", strings.Trim(actress, " "))
	fmt.Printf("alias is %s \n", alias)
	fmt.Printf("avatar is %s \n", avatar)
	doc.Find("h2").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			title := s.Text()
			fmt.Printf("title is %s \n", title)
			personal, _ := s.Next().Html()
			personal = strings.Replace(strings.Replace(strings.Replace(personal, "<span itemprop=\"height\">", "", -1), "<span itemprop=\"nationality\">", "", -1), "</span>", "", -1)
			personals := strings.Split(personal, "<br/>")
			//fmt.Println(personals)
			// birth := personals[0]        // 出生
			// measurements := personals[1] // 三围
			// cup_size := personals[2]     // 罩杯
			// debut_date := personals[3]   // 出道日期
			// star_sign := personals[4]    // 星座
			// blood_group := personals[5]  // 血型
			// stature := personals[6]      // 身高
			// nationality := personals[7]  // 国籍
			for i2, s2 := range personals {
				fmt.Printf("i is %d personal is %s \n", i2, s2)
			}
			introduction := s.Next().Next().Text() // 简介
			fmt.Printf("Introduction is %s \n", introduction)
		}
	})
}

func getDoc(url string) *goquery.Document {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}
