package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/url"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/disintegration/imaging"
	"golang.org/x/image/draw"
)

func TestNowTime(t *testing.T) {
	NowTime()
}

func TestGetLocalIP(t *testing.T) {
	ip, err := GetLocalIP()
	t.Logf("ip: %s err: %s\n", ip, err)
}

func TestAppendContentToFile(t *testing.T) {
	err := AppendContentToFile("test.log", []byte("test"))
	t.Log(err)
}

func TestWriteFile(t *testing.T) {
	err := WriteFile("test.log", "test")
	t.Log(err)
}

func TestGetMP4Duration(t *testing.T) {
	videoDir := "./assets/video"
	filePath := videoDir + "/" + "test.mp4"
	fil, _ := os.Open(filePath)
	duration, _ := GetMP4Duration(fil)
	t.Log(duration)
}

func TestReadFrameAsJpeg(t *testing.T) {
	videoDir := "./assets/video"
	posterDir := "./assets/image/poster"
	filePath := videoDir + "/" + "test.mp4"
	posterPath := posterDir + "/test.jpg"
	_, err := os.Stat(posterPath)
	if os.IsNotExist(err) {
		_ = ReadFrameAsJpeg(filePath, posterPath, "00:02:30")
	}
}

var _ = map[string]string{
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
	"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (11)": "",
	"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (12)": "",
	"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (13)": "",
	"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (14)": "",
	"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (15)": "",
}

func TestVideoRename(t *testing.T) {
	var nameMap = map[string]string{
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (10)": "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (11)": "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (12)": "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (13)": "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (14)": "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (15)": "",
	}
	// #上原聡美  #さとみ
	// #大月のの  #中津井加代
	// 朝桐光    南野あかり
	// 宮村恋  #歩
	var nameSlice = []string{
		"无码频道_tg关注_@AVWUMAYUANPIAN_每天更新_", "_tg关注_@AVWUMAYUANPIAN",
		"_一本道_无码AV_無碼AV", "_一本道_无码AV",
		"_加勒比_无码AV_無碼AV", "_加勒比_无码AV",
		"_人妻paco_无码AV", "_天然素人_无码AV", "_#Heyzo_无码AV", "_TG频道@TBBAD", "歩", "#", " "}
	var actressSlice = []string{"本真ゆり", "西内萌菜", "宮村恋", "田所三久", "栗原梢", "夏希アンジュ", "Heyzo-", "Vol.", "File.", "No."}
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

	var savePath string

	switch runtime.GOOS {
	case "linux":
	case "darwin":
		savePath = "/Users/v_weixiongwei/go/src/video/assets/image/thumbnail/"
	case "windows":
		savePath = "E:/video/assets/image/thumbnail/"
	}

	saveFile := "_s360" + path.Ext(src)
	err := DownloadImage(src, savePath, saveFile)
	t.Logf("err is %s\n", err)
}

func TestGetWebDocument(t *testing.T) {
	//url := Join("https://920share.com/?s=", "衣吹かのん")
	//url := Join("https://ggjav.com/main/search?string=", "小泉真希")
	//url := Join("https://netflav.com/search?type=title&keyword=", "杉浦花音")

	//url := Join("https://jable.tv/search/", "杉浦花音", "/")
	//url := Join("https://missav.com/cn/search/", "杉浦花音")
	//url := Join("https://dgaqp.com/search/", "杉浦花音")
	//url := Join("https://supjav.com/zh/?s=", "杉浦花音")
	//url := Join("https://ggjav.com/main/search?string", "杉浦花音")
	//Nanako Asahina, まーちゃん, 今井花菜, 小池愛菜, 恋乃みくる, 朝比奈京子, 朝比奈恭子, 浅野麻衣, 野々村あいり, 陽菜子,上岡流留香,冴島みのり,ななこ,せりな・愛・ちひろ,モカ

	//
	var url1 = Join("https://nowav.tv/?s=", "亀井ひとみ")
	doc, err := GetWebDocument("GET", url1, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(doc)
}

func TestAv6kCom(t *testing.T) {
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
				data[key] = Join("https://av6k.com", src)
			}
		})
		time.Sleep(time.Microsecond * 100)
	}

	fmt.Printf("%+v\n", data)
}

func TestAv1688Cc(t *testing.T) {
	url1 := Join("https://av1688.cc/?s=", "小泉真希")
	doc, err := GetWebDocument("GET", url1, nil)
	if err != nil {
		log.Fatal(err)
	}

	var page int
	doc.Find(".pagination").Find("li").Each(func(i int, s *goquery.Selection) {
		if i == doc.Find(".pagination").Find("li").Length()-1 {
			page, _ = strconv.Atoi(strings.Split(s.Text(), " ")[1])
		}
	})

	data := make(map[string]string)

	for i := 1; i <= page; i++ {
		if i > 1 {
			url1 = Join("https://av1688.cc/page/", strconv.Itoa(i), "?s=", "小泉真希")
			doc, err = GetWebDocument("GET", url1, nil)
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
}

func TestXslist(t *testing.T) {
	//瀬戸レイカ,目々澤めぐ,さおり,優希まこと,和泉紫乃,柏木もも,大久保ゆう,広瀬里香,上野真奈美,小鳥遊つばさ,川越ゆい,早坂咲重,海野真凜,翼みさき,桜田桃羽,榊原ひとみ,須間あいり,高田伸子,三花れな,美波ゆさ,小嶋ひより,花咲胡桃
	url1 := Join("https://xslist.org/search?query=", "優希まこと", "&lg=zh")
	doc, err := GetWebDocument("GET", url1, nil)
	if err != nil {
		log.Fatal(err)
	}
	href, _ := doc.Find("a").Attr("href")

	doc, err = GetWebDocument("GET", href, nil)
	if err != nil {
		log.Fatal(err)
	}
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
			// fmt.Printf("%+v\n", personals)
			for i2, s2 := range personals {
				fmt.Printf("i is %d personal is %s \n", i2, s2)
			}
			// birth := personals[0]                  // 出生
			// measurements := personals[1]           // 三围
			// cup_size := personals[2]               // 罩杯
			// debut_date := personals[3]             // 出道日期
			// star_sign := personals[4]              // 星座
			// blood_group := personals[5]            // 血型
			// stature := personals[6]                // 身高
			// nationality := personals[7]            // 国籍
			// introduction := s.Next().Next().Text() // 简介
			// fmt.Printf("Introduction is %s \n", introduction)
		}
	})
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
}

//func removeWatermark1(img image.Image, watermark image.Image) image.Image {
//	b := watermark.Bounds()
//	// 确保水印在图片内
//	if !b.In(img.Bounds()) {
//		return img
//	}
//	// 创建水印遮罩
//	mask := image.NewNRGBA(b)
//	draw.Draw(mask, mask.Bounds(), image.NewUniform(color.Transparent), image.ZP, draw.Src)
//	// 使用遮罩去除水印
//	watermark = imaging.Paste(img, mask, image.Point{X: 667, Y: 418})
//	return watermark
//}

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
	draw.Draw(canvas, rect, &image.Uniform{C: color.Transparent}, image.Point{X: 0, Y: 0}, draw.Src)

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

//func TestMain(T *testing.T) {
//	// 加载图像
//	img := gocv.IMRead("watermarked.jpg", gocv.IMReadColor)
//	defer img.Close()
//
//	// 检查图像是否成功加载
//	if img.Empty() {
//		fmt.Printf("Error reading image\n")
//		return
//	}
//
//	// 在这里添加你的图像处理代码来尝试去除水印
//	// ...
//
//	// 保存处理后的图像
//	gocv.IMWrite("output.jpg", img)
//}
