package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/disintegration/imaging"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"
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
	src := "https://cdn.njav.tv/resize/s360/5/e5/1pondo-122922_001/thumb_h.jpg"
	//savePath := "/Users/v_weixiongwei/go/src/video/assets/image/avatar/"
	savePath := "E:/video/assets/image/thumbnail/"
	saveFile := "s360" + path.Ext(src)
	err := DownloadImage(src, savePath, saveFile)
	t.Logf("err is %s\n", err)
}

// https://github.com/PuerkitoBio/goquery
// https://github.com/gocolly/colly
// https://xslist.org/search?query=小野寺梨紗&lg=zh
// https://www.9sex.tv/cn/search?_token=eKfWaNle2cSL9iHl65TplHFLXjRmMIxzTgkYFaf0&type=actresses&query=
// https://cn.airav.wiki/?search_type=actors&lng=zh-CN&search=

func TestPachong(t *testing.T) {
	//url := Join("https://920share.com/?s=", "衣吹かのん")
	url := Join("https://ggjav.com/main/search?string=", "小泉真希")
	doc, err := GetWebDocument("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(doc)
	//av6kCom()
	//av1688Cc()
}

func av6kCom() {
	param := url.Values{"q": {"小泉真希"}}
	doc, err := GetWebDocument("POST", "https://av6k.com/plus/search.php", strings.NewReader(param.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	var page int
	doc.Find(".pages_c").Find("td").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			s.Find("b").Each(func(i int, s *goquery.Selection) {
				if i == 1 {
					page, _ = strconv.Atoi(s.Text())
				}
			})
		}
	})
	//fmt.Println(page)

	data := make(map[string]string)

	for i := 1; i <= page; i++ {
		if i > 1 {
			doc, err = GetWebDocument("GET", Join("https://av6k.com/search/", "小泉真希", "-", strconv.Itoa(i), ".html"), nil)
			if err != nil {
				log.Fatal(err)
			}
		}

		doc.Find(".newVideoC").Find(".listA").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Find("img").Attr("src")
			key := strings.Trim(s.Find(".listACT").Text(), " ")
			a := key[10:11]
			b := key[11:12]
			c := strings.Contains(key, "heyzo_hd_")
			fmt.Println(a, b, key)

			if b == "]" || a == "-" || c {
				if b == "]" {
					key = key[1:11]
					if strings.Contains(key, "Heyzo") || strings.Contains(key, "HEYZO") {
						key = strings.ToUpper(strings.Replace(key, "-", "_", -1))
					} else {
						key = strings.Replace(key, "-", "_", -1)
						key = strings.Split(key, "_")[0]
					}
				} else if a == "-" {
					key = key[0:10]
					if strings.Contains(key, "Heyzo") || strings.Contains(key, "HEYZO") {
						key = strings.ToUpper(strings.Replace(key, "-", "_", -1))
					} else {
						key = strings.ToUpper(strings.Replace(key, "-", "_", -1))
						key = strings.Split(key, "_")[0]
					}
				} else if c {
					key = strings.ToUpper(strings.Replace(strings.Split(key, " ")[0], "_hd_", "_", -1))
				}
				//fmt.Println(key)
				data[key] = Join("https://av6k.com", src)
			}
		})
		time.Sleep(time.Microsecond * 100)
	}

	fmt.Printf("%+v\n", data)

	// https://av6k.com/uploads/allimg/220813/2-220Q30U3450-L.jpg
	// https://av6k.com/uploads/allimg/220813/2-220Q30U3450-L.jpg
}

func av1688Cc() {
	url := Join("https://av1688.cc/?s=", "小泉真希")
	doc, err := GetWebDocument("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	var page int
	doc.Find(".pagination").Find("li").Each(func(i int, s *goquery.Selection) {
		if i == doc.Find(".pagination").Find("li").Length()-1 {
			page, _ = strconv.Atoi(strings.Split(s.Text(), " ")[1])
		}
	})
	//fmt.Println(page)

	data := make(map[string]string)

	for i := 1; i <= page; i++ {
		if i > 1 {
			url = Join("https://av1688.cc/page/", strconv.Itoa(i), "?s=", "小泉真希")
			doc, err = GetWebDocument("GET", url, nil)
			if err != nil {
				log.Fatal(err)
			}
		}

		doc.Find("#posts").Find("a").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Find("img").Attr("data-src")
			key, _ := s.Find("img").Attr("alt")

			Pondo := strings.Contains(key, "1pondo")
			Caribbeancom := strings.Contains(key, "Caribbeancom")
			HEYZO := strings.Contains(key, "HEYZO")
			musume := strings.Contains(key, "10musume")
			Pacopacomama := strings.Contains(key, "Pacopacomama")
			a := strings.Contains(key, "加勒比")
			b := strings.Contains(key, "一本道")
			c := strings.Contains(key, "カリビアンコム")
			d := strings.Contains(key, "加勒比PPV动画")
			f := strings.Contains(key, "Caribbeancompr-")
			g := strings.Contains(key, "一本道1pon")

			if Pondo || Caribbeancom || HEYZO || musume || Pacopacomama || a || b || c || d || f || g {
				fmt.Println(key)
				if HEYZO {
					key = key[0:10]
					key = strings.ToUpper(strings.Replace(key, "-", "_", -1))
					key = strings.ToUpper(strings.Replace(key, " ", "_", -1))
				} else if f {
					key = strings.Split(key, " ")[0]
					key = strings.Split(key, "-")[1]
					key = strings.Split(key, "_")[0]
				} else if g {
					key = strings.Split(key, " ")[0]
					m := len(key) - 10
					n := len(key)
					key = strings.Replace(key[m:n], "-", "_", -1)
					key = strings.Split(key, "_")[0]
				} else {
					key = strings.Split(key, " ")[1]
					key = strings.Replace(key, "-", "_", -1)
					key = strings.Split(key, "_")[0]
				}
				data[key] = src
			}
		})

	}

	fmt.Printf("%+v\n", data)

	// https://av1688.cc/wp-content/uploads/2024/07/20240714_6692c1d00b490.jpg
	// https://av1688.cc/wp-content/uploads/2024/07/20240714_6692c1d00b490.jpg
}

