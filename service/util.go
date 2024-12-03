package service

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
)

type UtilService struct{}

func (us *UtilService) OneAddInfoToActress(actress string) error {
	url1 := utils.Join("https://xslist.org/search?query=", actress, "&lg=zh")
	doc, err := utils.GetWebDocument("GET", url1, nil)
	if err != nil {
		return err
	}
	if doc.Find("body").Text() == "No results found." {
		return errors.New("no results found")
	}
	href, _ := doc.Find("a").Attr("href")

	doc, err = utils.GetWebDocument("GET", href, nil)
	if err != nil {
		return err
	}
	sss1 := doc.Find("#sss1")
	alias := sss1.Find("p").Text()
	avatar, _ := sss1.Find("img").Attr("src")

	var savePath string
	switch runtime.GOOS {
	case "linux":
	case "darwin":
		savePath = "/Users/v_weixiongwei/go/src/video/assets/image/avatar/"
	case "windows":
		savePath = "E:/video/assets/image/avatar/"
	}
	_, saveFile := path.Split(href)
	err = utils.DownloadImage(avatar, savePath, saveFile)
	if err != nil {
		return err
	}

	if alias != "" {
		alias = strings.Trim(strings.Split(alias, ":")[1], " ")
	}
	m := model.Actress{}
	m.Alias = alias
	m.Avatar = avatar
	doc.Find("h2").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			personal, _ := s.Next().Html()
			personal = strings.Replace(strings.Replace(strings.Replace(personal, "<span itemprop=\"height\">", "", -1), "<span itemprop=\"nationality\">", "", -1), "</span>", "", -1)
			personals := strings.Split(personal, "<br/>")
			birth := personals[0]                  // 出生
			measurements := personals[1]           // 三围
			cupSze := personals[2]                 // 罩杯
			debutDate := personals[3]              // 出道日期
			starSign := personals[4]               // 星座
			bloodGroup := personals[5]             // 血型
			stature := personals[6]                // 身高
			nationality := personals[7]            // 国籍
			introduction := s.Next().Next().Text() // 简介

			m.Birth = strings.Trim(strings.Split(birth, ":")[1], " ")
			m.Measurements = strings.Trim(strings.Split(measurements, ":")[1], " ")
			m.CupSize = strings.Trim(strings.Split(cupSze, ":")[1], " ")
			m.DebutDate = strings.Trim(strings.Split(debutDate, ":")[1], " ")
			m.StarSign = strings.Trim(strings.Split(starSign, ":")[1], " ")
			m.BloodGroup = strings.Trim(strings.Split(bloodGroup, ":")[1], " ")
			m.Stature = strings.Trim(strings.Split(stature, ":")[1], " ")
			m.Nationality = strings.Trim(strings.Split(nationality, ":")[1], " ")
			m.Introduction = strings.Trim(strings.Split(introduction, ":")[1], " ")
		}
	})
	//fmt.Printf("%+v\n", m)
	if err = db.Model(&actress).Updates(m).Error; err != nil {
		return err
	}
	time.Sleep(time.Microsecond * 300)

	return nil
}

func (us *UtilService) AllAddInfoToActress() error {
	var actresss []model.Actress
	if err := db.Where("birth is null or birth = ''").Find(&actresss).Error; err != nil {
		return err
	}

	numCPU := runtime.NumCPU()
	ch := make(chan string, numCPU)
	wg := new(sync.WaitGroup)

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(ch chan string, wg *sync.WaitGroup, i int) {
			defer wg.Done()
			for actress := range ch {
				err := us.OneAddInfoToActress(actress)
				fmt.Printf("index: %d, actress: %s, error: %s\n", i, actress, err)
			}
		}(ch, wg, i)
	}

	for _, actress := range actresss {
		ch <- actress.Actress
	}
	close(ch)

	wg.Wait()

	return nil
}

