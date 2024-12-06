package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"image/color/palette"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/disintegration/imaging"
	"github.com/shiningrush/avatarbuilder"
	"github.com/shiningrush/avatarbuilder/calc"
	"github.com/tidwall/gjson"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func NowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GenerateAvatar 生成头衔
func GenerateAvatar(name, path string) error {
	p := palette.Plan9
	var bgColor color.Color = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
	var frontColor color.Color = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
	bgColor = p[randInt(len(p))]
	frontColor = p[randInt(len(p))]
	ab := avatarbuilder.NewAvatarBuilder("E:/video/assets/ttf/SourceHanSansSC-Medium.ttf", &calc.SourceHansSansSCMedium{})
	ab.SetBackgroundColor(bgColor)
	ab.SetFrontgroundColor(frontColor)
	ab.SetFontSize(300)
	ab.SetAvatarSize(460, 460)
	if err := ab.GenerateImageAndSave(name, path); err != nil {
		return err
	}
	return nil
}

var (
	rng   = rand.New(rand.NewSource(time.Now().UnixNano()))
	rngMu = new(sync.Mutex)
)

func randInt(n int) int {
	rngMu.Lock()
	defer rngMu.Unlock()
	return rng.Intn(n)
}

func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				return "", err
			}

			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					return ipNet.IP.String(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("unable to determine local IP address")
}

// BoxHeader 信息头
type BoxHeader struct {
	Size       uint32
	FourccType [4]byte
	Size64     uint64
}

// GetMP4Duration 获取视频时长，以秒计
func GetMP4Duration(reader io.ReaderAt) (lengthOfTime uint32, err error) {
	var info = make([]byte, 0x10)
	var boxHeader BoxHeader
	var offset int64 = 0
	// 获取moov结构偏移
	for {
		_, err = reader.ReadAt(info, offset)
		if err != nil {
			return
		}
		boxHeader = getHeaderBoxInfo(info)
		fourccType := getFourccType(boxHeader)
		if fourccType == "moov" {
			break
		}
		// 有一部分mp4 mdat尺寸过大需要特殊处理
		if fourccType == "mdat" {
			if boxHeader.Size == 1 {
				offset += int64(boxHeader.Size64)
				continue
			}
		}
		offset += int64(boxHeader.Size)
	}
	// 获取moov结构开头一部分
	moovStartBytes := make([]byte, 0x100)
	_, err = reader.ReadAt(moovStartBytes, offset)
	if err != nil {
		return
	}
	// 定义timeScale与Duration偏移
	timeScaleOffset := 0x1C
	durationOffest := 0x20
	timeScale := binary.BigEndian.Uint32(moovStartBytes[timeScaleOffset : timeScaleOffset+4])
	duration := binary.BigEndian.Uint32(moovStartBytes[durationOffest : durationOffest+4])
	lengthOfTime = duration / timeScale
	return
}

// getHeaderBoxInfo 获取头信息
func getHeaderBoxInfo(data []byte) (boxHeader BoxHeader) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &boxHeader)
	return
}

// getFourccType 获取信息头类型
func getFourccType(boxHeader BoxHeader) (fourccType string) {
	fourccType = string(boxHeader.FourccType[:])
	return
}

// ResolveTime 将秒转成时分秒格式
func ResolveTime(seconds uint32) string {
	var (
		h, m, s string
	)
	var day = seconds / (24 * 3600)
	hour := (seconds - day*3600*24) / 3600
	minute := (seconds - day*24*3600 - hour*3600) / 60
	second := seconds - day*24*3600 - hour*3600 - minute*60
	h = strconv.Itoa(int(hour))
	if hour < 10 {
		h = "0" + strconv.Itoa(int(hour))
	}
	m = strconv.Itoa(int(minute))
	if minute < 10 {
		m = "0" + strconv.Itoa(int(minute))
	}
	s = strconv.Itoa(int(second))
	if second < 10 {
		s = "0" + strconv.Itoa(int(second))
	}
	return fmt.Sprintf("%s:%s:%s", h, m, s)
}