func Pachong1() {
	elems := make([]string, 3)
	elems[0] = "https://xslist.org/search?query="
	elems[1] = "小野寺梨紗"
	elems[2] = "&lg=zh"
	url := strings.Join(elems, "")
	fmt.Println(url)

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
	//client := &http.Client{
	//	Timeout: time.Millisecond * 3000,
	//}
	//res, err := client.Get(url)
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

// 去除水印
func removeWatermark(inputPath, outputPath string) error {
	// 读取原始图片
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// 判断水印位置
	bounds := img.Bounds()
	x := bounds.Dx() - 0
	y := bounds.Dy() - 30

	// 去除水印
	img = imaging.Crop(img, image.Rect(0, 0, x, y))

	// 保存处理后的图片
	err = imaging.Save(img, outputPath)
	if err != nil {
		return err
	}

	return nil
}

// 修复水印
func fixWatermark(inputPath, watermarkPath, outputPath string) error {
	// 读取原始图片和水印图片
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	watermark, err := imaging.Open(watermarkPath)
	if err != nil {
		return err
	}

	// 修复水印
	img = imaging.OverlayCenter(img, watermark, 1.0)

	// 保存处理后的图片
	err = imaging.Save(img, outputPath)
	if err != nil {
		return err
	}

	return nil
}

// https://www.cnblogs.com/Finley/p/16589798.html
func TestShuiyin(t *testing.T) {
	savePath := "E:/video/assets/image/thumbnail/"
	inputPath := "101923_001_スケスケ襦袢姿で最高のおもてなし_りおん_s360.jpg"
	outputPath := "output.jpg"
	//watermarkPath := "watermark.png"

	err := removeWatermark(savePath+inputPath, outputPath)
	if err != nil {
		fmt.Println("Failed to remove watermark:", err)
		return
	}

	//err = fixWatermark(inputPath, watermarkPath, outputPath)
	//if err != nil {
	//	fmt.Println("Failed to fix watermark:", err)
	//	return
	//}

	fmt.Println("Watermark removed and fixed successfully!")

	removeImg(savePath + inputPath)

	removeImg1(savePath + inputPath)
}

func removeWatermark1(img image.Image, watermark image.Image) image.Image {
	b := watermark.Bounds()
	// 确保水印在图片内
	if !b.In(img.Bounds()) {
		return img
	}
	// 创建水印遮罩
	mask := image.NewNRGBA(b)
	draw.Draw(mask, mask.Bounds(), image.NewUniform(color.Transparent), image.ZP, draw.Src)
	// 使用遮罩去除水印
	watermark = imaging.Paste(img, mask, image.Point{X: 667, Y: 418})
	return watermark
}

func removeImg1(inputPath string) {
	// 打开原始图片和水印图片
	src, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}
	defer src.Close()

	watermark, err := os.Open("logo2.png")
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}
	defer watermark.Close()

	// 解码图片
	img, err := jpeg.Decode(src)
	if err != nil {
		log.Fatalf("jpeg.Decode: %v", err)
	}

	watermarkImg, err := png.Decode(watermark)
	if err != nil {
		log.Fatalf("png.Decode: %v", err)
	}

	// 去除水印
	result := removeWatermark1(img, watermarkImg)

	// 保存结果
	output, err := os.Create("result.jpg")
	if err != nil {
		log.Fatalf("os.Create: %v", err)
	}
	defer output.Close()

	err = jpeg.Encode(output, result, nil)
	if err != nil {
		log.Fatalf("jpeg.Encode: %v", err)
	}
}

func removeImg(inputPath string) {
	// 打开原始图片
	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// 读取原始图片
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建一个与原始图片大小相同的透明画布
	canvas := image.NewRGBA(img.Bounds())

	// 将原始图片绘制到画布上
	draw.Draw(canvas, img.Bounds(), img, image.Point{X: 0, Y: 0}, draw.Src)

	// 在画布上去除某一个对象（这里以一个矩形框为例）
	rect := image.Rect(667, 418, 800, 450)
	draw.Draw(canvas, rect, &image.Uniform{color.Transparent}, image.Point{X: 0, Y: 0}, draw.Src)

	// 存储处理后的图片
	outFile, err := os.Create("removed.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outFile.Close()
	png.Encode(outFile, canvas)
	fmt.Println("图片去除成功！")
}