func (us *UtilService) AllVideoPic(page, pageSize int, site string) error {
	var count int64
	if err := db.Model(&model.Video{}).Count(&count).Error; err != nil {
		return err
	}

	var actresss []model.Actress
	if err := db.Scopes(Paginate(page, pageSize, int(count))).Find(&actresss).Error; err != nil {
		return err
	}

	numCPU := runtime.NumCPU()
	ch := make(chan string, numCPU)
	wg := new(sync.WaitGroup)

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(ch chan string, site string, wg *sync.WaitGroup, i int) {
			wg.Done()

			for actress := range ch {
				var err error
				var videos []model.Video
				err = db.Where("actress = ?", actress).Find(&videos).Error
				fmt.Println(err)

				switch site {
				case "av1688Cc":
					err = av1688Cc(actress, videos)
				case "av6kCom":
					err = av6kCom(actress, videos)
				default:
					err = njavTv(actress, videos)
				}
				fmt.Printf("index: %d, actress: %s, error: %s\n", i, actress, err)
			}
		}(ch, site, wg, i)
	}

	for _, actress := range actresss {
		ch <- actress.Actress
	}
	close(ch)

	wg.Wait()

	//data := make(map[string]map[string]string)
	//err := utils.WriteMapToFile("E:/video/assets/data/data.json", data)
	return nil
}

func (us *UtilService) OneVideoPic(actress string, site string) error {
	var videos []model.Video
	if err := db.Where("actress = ?", actress).Find(&videos).Error; err != nil {
		return err
	}

	var err error
	switch site {
	case "av1688Cc":
		err = av1688Cc(actress, videos)
	case "av6kCom":
		err = av6kCom(actress, videos)
	default:
		err = njavTv(actress, videos)
	}

	return err
}