// CutVideoForGif 将视频剪切为 GIF
// ss 开始时间  例子：00:00:15
// t  持续时间  例子：00:00:06
// r 设定帧速率，默认为25
// s 设定画面的宽与高
func CutVideoForGif(videoPath, snapshotPath, ss string) error {
	err := ffmpeg.Input(videoPath, ffmpeg.KwArgs{"ss": ss}).
		Output(snapshotPath, ffmpeg.KwArgs{"pix_fmt": "rgb24", "t": "00:00:15", "r": "30"}).
		OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return err
	}
	return nil
}

// ReadFrameAsJpeg 将视频剪切为 JPG
func ReadFrameAsJpeg(inFileName, outFileName, ss string) error {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName, ffmpeg.KwArgs{"ss": ss}).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 0)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg", "qscale:v": 2}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		return err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		return err
	}

	if err = imaging.Save(img, outFileName); err != nil {
		return err
	}

	return nil
}

// GetVideoInfo 获取视频信息
func GetVideoInfo(inFileName string) (map[string]interface{}, error) {
	s, err := ffmpeg.Probe(inFileName)
	if err != nil {
		return nil, err
	}

	var creationTime time.Time
	if gjson.Get(s, "format.tags.creation_time").IsBool() {
		creationTime = gjson.Get(s, "format.tags.creation_time").Time().Add(0)
	} else {
		fi, err := os.Stat(inFileName)
		if err != nil {
			return nil, err
		}
		creationTime = fi.ModTime()
	}
	duration := gjson.Get(s, "format.duration").Float()
	size := gjson.Get(s, "format.size").Int()
	width := gjson.Get(s, "streams.0.width").Int()
	height := gjson.Get(s, "streams.0.height").Int()
	codecName0 := gjson.Get(s, "streams.0.codec_name").String()
	codecName1 := gjson.Get(s, "streams.1.codec_name").String()
	channelLayout := gjson.Get(s, "streams.1.channel_layout").String()

	// fmt.Println(a)
	// fmt.Printf("大小：%d\n", size)
	// fmt.Printf("尺寸：%d x %d\n", width, height)
	// fmt.Printf("编解码器：%s %s\n", strings.ToUpper(codecName1), strings.ToUpper(codecName0))
	// fmt.Printf("时长：%f\n", duration)
	// fmt.Printf("音频声道：%s\n", channelLayout)
	// fmt.Printf("创建时间：%s\n", creationTime.Format("2006:01:02 15:04:06"))

	data := map[string]interface{}{
		"duration":       duration,
		"size":           size,
		"creation_time":  creationTime,
		"width":          width,
		"height":         height,
		"codec_name0":    codecName0,
		"codec_name1":    codecName1,
		"channel_layout": channelLayout,
	}

	return data, nil
}

// Join 字符串拼接
func Join(s ...string) string {
	var b bytes.Buffer
	for i := 0; i < len(s); i++ {
		b.WriteString(s[i])
	}
	return b.String()
}

// ReadFileToMap 读取文件数据到 map
func ReadFileToMap(name string, v any) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

// AppendContentToFile 将内容追加到文件
func AppendContentToFile(filename string, b []byte) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	if _, err = f.Write(b); err != nil {
		_ = f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
}