func av1688Cc(actress string, videos []model.Video) error {
	doc, err := utils.GetWebDocument("GET", utils.Join("https://av1688.cc/?s=", actress), nil)
	if err != nil {
		return err
	}

	//fmt.Println(actress)

	var page int
	doc.Find(".pagination").Find("li").Each(func(i int, s *goquery.Selection) {
		if i == doc.Find(".pagination").Find("li").Length()-1 {
			page, _ = strconv.Atoi(strings.Split(s.Text(), " ")[1])
		}
	})

	data := make(map[string]string)

	for i := 1; i <= page; i++ {
		if i > 1 {
			doc, err = utils.GetWebDocument("GET", utils.Join("https://av1688.cc/page/", strconv.Itoa(i), "?s=", actress), nil)
			if err != nil {
				return err
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
	//fmt.Printf("%+v\n", data)
	// https://av1688.cc/wp-content/uploads/2024/07/20240714_6692c1d00b490.jpg
	// https://av1688.cc/wp-content/uploads/2024/07/20240714_6692c1d00b490.jpg
	for i := 0; i < len(videos); i++ {
		video := videos[i]

		HEYZO := strings.Contains(video.Title, "Heyzo")
		title := ""
		if HEYZO {
			title = video.Title[0:10]
		} else {
			title = video.Title[0:6]
		}

		//fmt.Println(title)

		src, ok := data[strings.ToUpper(title)]
		if ok {
			src = strings.Split(src, "?")[0]
			savePath := "E:/video/assets/image/thumbnail/"
			saveFile := video.Title + "_s360" + path.Ext(src)
			_, err = os.Stat(path.Join(savePath, saveFile))
			if os.IsNotExist(err) {
				err = utils.DownloadImage(src, savePath, saveFile)
				if err != nil {
					return err
				}
			}

			savePath = "E:/video/assets/image/preview/"
			saveFile = video.Title + path.Ext(src)
			_, err = os.Stat(path.Join(savePath, saveFile))
			if os.IsNotExist(err) {
				err = utils.DownloadImage(src, savePath, saveFile)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func av6kCom(actress string, videos []model.Video) error {
	param := url.Values{"q": {actress}}
	doc, err := utils.GetWebDocument("POST", "https://av6k.com/plus/search.php", strings.NewReader(param.Encode()))
	if err != nil {
		return err
	}

	data := make(map[string]string)

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

	for i := 1; i <= page; i++ {
		if i > 1 {
			doc, err = utils.GetWebDocument("GET", utils.Join("https://av6k.com/search/", "小泉真希", "-", strconv.Itoa(i), ".html"), nil)
			if err != nil {
				return err
			}
		}

		doc.Find(".newVideoC").Find(".listA").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Find("img").Attr("src")
			key := strings.Trim(s.Find(".listACT").Text(), " ")
			a := key[10:11]
			b := key[11:12]
			c := strings.Contains(key, "heyzo_hd_")
			//fmt.Println(a, b, key)

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
				data[key] = utils.Join("https://av6k.com", src)
			}
		})
		time.Sleep(time.Microsecond * 50)
	}
	//fmt.Printf("%+v\n", data)

	for i := 0; i < len(videos); i++ {
		video := videos[i]

		HEYZO := strings.Contains(video.Title, "Heyzo")
		title := ""
		if HEYZO {
			title = video.Title[0:10]
		} else {
			title = video.Title[0:6]
		}

		//fmt.Println(title)

		src, ok := data[strings.ToUpper(title)]
		if ok {
			src = strings.Split(src, "?")[0]
			savePath := "E:/video/assets/image/thumbnail/"
			saveFile := video.Title + "_s360" + path.Ext(src)
			_, err = os.Stat(path.Join(savePath, saveFile))
			if os.IsNotExist(err) {
				err = utils.DownloadImage(src, savePath, saveFile)
				if err != nil {
					return err
				}
			}

			savePath = "E:/video/assets/image/preview/"
			saveFile = video.Title + path.Ext(src)
			_, err = os.Stat(path.Join(savePath, saveFile))
			if os.IsNotExist(err) {
				err = utils.DownloadImage(src, savePath, saveFile)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func njavTv(actress string, videos []model.Video) error {
	p := math.Ceil(float64(len(videos))/12) * 2
	data := make(map[string]string)

	for i := 1; i < int(p); i++ {

	start:
		elems := make([]string, 3)
		elems[0] = "https://njav.tv/zh/search?keyword="
		elems[1] = actress
		elems[2] = "&page=" + strconv.Itoa(i)
		url1 := strings.Join(elems, "")

		doc, err := utils.GetWebDocument("GET", url1, nil)
		if err != nil {
			time.Sleep(time.Second * 2)
			goto start
		}

		doc.Find(".box-item").Each(func(i int, s *goquery.Selection) {
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				if i == 0 {
					src, _ := s.Find("img").Attr("data-src")
					key, _ := s.Find("img").Attr("title")

					Pondo := strings.Contains(key, "1Pondo")
					Caribbeancom := strings.Contains(key, "Caribbeancom")
					HEYZO := strings.Contains(key, "HEYZO")
					musume := strings.Contains(key, "10musume")
					Pacopacomama := strings.Contains(key, "Pacopacomama")

					if Pondo || Caribbeancom || HEYZO || musume || Pacopacomama {
						if Pondo || Caribbeancom || musume || Pacopacomama {
							key = strings.Split(key, "-")[1]
							key = strings.Split(key, "_")[0]
						} else {
							key = strings.ToUpper(strings.Replace(key, "-", "_", -1))
						}
						data[key] = src
					}
				}
				//else {
				//	data["title"] = s.Text()
				//}
			})
		})
		time.Sleep(time.Microsecond * 300)
	}
	//fmt.Printf("%+v\n", data)
	//data[actress.Actress] = data

	for i := 0; i < len(videos); i++ {
		video := videos[i]

		HEYZO := strings.Contains(video.Title, "Heyzo")
		title := ""
		if HEYZO {
			title = video.Title[0:10]
		} else {
			title = video.Title[0:6]
		}

		//fmt.Println(title)

		src, ok := data[strings.ToUpper(title)]
		if ok {
			src = strings.Split(src, "?")[0]
			savePath := "E:/video/assets/image/thumbnail/"
			saveFile := video.Title + "_s360" + path.Ext(src)
			_, err := os.Stat(path.Join(savePath, saveFile))
			if os.IsNotExist(err) {
				err = utils.DownloadImage(src, savePath, saveFile)
				if err != nil {
					return err
				}
			}

			savePath = "E:/video/assets/image/preview/"
			src = strings.Replace(src, "resize/s360", "images", -1)
			saveFile = video.Title + path.Ext(src)
			_, err = os.Stat(path.Join(savePath, saveFile))
			if os.IsNotExist(err) {
				err = utils.DownloadImage(src, savePath, saveFile)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