// WriteFile 将数据写入文件
func WriteFile(name string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = os.WriteFile(name, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func VideoRename(videoDir string, nameMap map[string]string, nameSlice, actressSlice []string) error {
	files, err := os.ReadDir(videoDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("no file")
	}

	for _, file := range files {
		oldFilename := file.Name()
		filename := oldFilename

		oldName := strings.Split(filename, ".")[0]
		newName, ok := nameMap[oldName]
		if ok {
			filename = strings.Replace(filename, oldName, newName, -1)
		}
		for _, v := range nameSlice {
			filename = strings.Replace(filename, v, "", -1)
		}
		if len(filename) > 0 && filename[6:7] == "-" {
			filename = strings.Replace(filename, filename[6:7], "_", -1)
		}
		if len(filename) > 0 && !strings.Contains(filename, filename[0:10]+"_") && strings.Contains(filename, filename[0:10]) {
			newStr := ""
			if !strings.Contains(filename, "Heyzo-") {
				re := regexp.MustCompile(`\d+_\d+`)
				newStr = re.FindString(filename[0:10])
			} else {
				newStr = filename[0:10]
			}
			filename = strings.Replace(filename, filename[0:10], newStr+"_", -1)
		}
		for _, v := range actressSlice {
			if strings.Contains(filename, v) {
				if v == "Vol." {
					filename = strings.Replace(filename, v, "Vol_", -1)
				} else if v == "Heyzo-" {
					filename = strings.Replace(filename, v, "Heyzo_", -1)
				} else if v == "File." {
					filename = strings.Replace(filename, v, "File_", -1)
				} else if v == "No." {
					filename = strings.Replace(filename, v, "No_", -1)
				} else if v == "__" {
					filename = strings.Replace(filename, v, "_", -1)
				} else {
					if !strings.Contains(filename, "_"+v) {
						filename = strings.Replace(filename, v, "_"+v, -1)
					}
				}
			}
		}

		oldPath := videoDir + "/" + oldFilename
		newPath := videoDir + "/" + filename
		if err = os.Rename(oldPath, newPath); err != nil {
			return err
		}
	}

	return nil
}

func GeneteSQL() string {
	var data = make(map[string]struct{})
	err := ReadFileToMap("data.json", &data)
	if err != nil {
		return ""
	}

	var avatarDir = "./assets/image/avatar"
	var actressSql = "INSERT OR REPLACE INTO video_Actress (actress, avatar, CreatedAt, UpdatedAt) VALUES "
	var root = "/Users/v_weixiongwei/go/src/video"
	root = "D:/GoLang/video"
	for actress := range data {
		avatarPath := avatarDir + "/" + actress + ".png"
		rootPath := root + "/assets/image/avatar" + "/" + actress + ".png"
		_, err := os.Stat(rootPath)
		if os.IsNotExist(err) {
			nameSlice := []rune(actress)
			if err := GenerateAvatar(string(nameSlice[0]), rootPath); err != nil {
				log.Fatal(err)
			}
		}
		actressSql += fmt.Sprintf("('%s', '%s', '%v', '%v'), ", actress, avatarPath, time.Now().Local(), time.Now().Local())
	}
	actressSqlBytes := []byte(actressSql)
	actressSql = string(actressSqlBytes[:len(actressSqlBytes)-2])

	return actressSql
}

func GetTime(tt string) (d time.Time) {
	now := time.Now()
	year, month, day := now.Date()

	// 今日日期
	today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	switch tt {
	case "today":
		d = today
	case "yesterday":
		// 昨日日期
		d = today.AddDate(0, 0, -1)
	case "weekStart":
		// 本周起始日期（周一）
		d = today.AddDate(0, 0, -int(today.Weekday())+1)
	case "monthStart":
		// 本月起始日期
		d = time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	default:
		d = now.Local()
	}
	return
}

func GetWebDocument(method, url string, body io.Reader) (*goquery.Document, error) {
	// Request the HTML page.
	client := http.Client{
		Timeout: 50 * time.Second,
	}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Printf("wrong http request: %s", err.Error())
		return nil, fmt.Errorf("wrong http request: %s", err.Error())
	}
	if method == "POST" {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	//request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//request.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	//request.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	//request.Header.Set("Cookie", "")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error status code:", resp.StatusCode)
		return nil, fmt.Errorf("wrong http status code: %d", resp.StatusCode)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func DownloadImage(url, savePath, saveFile string) error {
	// 发起 GET 请求获取图片数据
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if saveFile == "" {
		// 获取原文件名
		_, saveFile = path.Split(res.Request.URL.Path)
	}

	// 创建保存图片的文件
	file, err := os.Create(path.Join(savePath, saveFile))
	if err != nil {
		return err
	}
	defer file.Close()

	// 将响应体的数据写入文件
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}
	return nil
}
